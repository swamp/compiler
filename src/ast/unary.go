/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type UnaryExpression struct {
	token    token.Token
	left     Expression
	operator token.OperatorToken
}

func NewUnaryExpression(token token.Token, operator token.OperatorToken, left Expression) *UnaryExpression {
	return &UnaryExpression{
		token:    token,
		operator: operator,
		left:     left,
	}
}

func (i *UnaryExpression) Left() Expression {
	return i.left
}

func (i *UnaryExpression) OperatorType() token.Type {
	return i.token.Type()
}

func (i *UnaryExpression) OperatorToken() token.OperatorToken {
	return i.operator
}

func (i *UnaryExpression) FetchPositionLength() token.Range {
	return i.token.FetchPositionLength()
}

func (i *UnaryExpression) String() string {
	return fmt.Sprintf("(%v %v)", i.operator, i.left)
}

func (i *UnaryExpression) DebugString() string {
	return fmt.Sprintf("[infix %v]", i.operator.String())
}
