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
}

func NewArrayLiteral(startListParen token.ParenToken, expressions []Expression) *ArrayLiteral {
	return &ArrayLiteral{startListParen: startListParen, expressions: expressions}
}

func (i *ArrayLiteral) String() string {
	return fmt.Sprintf("[array-literal: %v]", i.expressions)
}

func (i *ArrayLiteral) Expressions() []Expression {
	return i.expressions
}

func (i *ArrayLiteral) PositionLength() token.PositionLength {
	return i.startListParen.FetchPositionLength()
}

func (i *ArrayLiteral) ParenToken() token.ParenToken {
	return i.startListParen
}

func (i *ArrayLiteral) DebugString() string {
	return fmt.Sprintf("[array-literal %v ]", i.expressions)
}
