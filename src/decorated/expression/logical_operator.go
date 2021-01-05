/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type LogicalOperatorType uint

const (
	LogicalAnd LogicalOperatorType = iota
	LogicalOr
)

type LogicalOperator struct {
	BinaryOperator
	operatorType LogicalOperatorType
}

func NewLogicalOperator(left DecoratedExpression, right DecoratedExpression, operatorType LogicalOperatorType, booleanType dtype.Type) (*LogicalOperator, decshared.DecoratedError) {
	a := &LogicalOperator{operatorType: operatorType}
	a.BinaryOperator.left = left
	a.BinaryOperator.right = right
	if left.Type() != booleanType {
		return nil, NewLogicalOperatorLeftMustBeBoolean(a, left)
	}

	if right.Type() != booleanType {
		return nil, NewLogicalOperatorRightMustBeBoolean(a, right)
	}

	a.BinaryOperator.DecoratedExpressionNode.decoratedType = left.Type()
	return a, nil
}

func (l *LogicalOperator) Left() DecoratedExpression {
	return l.BinaryOperator.left
}

func (l *LogicalOperator) Right() DecoratedExpression {
	return l.BinaryOperator.right
}

func (l *LogicalOperator) OperatorType() LogicalOperatorType {
	return l.operatorType
}

func (l *LogicalOperator) String() string {
	return fmt.Sprintf("[logical %v %v %v]", l.BinaryOperator.left, l.BinaryOperator.right, l.operatorType)
}

func (l *LogicalOperator) FetchPositionAndLength() token.PositionLength {
	return l.BinaryOperator.left.FetchPositionAndLength()
}
