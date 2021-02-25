/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type ArithmeticOperatorType uint

const (
	ArithmeticPlus ArithmeticOperatorType = iota
	ArithmeticMinus
	ArithmeticMultiply
	ArithmeticDivide
	ArithmeticAppend
	ArithmeticCons
	ArithmeticFixedMultiply
	ArithmeticFixedDivide
)

type ArithmeticOperator struct {
	BinaryOperator
	operatorType ArithmeticOperatorType
	infix        *ast.BinaryOperator
}

func NewArithmeticOperator(infix *ast.BinaryOperator, left Expression, right Expression, operatorType ArithmeticOperatorType) (*ArithmeticOperator, decshared.DecoratedError) {
	a := &ArithmeticOperator{operatorType: operatorType, infix: infix}
	a.BinaryOperator.left = left
	a.BinaryOperator.right = right
	if err := dectype.CompatibleTypes(left.Type(), right.Type()); err != nil {
		return nil, NewUnMatchingArithmeticOperatorTypes(infix, left, right)
	}
	a.BinaryOperator.ExpressionNode.decoratedType = left.Type()
	return a, nil
}

func (a *ArithmeticOperator) OperatorType() ArithmeticOperatorType {
	return a.operatorType
}

func (a *ArithmeticOperator) Left() Expression {
	return a.left
}

func (a *ArithmeticOperator) Right() Expression {
	return a.right
}

func arithmeticOperatorToString(t ArithmeticOperatorType) string {
	switch t {
	case ArithmeticPlus:
		return "PLUS"
	case ArithmeticMinus:
		return "MINUS"
	case ArithmeticMultiply:
		return "MULTIPLY"
	case ArithmeticDivide:
		return "DIVIDE"
	case ArithmeticAppend:
		return "APPEND"
	case ArithmeticCons:
		return "CONS"
	case ArithmeticFixedMultiply:
		return "FMULTIPLY"
	case ArithmeticFixedDivide:
		return "FDIVIDE"
	}
	panic("unknown arithmetic type")
}

func (a *ArithmeticOperator) String() string {
	return fmt.Sprintf("(arithmetic %v %v %v)", a.left, arithmeticOperatorToString(a.operatorType), a.right)
}

func (a *ArithmeticOperator) FetchPositionLength() token.SourceFileReference {
	return a.Left().FetchPositionLength()
}
