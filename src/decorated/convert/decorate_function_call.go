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
)

func getFunctionValueExpression(d DecorateStream, call *ast.FunctionCall, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	/* TODO: SUpport recur again
	functionExpressionIdentifier, wasIdentifier := call.FunctionExpression().(*ast.VariableIdentifier)

	if wasIdentifier && functionExpressionIdentifier.Name() == "recur" {
		fakeIdent := ast.NewVariableIdentifier(token.NewVariableSymbolToken("__self", token.SourceFileReference{}, 0))
		namedDef, resolveErr := context.ResolveVariable(fakeIdent)
		if resolveErr != nil {
			return nil, resolveErr
		}
		functionReference := decorated.NewFunctionReference(functionExpressionIdentifier, namedDef.(*decorated.FunctionExpression))
		return functionReference, nil
	}

	*/

	decoratedExpression, functionErr := DecorateExpression(d, call.FunctionExpression(), context)
	if functionErr != nil {
		return nil, functionErr
	}

	return decoratedExpression, nil
}

func determineEncounteredFunctionTypeAndArguments(d DecorateStream, call *ast.FunctionCall, functionValueExpressionFunctionType *dectype.FunctionAtom, encounteredCallParametersType *dectype.FunctionAtom, context *VariableContext) (*dectype.FunctionAtom, decshared.DecoratedError) {
	/* Smash functions */
	smashedFunctionType, smashErr := dectype.SmashFunctions(functionValueExpressionFunctionType, encounteredCallParametersType)
	if smashErr != nil {
		return nil, decorated.NewCouldNotSmashFunctions(call, functionValueExpressionFunctionType, encounteredCallParametersType, smashErr)
	}
	/* end of smash functions */

	/* Smash is not enough, we need to add extra parameter types? */
	var completeCalledFunctionParameterTypes []dtype.Type
	completeCalledFunctionParameterTypes = append(completeCalledFunctionParameterTypes, encounteredCallParametersType.FunctionParameterTypes()...)
	extraParameters := smashedFunctionType.FunctionParameterTypes()[len(encounteredCallParametersType.FunctionParameterTypes()):]
	//nolint: gosimple
	for _, extraParameter := range extraParameters {
		// log.Printf("extra parameter:%v\n", extraParameter)
		completeCalledFunctionParameterTypes = append(completeCalledFunctionParameterTypes, extraParameter)
	}
	/* */
	completeCalledFunctionType := dectype.NewFunctionAtom(nil, completeCalledFunctionParameterTypes)

	functionCompatibleErr := dectype.CompatibleTypes(smashedFunctionType, completeCalledFunctionType)
	if functionCompatibleErr != nil {
		return nil, decorated.NewFunctionCallTypeMismatch(functionCompatibleErr, call, smashedFunctionType, completeCalledFunctionType)
	}

	resolvedFunctionArguments, _ := completeCalledFunctionType.ParameterAndReturn()

	errorPosLength := call.FunctionExpression().FetchPositionLength()
	if encounteredCallParametersType.ParameterCount() > len(completeCalledFunctionParameterTypes) {
		return nil, decorated.NewExtraFunctionArguments(errorPosLength, resolvedFunctionArguments, encounteredCallParametersType.FunctionParameterTypes())
	}

	return completeCalledFunctionType, nil
}

func decorateFunctionCallInternal(d DecorateStream, call *ast.FunctionCall, functionValueExpression decorated.Expression, decoratedEncounteredArgumentExpressions []decorated.Expression, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	functionValueExpressionType := functionValueExpression.Type()
	functionValueExpressionFunctionType, wasFunction := functionValueExpressionType.(*dectype.FunctionAtom)
	if !wasFunction {
		return nil, decorated.NewExpectedFunctionTypeForCall(functionValueExpression)
	}

	var encounteredArgumentTypes []dtype.Type
	for _, encounteredArgumentExpression := range decoratedEncounteredArgumentExpressions {
		encounteredArgumentTypes = append(encounteredArgumentTypes, encounteredArgumentExpression.Type())
	}
	encounteredFunctionCallType := dectype.NewFunctionAtom(nil, encounteredArgumentTypes)

	completeCalledFunctionType, determineErr := determineEncounteredFunctionTypeAndArguments(d, call, functionValueExpressionFunctionType, encounteredFunctionCallType, context)
	if determineErr != nil {
		return nil, determineErr
	}

	/*
			/* Extra check here. Is it neccessary?
		for index, encounteredArgumentType := range encounteredArgumentTypes {
			expectedArgumentType := resolvedFunctionArguments[index]
			compatibleErr := dectype.CompatibleTypes(expectedArgumentType, encounteredArgumentType)
			if compatibleErr != nil {
				return nil, decorated.NewFunctionArgumentTypeMismatch(errorPosLength, nil, nil, expectedArgumentType, encounteredArgumentType, fmt.Errorf("%v %v", smashedFunctionType, compatibleErr))
			}
		}
	*/

	isCurrying := len(decoratedEncounteredArgumentExpressions) < completeCalledFunctionType.ParameterCount()-1
	if isCurrying {
		providedArgumentCount := len(decoratedEncounteredArgumentExpressions)
		allFunctionTypes := completeCalledFunctionType.FunctionParameterTypes()
		curryFunctionTypes := allFunctionTypes[providedArgumentCount:]
		curryFunctionType := dectype.NewFunctionAtom(nil, curryFunctionTypes)
		return decorated.NewCurryFunction(curryFunctionType, functionValueExpression, decoratedEncounteredArgumentExpressions), nil
	}

	returnType := completeCalledFunctionType.ReturnType()

	return decorated.NewFunctionCall(call, functionValueExpression, returnType, decoratedEncounteredArgumentExpressions), nil
}

func decorateFunctionCall(d DecorateStream, call *ast.FunctionCall, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	functionValueExpression, functionReferenceErr := getFunctionValueExpression(d, call, context)
	if functionReferenceErr != nil {
		return nil, functionReferenceErr
	}

	var decoratedEncounteredArgumentExpressions []decorated.Expression
	for _, rawExpression := range call.Arguments() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, rawExpression, context)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}
		decoratedEncounteredArgumentExpressions = append(decoratedEncounteredArgumentExpressions, decoratedExpression)
	}

	return decorateFunctionCallInternal(d, call, functionValueExpression, decoratedEncounteredArgumentExpressions, context)
}
