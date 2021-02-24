/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type NamedDecoratedExpression struct {
	fullyQualifiedName string
	expression         DecoratedExpression
	wasReferenced      bool
	mDef               *ModuleDefinition
}

func NewNamedDecoratedExpression(fullyQualifiedName string, mDef *ModuleDefinition,
	expression DecoratedExpression) *NamedDecoratedExpression {
	if fullyQualifiedName == "" {
		panic("must have qualified name")
	}

	if expression == nil {
		panic("must have a valid expression")
	}

	return &NamedDecoratedExpression{
		fullyQualifiedName: fullyQualifiedName, mDef: mDef,
		expression: expression, wasReferenced: false,
	}
}

func (n *NamedDecoratedExpression) FullyQualifiedName() string {
	return n.fullyQualifiedName
}

func (n *NamedDecoratedExpression) String() string {
	return fmt.Sprintf("[decoratedexpression %v %v]", n.fullyQualifiedName, n.expression)
}

func (n *NamedDecoratedExpression) SetReferenced() {
	n.wasReferenced = true
}

func (n *NamedDecoratedExpression) WasReferenced() bool {
	return n.wasReferenced
}

func (n *NamedDecoratedExpression) Expression() DecoratedExpression {
	return n.expression
}

func (n *NamedDecoratedExpression) ModuleDefinition() *ModuleDefinition {
	return n.mDef
}

func (n *NamedDecoratedExpression) Type() dtype.Type {
	return n.expression.Type()
}

func (n *NamedDecoratedExpression) FetchPositionLength() token.SourceFileReference {
	return token.SourceFileReference{}
}
