package debug

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func checkIfStructContainSingleVisibleField(reflectValue reflect.Value) (
	reflect.Value, reflect.StructField, bool, bool,
) {
	if reflectValue.Kind() == reflect.Interface {
		reflectValue = reflectValue.Elem()
	}
	for reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflectValue.Elem()
	}

	if reflectValue.Kind() != reflect.Struct {
		return reflectValue, reflect.StructField{}, false, false
	}

	fields := reflect.VisibleFields(reflectValue.Type())
	fieldFoundCount := 0
	var found reflect.Value
	var foundField reflect.StructField
	for _, field := range fields {
		foundDebug := field.Tag.Get("debug")
		if foundDebug == "" {
			//log.Printf("missing tag on %s %s", reflectValue.Type().Name(), field.Name)
			continue
		}

		reflectFieldValue := reflectValue.FieldByIndex(field.Index)
		if reflectFieldValue.Kind() == reflect.Slice {
			if reflectFieldValue.Len() == 0 {
				continue
			}
		}
		found = reflectFieldValue
		foundField = field
		fieldFoundCount++
	}

	if fieldFoundCount != 1 {
		return reflect.Value{}, reflect.StructField{}, false, true
	}

	return found, foundField, true, true
}

func writeAllSingleFieldStructs(fieldValue reflect.Value, structField reflect.StructField, localDepth int,
	writer io.Writer) (
	reflect.Value, bool,
) {
	if fieldValue.Kind() == reflect.Slice {
		if fieldValue.Len() == 0 {
			return reflect.Value{}, false
		}
	}
	if localDepth > 0 {
		fmt.Fprintf(writer, ".")
	}
	fmt.Fprintf(writer, "%s", structField.Name)

	foundFieldValue, foundSingleField, fieldWasFound, _ := checkIfStructContainSingleVisibleField(fieldValue)
	if fieldWasFound {
		if localDepth == 0 {
			fmt.Fprintf(writer, "(%s)", fieldValue.Type().String())
		}
		return writeAllSingleFieldStructs(foundFieldValue, foundSingleField, localDepth+1, writer)
	}
	fmt.Fprintf(writer, "(%s)", fieldValue.Type().String())
	return fieldValue, true
}

func getTrueTypeName(reflectValue reflect.Value) string {
	if reflectValue.Kind() == reflect.Interface {
		elemValue := reflectValue.Elem()
		reflectValue = elemValue
	}

	if !reflectValue.IsValid() {
		return "NIL"
	}

	return reflectValue.Type().String()
}

func tree(reflectValue reflect.Value, tab int, writer io.Writer) {
	if reflectValue.Kind() == reflect.Interface {
		reflectValue = reflectValue.Elem()
	}
	for reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflectValue.Elem()
	}

	tabs := strings.Repeat("..", tab)

	if !reflectValue.IsValid() {
		fmt.Fprintf(writer, "nil\n")
		return
	}

	switch reflectValue.Kind() {
	case reflect.Struct:
		fields := reflect.VisibleFields(reflectValue.Type())
		fieldWrittenCount := 0
		filteredOut := false
		for _, field := range fields {
			foundDebug := field.Tag.Get("debug")
			if foundDebug == "" {
				//log.Printf("missing tag on %s %s", reflectValue.Type().Name(), field.Name)
				continue
			}

			if len(field.Index) > 1 {
				//				filteredOut = true
				//				continue
			}
			reflectFieldValue := reflectValue.FieldByIndex(field.Index)

			var localWriter bytes.Buffer
			subReflectFieldValue, shouldUseThisFieldAfterAll := writeAllSingleFieldStructs(
				reflectFieldValue, field, 0, &localWriter,
			)
			if !shouldUseThisFieldAfterAll {
				filteredOut = true
				continue
			}

			if fieldWrittenCount >= 0 {
				fmt.Fprintf(writer, "\n%s", tabs)
			}

			writer.Write(localWriter.Bytes())

			fieldWrittenCount++
			tree(subReflectFieldValue, tab+1, writer)
		}

		if fieldWrittenCount == 0 && !filteredOut {
			panic(fmt.Errorf("not found any debug fields on '%s'", reflectValue.Type().String()))
		}

	case reflect.Slice:
		for i := 0; i < reflectValue.Len(); i++ {
			fmt.Fprintf(writer, "\n%s", tabs)
			item := reflectValue.Index(i)
			fmt.Fprintf(writer, "%d: (%s)", i, getTrueTypeName(item))
			tree(item, tab+1, writer)
		}

	case reflect.String:
		fmt.Fprintf(writer, " = '%s'", reflectValue.String())

	case reflect.Int:
	case reflect.Int32:
		fmt.Fprintf(writer, " = %d", reflectValue.Int())
	case reflect.Bool:
		fmt.Fprintf(writer, " = %v", reflectValue.Bool())
	default:
		fmt.Fprintf(writer, "unknown %d", reflectValue.Kind())
	}
}

func Tree(expr interface{}, writer io.Writer) {
	subType := reflect.ValueOf(expr)
	fmt.Fprintf(writer, "\n%s", getTrueTypeName(subType))
	tree(subType, 1, writer)
}

func TreeString(expr interface{}) string {
	var writer strings.Builder
	Tree(expr, &writer)
	return writer.String()
}
