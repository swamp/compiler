/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/swamp/compiler/src/ast"
)

type ModuleDefinitions struct {
	definitions        map[string]*ModuleDefinition
	orderedDefinitions []*ModuleDefinition
	ownedByModule      *Module
}

func NewModuleDefinitions(ownedByModule *Module) *ModuleDefinitions {
	if ownedByModule == nil {
		panic("sorry, all localDefinitions must be owned by a module")
	}
	return &ModuleDefinitions{
		ownedByModule: ownedByModule,
		definitions:   make(map[string]*ModuleDefinition),
	}
}

func (d *ModuleDefinitions) CopyFrom(other *ModuleDefinitions) error {
	for x, y := range other.orderedDefinitions {
		log.Printf("overwriting %v\n", x)
		d.definitions[y.FullyQualifiedVariableName().String()] = y
		d.orderedDefinitions = append(d.orderedDefinitions, y)
	}

	return nil
}

func (d *ModuleDefinitions) OwnedByModule() *Module {
	return d.ownedByModule
}

func (d *ModuleDefinitions) Definitions() []ModuleDef {
	var keys []ModuleDef
	for _, expr := range d.orderedDefinitions {
		keys = append(keys, expr)
	}

	return keys
}

func (d *ModuleDefinitions) FindDefinitionExpression(identifier *ast.VariableIdentifier) *ModuleDefinition {
	expressionDef, wasFound := d.definitions[identifier.Name()]
	if !wasFound {
		return nil
	}
	expressionDef.MarkAsReferenced()
	return expressionDef
}

func (d *ModuleDefinitions) AddDecoratedExpression(identifier *ast.VariableIdentifier, importModule *ImportedModule, expr Expression) error {
	existingDeclare := d.FindDefinitionExpression(identifier)
	if existingDeclare != nil {
		return fmt.Errorf("sorry, '%v' already declared", existingDeclare)
	}

	def := NewModuleDefinition(d, importModule, identifier, expr)
	d.definitions[identifier.Name()] = def
	d.orderedDefinitions = append(d.orderedDefinitions, def)

	return nil
}

func (d *ModuleDefinitions) AddEmptyExternalDefinition(identifier *ast.VariableIdentifier, importModule *ImportedModule) error {
	existingDeclare := d.FindDefinitionExpression(identifier)
	if existingDeclare != nil {
		return fmt.Errorf("sorry, '%v' already declared", existingDeclare)
	}

	def := NewModuleDefinition(d, importModule, identifier, nil)
	d.definitions[identifier.Name()] = def
	d.orderedDefinitions = append(d.orderedDefinitions, def)
	return nil
}

func (t *ModuleDefinitions) DebugString() string {
	s := "Module LocalDefinitions:\n"
	for _, definition := range t.definitions {
		s += fmt.Sprintf(".. %p %v\n", definition, definition)
	}

	return s
}

func (t *ModuleDefinitions) DebugOutput() {
	fmt.Println(t.DebugString())
}

func writeStructField(field reflect.StructField, writer io.Writer) {
	fmt.Fprintf(writer, ".%v", field.Name)
}

func tree(reflectValue reflect.Value, tab int, writer io.Writer) {
	if reflectValue.Kind() == reflect.Interface {
		log.Printf("type %s", reflectValue)
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
				fmt.Fprintf(writer, "(%s)", reflectValue.Type().Name())
			}
			if fieldWrittenCount > 0 {
				fmt.Fprintf(writer, "\n%s", tabs)
			}
			writeStructField(field, writer)
			fieldWrittenCount++
			tree(reflectField, tab+1, writer)
		}
		if fieldWrittenCount == 0 {
			panic(fmt.Errorf("not allowed '%s'", reflectValue.Type().Name()))
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

func (t *ModuleDefinitions) ShortString() string {
	var builder strings.Builder
	for _, expression := range t.orderedDefinitions {
		v := expression.Expression()
		subType := reflect.ValueOf(v)
		reflectExpression := subType.Elem()
		tree(reflectExpression, 0, &builder)
	}
	return builder.String()
}

func (t *ModuleDefinitions) String() string {
	return t.ShortString()
}
