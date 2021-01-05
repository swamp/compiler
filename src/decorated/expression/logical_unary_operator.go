/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	decshared "github.com/swamp/compiler/src/decorated/decshared"
)

type LogicalUnaryOperatorType uint

const (
	LogicalUnaryNot LogicalUnaryOperatorType = iota
)

type LogicalUnaryOperator struct {
	UnaryOperator
	operatorType LogicalUnaryOperatorType
}

func NewLogicalUnaryOperator(left DecoratedExpression, operatorType LogicalUnaryOperatorType) (*LogicalUnaryOperator, decshared.DecoratedError) {
	a := &LogicalUnaryOperator{operatorType: operatorType}
	a.UnaryOperator.left = left
	a.UnaryOperator.DecoratedExpressionNode.decoratedType = left.Type()
	return a, nil
}

func (a *LogicalUnaryOperator) OperatorType() LogicalUnaryOperatorType {
	return a.operatorType
}

func (a *LogicalUnaryOperator) Left() DecoratedExpression {
	return a.left
}

func (a *LogicalUnaryOperator) String() string {
	return fmt.Sprintf("[unarylogical %v left:%v]", a.operatorType, a.left)
}
