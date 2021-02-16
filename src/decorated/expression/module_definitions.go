/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
)

type ModuleDefinitions struct {
	definitions   map[string]*ModuleDefinition
	ownedByModule *Module
}

func NewModuleDefinitions(ownedByModule *Module) *ModuleDefinitions {
	if ownedByModule == nil {
		panic("sorry, all definitions must be owned by a module")
	}
	return &ModuleDefinitions{ownedByModule: ownedByModule,
		definitions: make(map[string]*ModuleDefinition)}
}

func (d *ModuleDefinitions) Definitions() []*ModuleDefinition {
	var keys []*ModuleDefinition
	for _, exprKey := range sortedExpressionKeys(d.definitions) {
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
	return expressionDef
}

func (d *ModuleDefinitions) AddDecoratedExpression(identifier *ast.VariableIdentifier, expr DecoratedExpression) error {
	existingDeclare := d.FindDefinitionExpression(identifier)
	if existingDeclare != nil {
		return fmt.Errorf("sorry, '%v' already declared", existingDeclare)
	}

	def := NewModuleDefinition(d, identifier, expr)
	d.definitions[identifier.Name()] = def

	return nil
}

func (t *ModuleDefinitions) DebugString() string {
	s := "Module Definitions:\n"
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

	definitionKeys := sortedExpressionKeys(t.definitions)
	for _, expressionKey := range definitionKeys {
		expression := t.definitions[expressionKey]
		s += fmt.Sprintf("%s\n", expression.ShortString())
	}
	return s
}

func (t *ModuleDefinitions) String() string {
	return t.ShortString()
}
