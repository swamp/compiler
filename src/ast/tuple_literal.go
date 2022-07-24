/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type TupleLiteral struct {
	expressions    []Expression
	startListParen token.ParenToken
	endListParen   token.ParenToken
	inclusive      token.SourceFileReference
}

func NewTupleLiteral(startListParen token.ParenToken, endListParen token.ParenToken, expressions []Expression) *TupleLiteral {
	inclusive := token.MakeInclusiveSourceFileReference(startListParen.SourceFileReference, endListParen.SourceFileReference)
	return &TupleLiteral{startListParen: startListParen, endListParen: endListParen, expressions: expressions, inclusive: inclusive}
}

func (i *TupleLiteral) String() string {
	return fmt.Sprintf("[TupleLiteral %v]", i.expressions)
}

func (i *TupleLiteral) Expressions() []Expression {
	return i.expressions
}

func (i *TupleLiteral) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *TupleLiteral) StartParenToken() token.ParenToken {
	return i.startListParen
}

func (i *TupleLiteral) EndParenToken() token.ParenToken {
	return i.endListParen
}

func (i *TupleLiteral) DebugString() string {
	return fmt.Sprintf("[tuple-literal %v ]", i.expressions)
}
