/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"github.com/swamp/compiler/src/ast"
)

type ModuleDefinitionsCombine struct {
	internalDefinitions *ModuleDefinitions
	importDefinitions   *ModuleReferenceDefinitions
}

func NewModuleDefinitionsCombine(internalDefinitions *ModuleDefinitions,
	importDefinitions *ModuleReferenceDefinitions) *ModuleDefinitionsCombine {
	return &ModuleDefinitionsCombine{internalDefinitions: internalDefinitions, importDefinitions: importDefinitions}
}

func (d *ModuleDefinitionsCombine) FindDefinitionExpression(identifier *ast.VariableIdentifier) *ModuleDefinition {
	foundDef := d.internalDefinitions.FindDefinitionExpression(identifier)
	if foundDef == nil {
		return d.importDefinitions.FindDefinition(identifier)
	}

	return foundDef
}
