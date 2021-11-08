/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/ast"
)

type ModuleDefinitions struct {
	definitions   map[string]*ModuleDefinition
	ownedByModule *Module
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
	for x, y := range other.definitions {
		log.Printf("overwriting %v\n", x)
		d.definitions[x] = y
	}

	return nil
}

func (d *ModuleDefinitions) OwnedByModule() *Module {
	return d.ownedByModule
}

func (d *ModuleDefinitions) Definitions() []ModuleDef {
	var keys []ModuleDef
	for _, exprKey := range sortedExpressionKeysDefinition(d.definitions) {
		expr := d.definitions[exprKey]
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

func (d *ModuleDefinitions) AddDecoratedExpression(identifier *ast.VariableIdentifier, expr Expression) error {
	existingDeclare := d.FindDefinitionExpression(identifier)
	if existingDeclare != nil {
		return fmt.Errorf("sorry, '%v' already declared", existingDeclare)
	}

	def := NewModuleDefinition(d, identifier, expr)
	d.definitions[identifier.Name()] = def

	return nil
}

func (d *ModuleDefinitions) AddEmptyExternalDefinition(identifier *ast.VariableIdentifier) error {
	existingDeclare := d.FindDefinitionExpression(identifier)
	if existingDeclare != nil {
		return fmt.Errorf("sorry, '%v' already declared", existingDeclare)
	}

	def := NewModuleDefinition(d, identifier, nil)
	d.definitions[identifier.Name()] = def

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

func (t *ModuleDefinitions) ShortString() string {
	s := ""

	definitionKeys := sortedExpressionKeysDefinition(t.definitions)
	for _, expressionKey := range definitionKeys {
		expression := t.definitions[expressionKey]
		s += fmt.Sprintf("%s\n", expression.String())
	}
	return s
}

func (t *ModuleDefinitions) String() string {
	return t.ShortString()
}
