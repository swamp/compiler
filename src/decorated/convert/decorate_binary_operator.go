/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func tryConvertToArithmeticOperator(operatorType token.Type) (decorated.ArithmeticOperatorType, bool) {
	switch operatorType {
	case token.OperatorPlus:
		return decorated.ArithmeticPlus, true
	case token.OperatorAppend:
		return decorated.ArithmeticAppend, true
	case token.OperatorCons:
		return decorated.ArithmeticCons, true
	case token.OperatorMinus:
		return decorated.ArithmeticMinus, true
	case token.OperatorMultiply:
		return decorated.ArithmeticMultiply, true
	case token.OperatorDivide:
		return decorated.ArithmeticDivide, true
	}
	return 0, false
}

func tryConvertToBooleanOperator(operatorType token.Type) (decorated.BooleanOperatorType, bool) {
	switch operatorType {
	case token.OperatorEqual:
		return decorated.BooleanEqual, true
	case token.OperatorNotEqual:
		return decorated.BooleanNotEqual, true
	case token.OperatorLess:
		return decorated.BooleanLess, true
	case token.OperatorLessOrEqual:
		return decorated.BooleanLessOrEqual, true
	case token.OperatorGreater:
		return decorated.BooleanGreater, true
	case token.OperatorGreaterOrEqual:
		return decorated.BooleanGreaterOrEqual, true
	}
	return 0, false
}

func tryConvertToLogicalOperator(operatorType token.Type) (decorated.LogicalOperatorType, bool) {
	switch operatorType {
	case token.OperatorOr:
		return decorated.LogicalOr, true
	case token.OperatorAnd:
		return decorated.LogicalAnd, true
	}
	return 0, false
}

func tryConvertToBitwiseOperator(operatorType token.Type) (decorated.BitwiseOperatorType, bool) {
	switch operatorType {
	case token.OperatorBitwiseAnd:
		return decorated.BitwiseAnd, true
	case token.OperatorBitwiseOr:
		return decorated.BitwiseOr, true
	case token.OperatorBitwiseXor:
		return decorated.BitwiseXor, true
	case token.OperatorBitwiseNot:
		return decorated.BitwiseNot, true
	}
	return 0, false
}

func decorateBinaryOperator(d DecorateStream, infix *ast.BinaryOperator, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	leftExpression, leftExpressionErr := DecorateExpression(d, infix.Left(), context)
	if leftExpressionErr != nil {
		return nil, leftExpressionErr
	}
	rightExpression, rightExpressionErr := DecorateExpression(d, infix.Right(), context)
	if rightExpressionErr != nil {
		return nil, rightExpressionErr
	}

	if infix.OperatorType() == token.OperatorCons {
		return decorated.NewConsOperator(leftExpression, rightExpression)
	}

	compatibleErr := dectype.CompatibleTypes(leftExpression.Type(), rightExpression.Type())
	if compatibleErr != nil {
		return nil, decorated.NewUnMatchingBinaryOperatorTypes(infix, leftExpression.Type(), rightExpression.Type())
	}

	arithmeticOperatorType, worked := tryConvertToArithmeticOperator(infix.OperatorType())
	if worked {
		compatibleErr := dectype.CompatibleTypes(leftExpression.Type(), rightExpression.Type())
		if compatibleErr != nil {
			return nil, decorated.NewUnMatchingArithmeticOperatorTypes(infix, leftExpression, rightExpression)
		}
		opType := leftExpression.Type()
		primitive, _ := opType.(*dectype.PrimitiveAtom)
		if primitive != nil {
			if primitive.AtomName() == "Fixed" {
				if arithmeticOperatorType == decorated.ArithmeticMultiply {
					arithmeticOperatorType = decorated.ArithmeticFixedMultiply
				} else if arithmeticOperatorType == decorated.ArithmeticDivide {
					arithmeticOperatorType = decorated.ArithmeticFixedDivide
				}
			}
		}
		return decorated.NewArithmeticOperator(infix, leftExpression, rightExpression, arithmeticOperatorType)
	}

	booleanOperatorType, isBoolean := tryConvertToBooleanOperator(infix.OperatorType())
	if isBoolean {
		incompatibleErr := dectype.CompatibleTypes(leftExpression.Type(), rightExpression.Type())
		if incompatibleErr != nil {
			return nil, decorated.NewUnMatchingBooleanOperatorTypes(infix, leftExpression, rightExpression)
		}
		boolType := d.TypeRepo().FindTypeFromName("Bool")
		if boolType == nil {
			return nil, decorated.NewTypeNotFound("Bool")
		}
		return decorated.NewBooleanOperator(infix, leftExpression, rightExpression, booleanOperatorType, boolType.(dtype.Type))
	}

	logicalOperatorType, isLogical := tryConvertToLogicalOperator(infix.OperatorType())
	if isLogical {
		incompatibleErr := dectype.CompatibleTypes(leftExpression.Type(), rightExpression.Type())
		if incompatibleErr != nil {
			return nil, decorated.NewLogicalOperatorsMustBeBoolean(infix, leftExpression, rightExpression)
		}
		boolType := d.TypeRepo().FindTypeFromName("Bool")
		if boolType == nil {
			return nil, decorated.NewTypeNotFound("Bool")
		}
		return decorated.NewLogicalOperator(leftExpression, rightExpression, logicalOperatorType, boolType.(dtype.Type))
	}

	bitwiseOperatorType, isBitwise := tryConvertToBitwiseOperator(infix.OperatorType())
	if isBitwise {
		incompatibleErr := dectype.CompatibleTypes(leftExpression.Type(), rightExpression.Type())
		if incompatibleErr != nil {
			return nil, decorated.NewUnmatchingBitwiseOperatorTypes(infix, leftExpression, rightExpression)
		}
		return decorated.NewBitwiseOperator(infix, leftExpression, rightExpression, bitwiseOperatorType)
	}

	return nil, decorated.NewUnknownBinaryOperator(infix, leftExpression, rightExpression)
}
