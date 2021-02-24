/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"sort"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func sortedTypeAtomKeys(m map[string]*ModuleDefinition) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortedTypes(m map[string]dtype.Type) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	return keys
}

func sortedExpressionKeys(m map[string]*ModuleDefinition) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

type ModuleReferenceDefinitions struct {
	referencedExpressions map[string]*ModuleDefinition

	ownedByModule *Module
}

func NewModuleReferenceDefinitions(ownedByModule *Module) *ModuleReferenceDefinitions {
	if ownedByModule == nil {
		panic("sorry, all definitions must be owned by a module")
	}
	return &ModuleReferenceDefinitions{
		ownedByModule:         ownedByModule,
		referencedExpressions: make(map[string]*ModuleDefinition),
	}
}

func (d *ModuleReferenceDefinitions) ReferencedDefinitions() []*ModuleDefinition {
	var all []*ModuleDefinition

	keys := sortedExpressionKeys(d.referencedExpressions)
	for _, key := range keys {
		all = append(all, d.referencedExpressions[key])
	}
	return all
}

func (d *ModuleReferenceDefinitions) AddDefinitions(definitions []*ModuleDefinition) error {
	for _, def := range definitions {
		if def == nil {
			panic("not allowed with empty definitions")
		}
		addErr := d.AddDefinition(def.Identifier(), def)
		if addErr != nil {
			return addErr
		}
	}
	return nil
}

func (d *ModuleReferenceDefinitions) AddDefinitionsWithModulePrefix(definitions []*ModuleDefinition, relative dectype.PackageRelativeModuleName) error {
	for _, def := range definitions {
		completeName := relative.JoinLocalName(def.Identifier())
		addErr := d.AddDefinition(ast.NewVariableIdentifier(token.NewVariableSymbolToken(completeName, token.SourceFileReference{}, 0)), def)
		if addErr != nil {
			return addErr
		}
	}
	return nil
}

func (d *ModuleReferenceDefinitions) FindDefinition(identifier *ast.VariableIdentifier) *ModuleDefinition {
	def, wasFound := d.referencedExpressions[identifier.Name()]
	if !wasFound {
		return nil
	}

	def.MarkAsReferenced()

	return def
}

func (d *ModuleReferenceDefinitions) AddDefinition(identifier *ast.VariableIdentifier, definition *ModuleDefinition) error {
	existingDeclare := d.FindDefinition(identifier)
	if existingDeclare != nil {
		return fmt.Errorf("sorry, '%v' already declared", existingDeclare)
	}
	d.referencedExpressions[identifier.Name()] = definition
	return nil
}

func (t *ModuleReferenceDefinitions) DebugString() string {
	s := "Module Definitions:\n"
	keys := sortedTypeAtomKeys(t.referencedExpressions)
	for _, key := range keys {
		definition := t.referencedExpressions[key]
		s += fmt.Sprintf(".. %v => %p %v\n", key, definition, definition)
	}

	return s
}

func (t *ModuleReferenceDefinitions) DebugOutput() {
	fmt.Println(t.DebugString())
}

func (t *ModuleReferenceDefinitions) ShortString() string {
	s := ""
	keys := sortedTypeAtomKeys(t.referencedExpressions)
	for _, key := range keys {
		definition := t.referencedExpressions[key]
		s += fmt.Sprintf(".. %v => %v\n", key, definition.ShortString())
	}
	return s
}

func (t *ModuleReferenceDefinitions) String() string {
	return t.ShortString()
}
