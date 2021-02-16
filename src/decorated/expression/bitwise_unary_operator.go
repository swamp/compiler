/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/decshared"
)

type BitwiseUnaryOperatorType uint

const (
	BitwiseUnaryNot BitwiseUnaryOperatorType = iota
)

type BitwiseUnaryOperator struct {
	UnaryOperator
	operatorType BitwiseUnaryOperatorType
}

func NewBitwiseUnaryOperator(left DecoratedExpression, operatorType BitwiseUnaryOperatorType) (*BitwiseUnaryOperator, decshared.DecoratedError) {
	a := &BitwiseUnaryOperator{operatorType: operatorType}
	a.UnaryOperator.left = left
	a.UnaryOperator.DecoratedExpressionNode.decoratedType = left.Type()
	return a, nil
}

func (a *BitwiseUnaryOperator) OperatorType() BitwiseUnaryOperatorType {
	return a.operatorType
}

func (a *BitwiseUnaryOperator) Left() DecoratedExpression {
	return a.left
}

func (a *BitwiseUnaryOperator) String() string {
	return fmt.Sprintf("[unarybitwise %v left:%v]", a.operatorType, a.left)
}
