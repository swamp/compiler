/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"log"
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/concretize"
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

func concretizeFunctionLike(functionLike decorated.Expression, argumentTypes []dtype.Type) (
	dtype.Type, decshared.DecoratedError,
) {
	var returnType dtype.Type

	argumentTypeCount := len(argumentTypes)
	lastArgumentType := argumentTypes[argumentTypeCount-1]
	if lastArgumentType == nil {
		log.Printf("argumentType is nil")
	}

	switch t := functionLike.(type) {
	case *decorated.CurryFunction:
		functionAtom := t.OriginalFunctionType()
		if functionAtom == nil {
			panic(fmt.Errorf("can not convert to function type:%v", t.Type()))
		}
		returnType = functionAtom.ReturnType()
		indexToCheck := len(t.ArgumentsToSave())
		log.Printf("CurryFunction %v %v %v", functionAtom, lastArgumentType, returnType)
		compareErr := dectype.CompatibleTypes(functionAtom.FunctionParameterTypes()[indexToCheck], lastArgumentType)
		if compareErr != nil {
			return nil, decorated.NewInternalError(compareErr)
		}
	case *decorated.FunctionReference:
		switch u := t.FunctionValue().Type().(type) {
		case *dectype.LocalTypeNameOnlyContext:
			ref := dectype.NewLocalTypeNameContextReference(nil, u)
			addedAnyAtEnd := append([]dtype.Type{}, argumentTypes...)
			addedAnyAtEnd = append(addedAnyAtEnd, dectype.NewAnyType())
			resolved, err := concretize.ConcretizeLocalTypeContextUsingArguments(ref,
				addedAnyAtEnd)
			if err != nil {
				return nil, err
			}

			atom, atomErr := resolved.Resolve()
			if atomErr != nil {
				return nil, decorated.NewInternalError(atomErr)
			}
			functionType, _ := atom.(*dectype.FunctionAtom)

			return functionType.ReturnType(), nil
		case *dectype.FunctionAtom:
			returnType = u.ReturnType()
			if len(argumentTypes) != u.ParameterCount()-1 {
				return nil, decorated.NewInternalError(fmt.Errorf("wrong number of arguments"))
			}
			for index, encounteredArgumentType := range argumentTypes {
				requiredArgumentType := u.FunctionParameterTypes()[index]
				compareErr := dectype.CompatibleTypes(requiredArgumentType,
					encounteredArgumentType)
				if compareErr != nil {
					return nil, decorated.NewFunctionArgumentTypeMismatch(t.FetchPositionLength(), nil, nil,
						requiredArgumentType, encounteredArgumentType, compareErr)
				}
			}
		default:
			panic(fmt.Errorf("unknown function ref decorated %T", t.FunctionValue().Type()))
		}
	default:
		panic(fmt.Errorf("unknown function like decorated %T", functionLike))
	}

	return returnType, nil
}

func generateFunctionCall(d DecorateStream, nonComplete ast.Expression, lastArgumentThatMightBeNonComplete dtype.Type,
	context *VariableContext) (*decorated.IncompleteFunctionCall,
	decshared.DecoratedError) {
	functionValueExpression := nonComplete

	var argumentTypes []dtype.Type

	var arguments []decorated.Expression

	lastArgumentType := lastArgumentThatMightBeNonComplete

	unaliasNonComplete := dectype.Unalias(lastArgumentThatMightBeNonComplete)
	switch t := unaliasNonComplete.(type) {
	case *dectype.FunctionAtom:
		log.Printf("function %T", t)
		lastArgumentType = t.FunctionParameterTypes()[len(t.FunctionParameterTypes())-1]
	}

	switch t := nonComplete.(type) {
	case *ast.FunctionCall:
		functionValueExpression = t.FunctionExpression()
		for _, astArg := range t.Arguments() {
			decArg, err := DecorateExpression(d, astArg, context)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, decArg)
			argumentTypes = append(argumentTypes, decArg.Type())
		}
	default:
		log.Printf("what is this %T %v", nonComplete, nonComplete)
	}

	decoratedFunctionValueExpression, functionErr := DecorateExpression(d, functionValueExpression, context)
	if functionErr != nil {
		return nil, functionErr
	}

	argumentTypes = append(argumentTypes, lastArgumentType)

	correctReturnType, err := concretizeFunctionLike(decoratedFunctionValueExpression, argumentTypes)
	if err != nil {
		return nil, err
	}

	return decorated.NewIncompleteFunctionCall(decoratedFunctionValueExpression, arguments, correctReturnType), nil
}

type Node interface {
	FetchPositionLength() token.SourceFileReference
	String() string
}

type Expression interface {
	Node
	Type() dtype.Type
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

	leftSideReturns := leftDecorated.Type()

	incompleteRightFunctionCall, genErr := generateFunctionCall(d, right, leftSideReturns, context)
	if genErr != nil {
		return nil, genErr
	}

	return decorated.NewPipeRightOperator(leftDecorated, incompleteRightFunctionCall), nil
}

func decoratePipeLeft(d DecorateStream, infix *ast.BinaryOperator, context *VariableContext) (
	decorated.Expression, decshared.DecoratedError,
) {
	left := infix.Left()
	right := infix.Right()

	tabs := strings.Repeat("..", G_depth)

	log.Printf("%s pipeLeft: decorating right %T", tabs, right)
	rightDecorated, rightErr := DecorateExpression(d, right, context)
	if rightErr != nil {
		return nil, rightErr
	}

	rightSideReturns := rightDecorated.Type()
	log.Printf("%s pipeLeft: rightReturns %v %s", tabs, right, debug.TreeString(rightSideReturns))
	incompleteLeftFunctionCall, genErr := generateFunctionCall(d, left, rightSideReturns, context)
	if genErr != nil {
		log.Printf("err: %v", genErr)
		return nil, genErr
	}

	log.Printf("%s pipeLeft: rightIncompleteReturns %v %s", tabs, incompleteLeftFunctionCall.Type(),
		debug.TreeString(incompleteLeftFunctionCall.Type()))

	log.Printf("%s pipeLeft:  rightFunctionNode %v %s", tabs, right, debug.TreeString(rightSideReturns))

	return decorated.NewPipeLeftOperator(incompleteLeftFunctionCall, rightDecorated), nil
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
