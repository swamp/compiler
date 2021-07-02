/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type ModuleImportedDefinitions struct {
	importedDefinitions map[string]*ImportedDefinition

	ownedByModule *Module
}

func NewModuleImportedDefinitions(ownedByModule *Module) *ModuleImportedDefinitions {
	if ownedByModule == nil {
		panic("sorry, all localDefinitions must be owned by a module")
	}
	return &ModuleImportedDefinitions{
		ownedByModule:       ownedByModule,
		importedDefinitions: make(map[string]*ImportedDefinition),
	}
}

func (d *ModuleImportedDefinitions) ReferencedDefinitions() []*ImportedDefinition {
	var all []*ImportedDefinition

	keys := sortedExpressionKeys(d.importedDefinitions)
	for _, key := range keys {
		all = append(all, d.importedDefinitions[key])
	}
	return all
}

func (d *ModuleImportedDefinitions) AddDefinitions(definitions []*ImportedDefinition) error {
	for _, def := range definitions {
		if def == nil {
			panic("not allowed with empty localDefinitions")
		}
		addErr := d.AddDefinition(def.Identifier(), def)
		if addErr != nil {
			return addErr
		}
	}
	return nil
}

func (d *ModuleImportedDefinitions) AddDefinitionsWithModulePrefix(definitions []*ImportedDefinition, relative dectype.PackageRelativeModuleName) error {
	for _, def := range definitions {
		completeName := relative.JoinLocalName(def.Identifier())
		addErr := d.AddDefinition(ast.NewVariableIdentifier(token.NewVariableSymbolToken(completeName, token.SourceFileReference{}, 0)), def)
		if addErr != nil {
			return addErr
		}
	}
	return nil
}

func (d *ModuleImportedDefinitions) FindDefinition(identifier *ast.VariableIdentifier) *ImportedDefinition {
	def, wasFound := d.importedDefinitions[identifier.Name()]
	if !wasFound {
		return nil
	}

	def.MarkAsReferenced()

	return def
}

func (d *ModuleImportedDefinitions) AddDefinition(identifier *ast.VariableIdentifier, definition *ImportedDefinition) error {
	existingDeclare := d.FindDefinition(identifier)
	if existingDeclare != nil {
		return fmt.Errorf("sorry, '%v' already declared", existingDeclare)
	}
	d.importedDefinitions[identifier.Name()] = NewImportedDefinition(nil, identifier, definition)
	return nil
}

func (t *ModuleImportedDefinitions) DebugString() string {
	s := "Module Definitions:\n"
	keys := sortedTypeAtomKeys(t.importedDefinitions)
	for _, key := range keys {
		definition := t.importedDefinitions[key]
		s += fmt.Sprintf(".. %v => %p %v\n", key, definition, definition)
	}

	return s
}

func (t *ModuleImportedDefinitions) DebugOutput() {
	fmt.Println(t.DebugString())
}

func (t *ModuleImportedDefinitions) ShortString() string {
	s := ""
	keys := sortedTypeAtomKeys(t.importedDefinitions)
	for _, key := range keys {
		definition := t.importedDefinitions[key]
		s += fmt.Sprintf(".. %v => %v\n", key, definition.String())
	}
	return s
}

func (t *ModuleImportedDefinitions) String() string {
	return t.ShortString()
}
