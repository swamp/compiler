/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

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

/*
func parsePipeLeftExpression(p ParseStream, operatorToken token.OperatorToken, startIndentation int, precedence Precedence, left ast.FunctionExpression) (ast.FunctionExpression, parerr.ParseError) {
	_, spaceErr := p.eatOneSpace("space after pipe left")
	if spaceErr != nil {
		return nil, spaceErr
	}
	right, rightErr := p.parseExpressionNormal(startIndentation)
	if rightErr != nil {
		return nil, rightErr
	}

	leftCall, _ := left.(ast.FunctionCaller)
	if leftCall == nil {
		leftVar, _ := left.(*ast.VariableIdentifier)
		if leftVar == nil {
			return nil, parerr.NewLeftPartOfPipeMustBeFunctionCallError(operatorToken)
		}
		leftCall = ast.NewFunctionCall(leftVar, nil)
	}

	rightCall, _ := right.(ast.FunctionCaller)
	if rightCall == nil {
		return nil, parerr.NewRightPartOfPipeMustBeFunctionCallError(operatorToken)
	}

	args := leftCall.Arguments()
	args = append(args, rightCall)
	leftCall.OverwriteArguments(args)

	return leftCall, nil
}

*/

func defToFunctionReference(def *decorated.NamedDecoratedExpression, ident ast.ScopedOrNormalVariableIdentifier) *decorated.FunctionReference {
	lookupExpression := def.Expression()
	functionValue, _ := lookupExpression.(*decorated.FunctionValue)

	fromModule := def.ModuleDefinition().OwnedByModule()
	moduleRef := decorated.NewModuleReference(ast.NewModuleReference(fromModule.FullyQualifiedModuleName().Path().Parts()), fromModule)
	nameWithModuleRef := decorated.NewNamedDefinitionReference(moduleRef, ident)
	return decorated.NewFunctionReference(nameWithModuleRef, functionValue)
}

func decorateHalfOfAFunctionCall(d DecorateStream, left ast.Expression, context *VariableContext) (*ast.FunctionCall, decorated.Expression, []decorated.Expression, decshared.DecoratedError) {
	var arguments []decorated.Expression
	var functionExpression decorated.Expression
	var leftAstCall *ast.FunctionCall
	switch t := left.(type) {
	case *ast.FunctionCall:
		funcExpr, funcExprErr := DecorateExpression(d, t.FunctionExpression(), context)
		if funcExprErr != nil {
			return nil, nil, nil, funcExprErr
		}
		functionExpression = funcExpr
		for _, astArgument := range t.Arguments() {
			expr, exprErr := DecorateExpression(d, astArgument, context)
			if exprErr != nil {
				return nil, nil, nil, exprErr
			}
			arguments = append(arguments, expr)
		}
		leftAstCall = t
	case *ast.VariableIdentifier:
		def := context.FindNamedDecoratedExpression(t)
		if def == nil {
			return nil, nil, nil, decorated.NewInternalError(fmt.Errorf("couldn't find %v", t))
		}

		functionReference := defToFunctionReference(def, t)
		functionExpression = functionReference
		leftAstCall = ast.NewFunctionCall(functionReference, nil)
	case *ast.VariableIdentifierScoped:
		def := context.FindScopedNamedDecoratedExpression(t)
		if def == nil {
			return nil, nil, nil, decorated.NewInternalError(fmt.Errorf("couldn't find %v", t))
		}
		functionReference := defToFunctionReference(def, t)
		functionExpression = functionReference
		leftAstCall = ast.NewFunctionCall(functionReference, nil)
	}
	return leftAstCall, functionExpression, arguments, nil
}

func decoratePipeLeft(d DecorateStream, infix *ast.BinaryOperator, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	left := infix.Left()
	right := infix.Right()

	rightDecorated, rightErr := DecorateExpression(d, right, context)
	if rightErr != nil {
		return nil, rightErr
	}

	leftAstCall, functionExpression, arguments, halfErr := decorateHalfOfAFunctionCall(d, left, context)
	if halfErr != nil {
		return nil, halfErr
	}

	var allArguments []decorated.Expression
	allArguments = append(arguments, rightDecorated)

	fullLeftFunctionCall, functionCallErr := decorateFunctionCallInternal(d, leftAstCall, functionExpression, allArguments, context)
	if functionCallErr != nil {
		return nil, functionCallErr
	}

	calculatedFunctionCallType := functionExpression.Type().(*dectype.FunctionTypeReference).FunctionAtom()

	halfLeftSideFunctionCall := decorated.NewFunctionCall(leftAstCall, functionExpression, calculatedFunctionCallType, arguments)

	return decorated.NewPipeLeftOperator(halfLeftSideFunctionCall, rightDecorated, fullLeftFunctionCall), nil
}

func decoratePipeRight(d DecorateStream, infix *ast.BinaryOperator, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	left := infix.Left()
	right := infix.Right()

	leftDecorated, leftErr := DecorateExpression(d, left, context)
	if leftErr != nil {
		return nil, leftErr
	}

	rightAstCall, functionExpression, arguments, halfErr := decorateHalfOfAFunctionCall(d, right, context)
	if halfErr != nil {
		return nil, halfErr
	}

	var allArguments []decorated.Expression
	allArguments = append(arguments, leftDecorated)

	fullRightFunctionCall, functionCallErr := decorateFunctionCallInternal(d, rightAstCall, functionExpression, allArguments, context)
	if functionCallErr != nil {
		return nil, functionCallErr
	}

	functionCall := fullRightFunctionCall.(*decorated.FunctionCall)

	halfRightFunctionCall := decorated.NewFunctionCall(rightAstCall, functionCall, functionCall.CompleteCalledFunctionType(), arguments)

	return decorated.NewPipeRightOperator(leftDecorated, halfRightFunctionCall, fullRightFunctionCall), nil
}

func decorateBinaryOperator(d DecorateStream, infix *ast.BinaryOperator, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	if infix.OperatorType() == token.OperatorPipeLeft {
		return decoratePipeLeft(d, infix, context)
	} else if infix.OperatorType() == token.OperatorPipeRight {
		return decoratePipeRight(d, infix, context)
	} else {
		return decorateBinaryOperatorSameType(d, infix, context)
	}
}

func decorateBinaryOperatorSameType(d DecorateStream, infix *ast.BinaryOperator, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
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
		boolType := d.TypeRepo().FindBuiltInType("Bool")
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
		boolType := d.TypeRepo().FindBuiltInType("Bool")
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
