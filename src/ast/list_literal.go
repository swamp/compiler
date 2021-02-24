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
}

func NewListLiteral(startListParen token.ParenToken, expressions []Expression) *ListLiteral {
	return &ListLiteral{startListParen: startListParen, expressions: expressions}
}

func (i *ListLiteral) String() string {
	return fmt.Sprintf("[list-literal: %v]", i.expressions)
}

func (i *ListLiteral) Expressions() []Expression {
	return i.expressions
}

func (i *ListLiteral) FetchPositionLength() token.SourceFileReference {
	return i.startListParen.SourceFileReference
}

func (i *ListLiteral) ParenToken() token.ParenToken {
	return i.startListParen
}

func (i *ListLiteral) DebugString() string {
	return fmt.Sprintf("[list-literal %v ]", i.expressions)
}
