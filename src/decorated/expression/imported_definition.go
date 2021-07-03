/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
)

type ImportedDefinition struct {
	localIdentifier      *ast.VariableIdentifier
	createdBy            *ImportedModule
	referencedDefinition ModuleDef
	wasReferenced        bool
}

func NewImportedDefinition(createdBy *ImportedModule, identifier *ast.VariableIdentifier, referencedDefinition ModuleDef) *ImportedDefinition {
	if referencedDefinition == nil {
		panic("not allowed")
	}
	return &ImportedDefinition{
		createdBy: createdBy, localIdentifier: identifier, referencedDefinition: referencedDefinition,
	}
}

func (d *ImportedDefinition) Identifier() *ast.VariableIdentifier {
	return d.localIdentifier
}

func (d *ImportedDefinition) OwnedByModule() *Module {
	return d.referencedDefinition.OwnedByModule()
}

func (d *ImportedDefinition) CreatedBy() *ImportedModule {
	return d.createdBy
}

func (d *ImportedDefinition) FullyQualifiedVariableName() *FullyQualifiedVariableName {
	return d.referencedDefinition.OwnedByModule().FullyQualifiedName(d.localIdentifier)
}

func (d *ImportedDefinition) String() string {
	return fmt.Sprintf("[imdefx %v = %v]", d.localIdentifier, d.referencedDefinition)
}

func (d *ImportedDefinition) Definition() ModuleDef {
	return d.referencedDefinition
}

func (d *ImportedDefinition) Expression() Expression {
	return d.referencedDefinition.Expression()
}

func (d *ImportedDefinition) MarkAsReferenced() {
	d.wasReferenced = true
	d.referencedDefinition.MarkAsReferenced()
}

func (d *ImportedDefinition) IsInternal() bool {
	return d.OwnedByModule().IsInternal()
}

func (d *ImportedDefinition) WasReferenced() bool {
	return d.wasReferenced
}
