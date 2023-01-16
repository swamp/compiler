/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type ListLiteral struct {
	expressions    []Expression
	startListParen token.ParenToken
	endListParen   token.ParenToken
	inclusive      token.SourceFileReference
}

func NewListLiteral(startListParen token.ParenToken, endListParen token.ParenToken, expressions []Expression) *ListLiteral {
	inclusive := token.MakeInclusiveSourceFileReference(startListParen.SourceFileReference, endListParen.SourceFileReference)
	return &ListLiteral{startListParen: startListParen, endListParen: endListParen, expressions: expressions, inclusive: inclusive}
}

func (i *ListLiteral) String() string {
	return fmt.Sprintf("[ListLiteral %v]", i.expressions)
}

func (i *ListLiteral) Expressions() []Expression {
	return i.expressions
}

func (i *ListLiteral) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *ListLiteral) StartParenToken() token.ParenToken {
	return i.startListParen
}

func (i *ListLiteral) EndParenToken() token.ParenToken {
	return i.endListParen
}

func (i *ListLiteral) DebugString() string {
	return fmt.Sprintf("[list-literal %v ]", i.expressions)
}
