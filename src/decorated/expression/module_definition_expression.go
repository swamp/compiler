/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"github.com/swamp/compiler/src/ast"
)

type ModuleDefinition struct {
	localIdentifier *ast.VariableIdentifier
	createdIn       *ModuleDefinitions
	expr            DecoratedExpression
}

func NewModuleDefinition(createdIn *ModuleDefinitions, identifier *ast.VariableIdentifier, expr DecoratedExpression) *ModuleDefinition {
	return &ModuleDefinition{createdIn: createdIn, localIdentifier: identifier, expr: expr}
}
func (d *ModuleDefinition) Identifier() *ast.VariableIdentifier {
	return d.localIdentifier
}

func (d *ModuleDefinition) FullyQualifiedVariableName() *FullyQualifiedVariableName {
	return d.createdIn.ownedByModule.FullyQualifiedName(d.localIdentifier)
}

func (d *ModuleDefinition) String() string {
	return fmt.Sprintf("[mdefx %v = %v]", d.localIdentifier, d.expr)
}

func (d *ModuleDefinition) ShortString() string {
	return fmt.Sprintf("%v = %v", d.localIdentifier.Name(), d.expr)
}

func (d *ModuleDefinition) Expression() DecoratedExpression {
	return d.expr
}
