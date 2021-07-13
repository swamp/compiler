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
	importDefinitions   *ModuleImportedDefinitions
	importedModules     *ModuleImports
}

func NewModuleDefinitionsCombine(internalDefinitions *ModuleDefinitions,
	importDefinitions *ModuleImportedDefinitions, importedModules *ModuleImports) *ModuleDefinitionsCombine {
	return &ModuleDefinitionsCombine{internalDefinitions: internalDefinitions, importDefinitions: importDefinitions, importedModules: importedModules}
}

func (d *ModuleDefinitionsCombine) FindDefinitionExpression(identifier *ast.VariableIdentifier) ModuleDef {
	foundDef := d.internalDefinitions.FindDefinitionExpression(identifier)
	if foundDef == nil {
		importedDef := d.importDefinitions.FindDefinition(identifier)
		if importedDef == nil {
			return nil
		}
		return importedDef
	}

	foundDef.MarkAsReferenced()

	return foundDef
}

func (d *ModuleDefinitionsCombine) FindScopedDefinitionExpression(identifier *ast.VariableIdentifierScoped) ModuleDef {
	if d.importedModules == nil {
		log.Printf("it was scoped, but I dont have any imported modules %v", identifier)
	}
	foundModule := d.importedModules.FindModule(identifier.ModuleReference())
	if foundModule == nil {
		log.Printf("couldn't find module %v", identifier.ModuleReference())
		return nil
	}
	referencedModule := foundModule.referencedModule
	foundDef := referencedModule.exposedDefinitions.FindDefinition(identifier.AstVariableReference())
	if foundDef == nil {
		log.Printf("couldn't find definition in module %v\n%v\n", identifier, foundModule)
		return nil
	}

	foundModule.MarkAsReferenced()
	foundDef.MarkAsReferenced()

	return foundDef
}
