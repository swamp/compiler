/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"log"

	"github.com/swamp/compiler/src/ast"
)

type ModuleDefinitionsCombine struct {
	internalDefinitions *ModuleDefinitions
	importDefinitions   *ModuleReferenceDefinitions
	importedModules     *ModuleImports
}

func NewModuleDefinitionsCombine(internalDefinitions *ModuleDefinitions,
	importDefinitions *ModuleReferenceDefinitions, importedModules *ModuleImports) *ModuleDefinitionsCombine {
	return &ModuleDefinitionsCombine{internalDefinitions: internalDefinitions, importDefinitions: importDefinitions, importedModules: importedModules}
}

func (d *ModuleDefinitionsCombine) FindDefinitionExpression(identifier *ast.VariableIdentifier) *ModuleDefinition {
	foundDef := d.internalDefinitions.FindDefinitionExpression(identifier)
	if foundDef == nil {
		return d.importDefinitions.FindDefinition(identifier)
	}

	foundDef.MarkAsReferenced()

	return foundDef
}

func (d *ModuleDefinitionsCombine) FindScopedDefinitionExpression(identifier *ast.VariableIdentifierScoped) *ModuleDefinition {
	if d.importedModules == nil {
		log.Printf("it was scoped, but I dont have any imported modules %v", identifier)
	}
	foundModule := d.importedModules.FindModule(identifier.ModuleReference())
	if foundModule == nil {
		return nil
	}
	foundDef := foundModule.exposedDefinitions.FindDefinition(identifier.AstVariableReference())
	if foundDef == nil {
		return nil
	}

	foundDef.MarkAsReferenced()

	return foundDef
}
