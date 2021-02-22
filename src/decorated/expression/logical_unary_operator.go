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

type LogicalUnaryOperatorType uint

const (
	LogicalUnaryNot LogicalUnaryOperatorType = iota
)

type LogicalUnaryOperator struct {
	UnaryOperator
	operatorType LogicalUnaryOperatorType
	unary        *ast.UnaryExpression
}

func NewLogicalUnaryOperator(unary *ast.UnaryExpression, left Expression, operatorType LogicalUnaryOperatorType) (*LogicalUnaryOperator, decshared.DecoratedError) {
	a := &LogicalUnaryOperator{operatorType: operatorType, unary: unary}
	a.UnaryOperator.left = left
	a.UnaryOperator.ExpressionNode.decoratedType = left.Type()
	return a, nil
}

func (a *LogicalUnaryOperator) OperatorType() LogicalUnaryOperatorType {
	return a.operatorType
}

func (a *LogicalUnaryOperator) Left() Expression {
	return a.left
}

func (a *LogicalUnaryOperator) String() string {
	return fmt.Sprintf("[unarylogical %v left:%v]", a.operatorType, a.left)
}

func (a *LogicalUnaryOperator) FetchPositionLength() token.SourceFileReference {
	return a.unary.FetchPositionLength()
}
