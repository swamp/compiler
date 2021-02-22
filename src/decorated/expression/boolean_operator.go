/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type BooleanOperatorType uint

const (
	BooleanEqual BooleanOperatorType = iota
	BooleanNotEqual
	BooleanLess
	BooleanLessOrEqual
	BooleanGreater
	BooleanGreaterOrEqual
)

type BooleanOperator struct {
	BinaryOperator
	operatorType BooleanOperatorType
}

func NewBooleanOperator(infix *ast.BinaryOperator, left DecoratedExpression, right DecoratedExpression, operatorType BooleanOperatorType, booleanType dtype.Type) (*BooleanOperator, decshared.DecoratedError) {
	a := &BooleanOperator{operatorType: operatorType}
	a.BinaryOperator.left = left
	a.BinaryOperator.right = right
	if err := dectype.CompatibleTypes(left.Type(), right.Type()); err != nil {
		return nil, NewUnMatchingBooleanOperatorTypes(infix, left, right)
	}
	a.BinaryOperator.DecoratedExpressionNode.decoratedType = booleanType
	return a, nil
}

func booleanOperatorToString(t BooleanOperatorType) string {
	switch t {
	case BooleanEqual:
		return "EQ"
	case BooleanNotEqual:
		return "NEQ"
	case BooleanLess:
		return "LESS"
	case BooleanLessOrEqual:
		return "LE"
	case BooleanGreater:
		return "GR"
	case BooleanGreaterOrEqual:
		return "GRE"
	}

	panic("unknown boolean type")
}

func (l *BooleanOperator) String() string {
	return fmt.Sprintf("(boolop %v %v %v)", l.BinaryOperator.left,
		booleanOperatorToString(l.operatorType), l.BinaryOperator.right)
}

func (l *BooleanOperator) Left() DecoratedExpression {
	return l.BinaryOperator.left
}

func (l *BooleanOperator) Right() DecoratedExpression {
	return l.BinaryOperator.right
}

func (l *BooleanOperator) OperatorType() BooleanOperatorType {
	return l.operatorType
}

func (l *BooleanOperator) FetchPositionLength() token.Range {
	return l.BinaryOperator.left.FetchPositionLength()
}
