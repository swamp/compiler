/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/debug"
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
	case token.OperatorRemainder:
		return decorated.ArithmeticRemainder, true
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
	case token.OperatorUpdate:
		return decorated.BitwiseOr, true
	case token.OperatorBitwiseXor:
		return decorated.BitwiseXor, true
	case token.OperatorBitwiseNot:
		return decorated.BitwiseNot, true
	case token.OperatorBitwiseShiftLeft:
		return decorated.BitwiseShiftLeft, true
	case token.OperatorBitwiseShiftRight:
		return decorated.BitwiseShiftRight, true
	}
	return 0, false
}

func tryConvertCastOperator(infix *ast.BinaryOperator, left decorated.Expression, right *decorated.AliasReference) (
	*decorated.CastOperator, decshared.DecoratedError,
) {
	a := decorated.NewCastOperator(left, right, infix)
	if err := dectype.CompatibleTypes(left.Type(), right.Type()); err != nil {
		return nil, decorated.NewUnmatchingBitwiseOperatorTypes(infix, left, nil)
	}
	return a, nil
}

func concretizeFunctionLike(functionLike decorated.Expression, singleArgumentType dtype.Type) (
	dtype.Type, decshared.DecoratedError,
) {
	var returnType dtype.Type

	switch t := functionLike.(type) {
	case *decorated.CurryFunction:
		functionAtom := t.OriginalFunctionType()
		if functionAtom == nil {
			panic(fmt.Errorf("can not convert to function type:%v", t.Type()))
		}
		returnType = functionAtom.ReturnType()
		indexToCheck := len(t.ArgumentsToSave())
		compareErr := dectype.CompatibleTypes(functionAtom.FunctionParameterTypes()[indexToCheck], singleArgumentType)
		if compareErr != nil {
			return nil, decorated.NewInternalError(compareErr)
		}
		log.Printf("CurryFunction %v %v %v", functionAtom, singleArgumentType, returnType)
	case *decorated.FunctionReference:
		originalRightFunctionAtom := t.FunctionValue().Type().(*dectype.FunctionAtom)
		returnType = originalRightFunctionAtom.ReturnType()
		indexToCheck := 0
		compareErr := dectype.CompatibleTypes(originalRightFunctionAtom.FunctionParameterTypes()[indexToCheck],
			singleArgumentType)
		if compareErr != nil {
			return nil, decorated.NewInternalError(compareErr)
		}
		log.Printf("right funcRef  %v %v %v", originalRightFunctionAtom, singleArgumentType, returnType)
	default:
		panic(fmt.Errorf("unknown function like decorated %T", functionLike))
	}

	return returnType, nil
}

func decoratePipeRight(d DecorateStream, infix *ast.BinaryOperator, context *VariableContext) (
	decorated.Expression, decshared.DecoratedError,
) {
	left := infix.Left()
	right := infix.Right()

	leftDecorated, leftErr := DecorateExpression(d, left, context)
	if leftErr != nil {
		return nil, leftErr
	}

	rightDecorated, rightErr := DecorateExpression(d, right, context)
	if rightErr != nil {
		return nil, rightErr
	}

	leftSideReturns := leftDecorated.Type()
	log.Printf("leftFunctionNode %v %s", left, debug.TreeString(leftSideReturns))

	resultingReturnType, concreteErr := concretizeFunctionLike(rightDecorated, leftSideReturns)
	if concreteErr != nil {
		return nil, concreteErr
	}

	return decorated.NewPipeRightOperator(leftDecorated, rightDecorated, resultingReturnType), nil
}

func decoratePipeLeft(d DecorateStream, infix *ast.BinaryOperator, context *VariableContext) (
	decorated.Expression, decshared.DecoratedError,
) {
	left := infix.Left()
	right := infix.Right()

	leftDecorated, leftErr := DecorateExpression(d, left, context)
	if leftErr != nil {
		return nil, leftErr
	}

	rightDecorated, rightErr := DecorateExpression(d, right, context)
	if rightErr != nil {
		return nil, rightErr
	}

	log.Printf("pipeLeft\nleft:\n %s\nright:\n %s", debug.TreeString(leftDecorated), debug.TreeString(rightDecorated))

	rightSideReturns := rightDecorated.Type()
	//log.Printf("pipeLeft %v %s", right, debug.TreeString(rightSideReturns))

	resultingReturnType, concreteErr := concretizeFunctionLike(leftDecorated, rightSideReturns)
	if concreteErr != nil {
		return nil, concreteErr
	}

	return decorated.NewPipeLeftOperator(leftDecorated, rightDecorated, resultingReturnType), nil
}

func decorateBinaryOperator(d DecorateStream, infix *ast.BinaryOperator, context *VariableContext) (
	decorated.Expression, decshared.DecoratedError,
) {
	if infix.OperatorType() == token.OperatorPipeLeft {
		return decoratePipeLeft(d, infix, context)
	} else if infix.OperatorType() == token.OperatorPipeRight {
		return decoratePipeRight(d, infix, context)
	} else {
		return decorateBinaryOperatorSameType(d, infix, context)
	}
}

func decorateBinaryOperatorSameType(d DecorateStream, infix *ast.BinaryOperator, context *VariableContext) (
	decorated.Expression, decshared.DecoratedError,
) {
	leftExpression, leftExpressionErr := DecorateExpression(d, infix.Left(), context)
	if leftExpressionErr != nil {
		return nil, leftExpressionErr
	}
	rightExpression, rightExpressionErr := DecorateExpression(d, infix.Right(), context)
	if rightExpressionErr != nil {
		return nil, rightExpressionErr
	}

	if infix.OperatorType() == token.OperatorCons {
		listType, err := dectype.GetListType(rightExpression.Type())
		if err != nil {
			return nil, decorated.NewUnExpectedListTypeForCons(infix, leftExpression, rightExpression)
		}
		listTypeType := listType.ParameterTypes()[0]
		compatibleErr := dectype.CompatibleTypes(leftExpression.Type(), listTypeType)
		if compatibleErr != nil {
			return nil, decorated.NewUnMatchingBinaryOperatorTypes(infix, leftExpression.Type(), rightExpression.Type())
		}
		return decorated.NewConsOperator(leftExpression, rightExpression, d.TypeReferenceMaker())
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
		opTypeUnreferenced := dectype.UnReference(opType)
		primitive, _ := opTypeUnreferenced.(*dectype.PrimitiveAtom)
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
		boolType := d.TypeReferenceMaker().FindBuiltInType("Bool", infix.FetchPositionLength())
		if boolType == nil {
			return nil, decorated.NewTypeNotFound("Bool")
		}
		return decorated.NewBooleanOperator(infix, leftExpression, rightExpression, booleanOperatorType, boolType)
	}

	logicalOperatorType, isLogical := tryConvertToLogicalOperator(infix.OperatorType())
	if isLogical {
		incompatibleErr := dectype.CompatibleTypes(leftExpression.Type(), rightExpression.Type())
		if incompatibleErr != nil {
			return nil, decorated.NewLogicalOperatorsMustBeBoolean(infix, leftExpression, rightExpression)
		}
		boolType := d.TypeReferenceMaker().FindBuiltInType("Bool", infix.FetchPositionLength())
		if boolType == nil {
			return nil, decorated.NewTypeNotFound("Bool")
		}
		return decorated.NewLogicalOperator(leftExpression, rightExpression, logicalOperatorType, boolType)
	}

	bitwiseOperatorType, isBitwise := tryConvertToBitwiseOperator(infix.OperatorType())
	if isBitwise {
		incompatibleErr := dectype.CompatibleTypes(leftExpression.Type(), rightExpression.Type())
		if incompatibleErr != nil {
			return nil, decorated.NewUnmatchingBitwiseOperatorTypes(infix, leftExpression, rightExpression)
		}
		return decorated.NewBitwiseOperator(infix, leftExpression, rightExpression, bitwiseOperatorType)
	}

	if infix.OperatorType() == token.Colon {
		aliasReference, _ := rightExpression.(*decorated.AliasReference)
		return tryConvertCastOperator(infix, leftExpression, aliasReference)
	}

	return nil, decorated.NewUnknownBinaryOperator(infix, leftExpression, rightExpression)
}
