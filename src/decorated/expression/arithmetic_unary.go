/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	decshared "github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/token"
)

type ArithmeticUnaryOperatorType uint

const (
	ArithmeticUnaryMinus ArithmeticUnaryOperatorType = iota
)

type ArithmeticUnaryOperator struct {
	UnaryOperator
	operatorType ArithmeticUnaryOperatorType
	unary        *ast.UnaryExpression
}

func NewArithmeticUnaryOperator(unary *ast.UnaryExpression, left Expression, operatorType ArithmeticUnaryOperatorType) (*ArithmeticUnaryOperator, decshared.DecoratedError) {
	a := &ArithmeticUnaryOperator{operatorType: operatorType, unary: unary}
	a.UnaryOperator.left = left
	a.UnaryOperator.ExpressionNode.decoratedType = left.Type()
	return a, nil
}

func (a *ArithmeticUnaryOperator) OperatorType() ArithmeticUnaryOperatorType {
	return a.operatorType
}

func (a *ArithmeticUnaryOperator) Left() Expression {
	return a.left
}

func (a *ArithmeticUnaryOperator) String() string {
	return fmt.Sprintf("[unaryarithmetic %v left:%v]", a.operatorType, a.left)
}

func (a *ArithmeticUnaryOperator) FetchPositionLength() token.SourceFileReference {
	return a.unary.FetchPositionLength()
}
