/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/token"
)

func tryConvertToBitwiseUnaryOperator(operatorType token.Type) (decorated.BitwiseUnaryOperatorType, bool) {
	switch operatorType {
	case token.OperatorBitwiseNot:
		return decorated.BitwiseUnaryNot, true
	}
	return 0, false
}

func tryConvertToLogicalUnaryOperator(operatorType token.Type) (decorated.LogicalUnaryOperatorType, bool) {
	switch operatorType {
	case token.OperatorNot:
		return decorated.LogicalUnaryNot, true
	}
	return 0, false
}

func decorateUnary(d DecorateStream, unary *ast.UnaryExpression, context *VariableContext) (decorated.DecoratedExpression, decshared.DecoratedError) {
	bitwiseUnaryOperatorType, isUnaryBitwise := tryConvertToBitwiseUnaryOperator(unary.OperatorType())
	if isUnaryBitwise {
		leftExpression, leftExpressionErr := DecorateExpression(d, unary.Left(), context)
		if leftExpressionErr != nil {
			return nil, leftExpressionErr
		}
		return decorated.NewBitwiseUnaryOperator(leftExpression, bitwiseUnaryOperatorType)
	}

	logicalUnaryOperatorType, isLogicalUnary := tryConvertToLogicalUnaryOperator(unary.OperatorType())
	if isLogicalUnary {
		leftExpression, leftExpressionErr := DecorateExpression(d, unary.Left(), context)
		if leftExpressionErr != nil {
			return nil, leftExpressionErr
		}
		return decorated.NewLogicalUnaryOperator(leftExpression, logicalUnaryOperatorType)
	}
	panic("unknown unary")
}
