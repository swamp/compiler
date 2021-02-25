/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/token"
)

type BitwiseUnaryOperatorType uint

const (
	BitwiseUnaryNot BitwiseUnaryOperatorType = iota
)

type BitwiseUnaryOperator struct {
	UnaryOperator
	operatorType BitwiseUnaryOperatorType
	unary        *ast.UnaryExpression
}

func NewBitwiseUnaryOperator(unary *ast.UnaryExpression, left Expression, operatorType BitwiseUnaryOperatorType) (*BitwiseUnaryOperator, decshared.DecoratedError) {
	a := &BitwiseUnaryOperator{operatorType: operatorType}
	a.unary = unary
	a.UnaryOperator.left = left
	a.UnaryOperator.ExpressionNode.decoratedType = left.Type()
	return a, nil
}

func (a *BitwiseUnaryOperator) OperatorType() BitwiseUnaryOperatorType {
	return a.operatorType
}

func (a *BitwiseUnaryOperator) Left() Expression {
	return a.left
}

func (a *BitwiseUnaryOperator) String() string {
	return fmt.Sprintf("[unarybitwise %v left:%v]", a.operatorType, a.left)
}

func (a *BitwiseUnaryOperator) FetchPositionLength() token.SourceFileReference {
	return a.unary.FetchPositionLength()
}
