/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	decshared "github.com/swamp/compiler/src/decorated/decshared"
)

type ArithmeticUnaryOperatorType uint

const (
	ArithmeticUnaryMinus ArithmeticUnaryOperatorType = iota
)

type ArithmeticUnaryOperator struct {
	UnaryOperator
	operatorType ArithmeticUnaryOperatorType
}

func NewArithmeticUnaryOperator(left DecoratedExpression, operatorType ArithmeticUnaryOperatorType) (*ArithmeticUnaryOperator, decshared.DecoratedError) {
	a := &ArithmeticUnaryOperator{operatorType: operatorType}
	a.UnaryOperator.left = left
	a.UnaryOperator.DecoratedExpressionNode.decoratedType = left.Type()
	return a, nil
}

func (a *ArithmeticUnaryOperator) OperatorType() ArithmeticUnaryOperatorType {
	return a.operatorType
}

func (a *ArithmeticUnaryOperator) Left() DecoratedExpression {
	return a.left
}

func (a *ArithmeticUnaryOperator) String() string {
	return fmt.Sprintf("[unaryarithmetic %v left:%v]", a.operatorType, a.left)
}
