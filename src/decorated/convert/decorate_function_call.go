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

func DecorateFunctionValueForCall(symbolToken token.Range, resolvedFunction *dectype.FunctionAtom, encounteredArgumentTypes []dtype.Type) (bool, dtype.Type, *dectype.FunctionAtom, decshared.DecoratedError) {
	resolvedFunctionArguments, _ := resolvedFunction.ParameterAndReturn()
	isCurrying := len(encounteredArgumentTypes) < len(resolvedFunctionArguments)

	if len(encounteredArgumentTypes) > len(resolvedFunctionArguments) {
		return false, nil, nil, decorated.NewExtraFunctionArguments(symbolToken, resolvedFunctionArguments, encounteredArgumentTypes)
	}
	for index, encounteredArgumentType := range encounteredArgumentTypes {
		expectedArgumentType := resolvedFunctionArguments[index]
		compatibleErr := dectype.CompatibleTypes(expectedArgumentType, encounteredArgumentType)
		if compatibleErr != nil {
			return false, nil, nil, decorated.NewFunctionArgumentTypeMismatch(symbolToken, nil, nil, expectedArgumentType, encounteredArgumentType, fmt.Errorf("%v %v", resolvedFunction, compatibleErr))
		}
	}

	if isCurrying {
		providedArgumentCount := len(encounteredArgumentTypes)
		allFunctionTypes := resolvedFunction.FunctionParameterTypes()
		curryFunctionTypes := allFunctionTypes[providedArgumentCount:]
		curryFunctionType := dectype.NewFunctionAtom(curryFunctionTypes)
		return isCurrying, curryFunctionType, curryFunctionType, nil
	}

	returnType := resolvedFunction.ReturnType()
	return isCurrying, returnType, resolvedFunction, nil
}

func decorateFunctionCall(d DecorateStream, call *ast.FunctionCall, context *VariableContext) (decorated.DecoratedExpression, decshared.DecoratedError) {
	var decoratedExpressions []decorated.DecoratedExpression

	for _, rawExpression := range call.Arguments() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, rawExpression, context)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}
		decoratedExpressions = append(decoratedExpressions, decoratedExpression)
	}

	var encounteredArgumentTypes []dtype.Type
	for _, encounteredArgumentExpression := range decoratedExpressions {
		encounteredArgumentTypes = append(encounteredArgumentTypes, encounteredArgumentExpression.Type())
	}

	var decoratedFunctionExpression decorated.DecoratedExpression

	beforeFakeIdentifier, wasIdentifier := call.FunctionExpression().(*ast.VariableIdentifier)

	if wasIdentifier && beforeFakeIdentifier.Name() == "recur" {
		fakeIdent := ast.NewVariableIdentifier(token.NewVariableSymbolToken("__self", nil, token.Range{}, 0))
		namedDef := context.ResolveVariable(fakeIdent)

		fakeFunctionName := ast.NewVariableIdentifier(token.NewVariableSymbolToken(namedDef.FullyQualifiedName(), nil, token.Range{}, 0))
		getVar := decorated.NewGetVariable(fakeFunctionName, namedDef)
		decoratedFunctionExpression = getVar
	} else {
		var functionErr decshared.DecoratedError
		decoratedFunctionExpression, functionErr = DecorateExpression(d, call.FunctionExpression(), context)
		if functionErr != nil {
			return nil, functionErr
		}
	}

	if decoratedFunctionExpression == nil {
		return nil, decorated.NewInternalError(fmt.Errorf("expression was not just a variable identifier %v %v", call.FunctionExpression(), call.FetchPositionLength()))
	}

	callFunctionType := dectype.NewFunctionAtom(encounteredArgumentTypes)
	hopefullyFunctionType := decoratedFunctionExpression.Type()
	functionTypeOriginal, wasFunction := hopefullyFunctionType.(*dectype.FunctionAtom)
	if !wasFunction {
		return nil, decorated.NewExpectedFunctionTypeForCall(decoratedFunctionExpression)
	}

	functionType, smashErr := dectype.SmashFunctions(functionTypeOriginal, callFunctionType)
	if smashErr != nil {
		return nil, decorated.NewCouldNotSmashFunctions(call, functionTypeOriginal, callFunctionType, smashErr)
	}
	fakeIdentifier, _ := decoratedFunctionExpression.(*decorated.GetVariableOrReferenceFunction)

	// fmt.Printf("\n\ncall %v\n", fakeIdentifier)
	var complete []dtype.Type
	complete = append(complete, callFunctionType.FunctionParameterTypes()...)
	extraParameters := functionType.FunctionParameterTypes()[len(callFunctionType.FunctionParameterTypes()):]

	//nolint: gosimple
	for _, extraParameter := range extraParameters {
		//	fmt.Printf("extra parameter:%v\n", extraParameter)
		complete = append(complete, extraParameter)
	}

	completeCalledFunction := dectype.NewFunctionAtom(complete)

	functionCompatibleErr := dectype.CompatibleTypes(functionType, completeCalledFunction)
	if functionCompatibleErr != nil {
		return nil, decorated.NewFunctionCallTypeMismatch(functionCompatibleErr, call, functionType, completeCalledFunction)
	}

	errorPosLength := call.FunctionExpression().FetchPositionLength()

	isCurrying, fn, _, err := DecorateFunctionValueForCall(errorPosLength, functionType, encounteredArgumentTypes)
	if err != nil {
		return nil, err
	}

	functionValueExpression := NewFakeExpression(fn)
	functionValueDecoratedExpression := decorated.NewNamedDecoratedExpression("x", nil, functionValueExpression)
	functionValueDecoratedExpression.SetReferenced()

	fakeVariable := ast.NewVariableIdentifier(token.NewVariableSymbolToken(fakeIdentifier.Identifier().Name(), nil, token.Range{}, 8))
	getVariableExpression := decorated.NewGetVariable(fakeVariable, functionValueDecoratedExpression)

	if isCurrying {
		return decorated.NewCurryFunction(getVariableExpression, decoratedExpressions), nil
	}
	returnType := functionType.ReturnType()

	return decorated.NewFunctionCall(getVariableExpression, returnType, decoratedExpressions), nil
}
