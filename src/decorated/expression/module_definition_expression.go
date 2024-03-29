/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
)

type ModuleDef interface {
	MarkAsReferenced()
	Expression() Expression
	Identifier() *ast.VariableIdentifier
	FullyQualifiedVariableName() *FullyQualifiedPackageVariableName
	OwnedByModule() *Module
	CreatedBy() *ImportedModule
	WasReferenced() bool
	IsInternal() bool
	IsExternal() bool
	ShortString() string
}

type ModuleDefinition struct {
	localIdentifier *ast.VariableIdentifier
	createdIn       *ModuleDefinitions
	createdBy       *ImportedModule
	expr            Expression
	wasReferenced   bool
}

func NewModuleDefinition(createdIn *ModuleDefinitions, createdBy *ImportedModule, identifier *ast.VariableIdentifier, expr Expression) *ModuleDefinition {
	return &ModuleDefinition{
		createdIn: createdIn, createdBy: createdBy, localIdentifier: identifier, expr: expr,
	}
}

func NewModuleDefinitionExternal(createdIn *ModuleDefinitions, identifier *ast.VariableIdentifier) *ModuleDefinition {
	return &ModuleDefinition{
		createdIn: createdIn, localIdentifier: identifier, expr: nil,
	}
}

func (d *ModuleDefinition) Identifier() *ast.VariableIdentifier {
	return d.localIdentifier
}

func (d *ModuleDefinition) CreatedBy() *ImportedModule {
	return d.createdBy
}

func (d *ModuleDefinition) ParentDefinitions() *ModuleDefinitions {
	return d.createdIn
}

func (d *ModuleDefinition) OwnedByModule() *Module {
	if d.createdIn == nil {
		panic("unknown ownedbymodule")
	}
	return d.createdIn.OwnedByModule()
}

func (d *ModuleDefinition) FullyQualifiedVariableName() *FullyQualifiedPackageVariableName {
	return d.createdIn.ownedByModule.FullyQualifiedName(d.localIdentifier)
}

func (d *ModuleDefinition) String() string {
	return fmt.Sprintf("[ModuleDef %v = %v]", d.localIdentifier, d.expr)
}

func (d *ModuleDefinition) ShortString() string {
	return fmt.Sprintf("[ModuleDef %v = %v]", d.localIdentifier, d.expr.Type())
}

func (d *ModuleDefinition) Expression() Expression {
	return d.expr
}

func (d *ModuleDefinition) IsInternal() bool {
	return d.OwnedByModule().IsInternal()
}

func (d *ModuleDefinition) IsExternal() bool {
	return d.expr == nil
}

func (d *ModuleDefinition) MarkAsReferenced() {
	d.wasReferenced = true
}

func (d *ModuleDefinition) WasReferenced() bool {
	return d.wasReferenced
}
