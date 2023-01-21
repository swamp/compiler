/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type LogicalOperatorType uint

const (
	LogicalAnd LogicalOperatorType = iota
	LogicalOr
)

func LogicalOperatorToString(operatorType LogicalOperatorType) string {
	switch operatorType {
	case LogicalOr:
		return "or"
	case LogicalAnd:
		return "and"
	}

	panic("not possible")
}

type LogicalOperator struct {
	BinaryOperator
	operatorType LogicalOperatorType
}

func NewLogicalOperator(left Expression, right Expression, operatorType LogicalOperatorType, booleanType dtype.Type) (*LogicalOperator, decshared.DecoratedError) {
	a := &LogicalOperator{operatorType: operatorType}
	a.BinaryOperator.left = left
	a.BinaryOperator.right = right
	if err := dectype.CompatibleTypes(left.Type(), booleanType); err != nil {
		return nil, NewLogicalOperatorLeftMustBeBoolean(a, left, booleanType)
	}

	if err := dectype.CompatibleTypes(right.Type(), booleanType); err != nil {
		return nil, NewLogicalOperatorRightMustBeBoolean(a, right, booleanType)
	}

	a.BinaryOperator.ExpressionNode.decoratedType = left.Type()
	return a, nil
}

func (l *LogicalOperator) Left() Expression {
	return l.BinaryOperator.left
}

func (l *LogicalOperator) Right() Expression {
	return l.BinaryOperator.right
}

func (l *LogicalOperator) OperatorType() LogicalOperatorType {
	return l.operatorType
}

func (l *LogicalOperator) String() string {
	return fmt.Sprintf("[Logical %v %v %v]", l.BinaryOperator.left, LogicalOperatorToString(l.operatorType), l.BinaryOperator.right)
}

func (l *LogicalOperator) FetchPositionLength() token.SourceFileReference {
	return l.BinaryOperator.FetchPositionLength()
}
