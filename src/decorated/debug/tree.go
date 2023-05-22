package debug

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
)

func writeStructField(field reflect.StructField, tabs string, writer io.Writer) {
	fmt.Fprintf(writer, "\n%s%v", tabs, field.Name)
}

func tree(reflectValue reflect.Value, tab int, writer io.Writer) {
	if reflectValue.Kind() == reflect.Interface {
		reflectValue = reflectValue.Elem()
	}
	for reflectValue.Kind() == reflect.Pointer {
		reflectValue = reflectValue.Elem()
	}

	tabs := strings.Repeat("..", tab)

	switch reflectValue.Kind() {
	case reflect.Struct:
		fields := reflect.VisibleFields(reflectValue.Type())
		fieldWrittenCount := 0
		for _, field := range fields {
			foundDebug := field.Tag.Get("debug")
			if foundDebug == "" {
				log.Printf("missing tag on %s %s", reflectValue.Type().Name(), field.Name)
				continue
			}
			reflectField := reflectValue.FieldByIndex(field.Index)
			if fieldWrittenCount == 0 {
				fmt.Fprintf(writer, " %s", reflectValue.Type().String())
			}
			writeStructField(field, tabs, writer)
			fieldWrittenCount++
			tree(reflectField, tab+1, writer)
		}
		if fieldWrittenCount == 0 {
			panic(fmt.Errorf("not allowed '%s'", reflectValue.Type().String()))
		}
	case reflect.Slice:
		for i := 0; i < reflectValue.Len(); i++ {
			if i >= 0 {
				fmt.Fprintf(writer, "\n%s", tabs)
			}
			item := reflectValue.Index(i)
			fmt.Fprintf(writer, "%d: ", i)
			tree(item, tab+1, writer)
		}
	case reflect.String:
		fmt.Fprintf(writer, " = '%s'", reflectValue.String())
	case reflect.Int:
	case reflect.Int32:
		fmt.Fprintf(writer, " = %d", reflectValue.Int())
	default:
		fmt.Fprintf(writer, "unknown %d", reflectValue.Kind())
	}
}

func Tree(expr interface{}, writer io.Writer) {
	subType := reflect.ValueOf(expr)
	reflectExpression := subType.Elem()
	tree(reflectExpression, 0, writer)
}
