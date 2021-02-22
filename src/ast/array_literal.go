/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type ArrayLiteral struct {
	expressions    []Expression
	startListParen token.ParenToken
	endListParen   token.ParenToken
	inclusive      token.SourceFileReference
}

func NewArrayLiteral(startListParen token.ParenToken, endListParen token.ParenToken, expressions []Expression) *ArrayLiteral {
	inclusive := token.MakeInclusiveSourceFileReference(startListParen.SourceFileReference, endListParen.SourceFileReference)
	return &ArrayLiteral{startListParen: startListParen, expressions: expressions, endListParen: endListParen, inclusive: inclusive}
}

func (i *ArrayLiteral) String() string {
	return fmt.Sprintf("[array-literal: %v]", i.expressions)
}

func (i *ArrayLiteral) Expressions() []Expression {
	return i.expressions
}

func (i *ArrayLiteral) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *ArrayLiteral) StartParenToken() token.ParenToken {
	return i.startListParen
}

func (i *ArrayLiteral) EndParenToken() token.ParenToken {
	return i.endListParen
}

func (i *ArrayLiteral) DebugString() string {
	return fmt.Sprintf("[array-literal %v ]", i.expressions)
}
