/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"bytes"
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type BinaryOperator struct {
	token    token.Token
	left     Expression
	operator token.OperatorToken
	right    Expression
}

func NewBinaryOperator(token token.Token, operator token.OperatorToken, left Expression,
	right Expression) *BinaryOperator {
	if left == nil {
		panic("left is nil")
	}

	if right == nil {
		panic("right is nil")
	}

	return &BinaryOperator{
		token:    token,
		operator: operator,
		left:     left,
		right:    right,
	}
}

func (i *BinaryOperator) Left() Expression {
	return i.left
}

func (i *BinaryOperator) OperatorType() token.Type {
	return i.token.Type()
}

func (i *BinaryOperator) Right() Expression {
	return i.right
}

func (i *BinaryOperator) Token() token.Token {
	return i.token
}

func (i *BinaryOperator) OperatorToken() token.OperatorToken {
	return i.operator
}

func (i *BinaryOperator) FetchPositionLength() token.SourceFileReference {
	return i.token.FetchPositionLength()
}

func (i *BinaryOperator) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.left.String())
	out.WriteString(" " + i.operator.String() + " ")
	out.WriteString(i.right.String())
	out.WriteString(")")

	return out.String()
}

func (i *BinaryOperator) DebugString() string {
	return fmt.Sprintf("[binaryop %v]", i.operator.String())
}
