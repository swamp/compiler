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

type BitwiseOperatorType uint

const (
	BitwiseOr BitwiseOperatorType = iota
	BitwiseAnd
	BitwiseXor
	BitwiseShiftLeft
	BitwiseShiftRight
	BitwiseNot
)

type BitwiseOperator struct {
	BinaryOperator
	operatorType BitwiseOperatorType
}

func NewBitwiseOperator(infix *ast.BinaryOperator, left Expression, right Expression,
	operatorType BitwiseOperatorType) (*BitwiseOperator, decshared.DecoratedError) {
	a := &BitwiseOperator{operatorType: operatorType}
	a.BinaryOperator.left = left
	a.BinaryOperator.right = right
	if err := dectype.CompatibleTypes(left.Type(), right.Type()); err != nil {
		return nil, NewUnmatchingBitwiseOperatorTypes(infix, left, right)
	}
	a.BinaryOperator.ExpressionNode.decoratedType = left.Type()
	return a, nil
}

func (a *BitwiseOperator) OperatorType() BitwiseOperatorType {
	return a.operatorType
}

func (a *BitwiseOperator) Left() Expression {
	return a.left
}

func (a *BitwiseOperator) Right() Expression {
	return a.right
}

func (a *BitwiseOperator) String() string {
	return fmt.Sprintf("[bitwise %v left:%v right:%v]", a.operatorType, a.left, a.right)
}

func (a *BitwiseOperator) FetchPositionLength() token.SourceFileReference {
	return a.left.FetchPositionLength()
}
