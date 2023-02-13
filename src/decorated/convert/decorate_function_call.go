/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/concretize"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func getFunctionValueExpression(d DecorateStream, call *ast.FunctionCall, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	decoratedExpression, functionErr := DecorateExpression(d, call.FunctionExpression(), context)
	if functionErr != nil {
		return nil, functionErr
	}

	return decoratedExpression, nil
}

func determineEncounteredFunctionTypeAndArguments(d DecorateStream, call *ast.FunctionCall, functionValueExpressionFunctionType *dectype.FunctionAtom, encounteredCallParametersType *dectype.FunctionAtom, context *VariableContext) (*dectype.FunctionAtom, decshared.DecoratedError) {

	/*

		smashedFunctionType, smashErr := concretize.FunctionType(functionValueExpressionFunctionType, encounteredCallParametersType.FunctionParameterTypes(), context)
		if smashErr != nil {
			return nil, decorated.NewCouldNotSmashFunctions(call, functionValueExpressionFunctionType, encounteredCallParametersType, smashErr)
		}
	*/

	/* Smash is not enough, we need to add extra parameter types? */
	var completeCalledFunctionParameterTypes []dtype.Type
	completeCalledFunctionParameterTypes = append(completeCalledFunctionParameterTypes, encounteredCallParametersType.FunctionParameterTypes()...)
	extraParameters := functionValueExpressionFunctionType.FunctionParameterTypes()[len(encounteredCallParametersType.FunctionParameterTypes()):]

	completeCalledFunctionParameterTypes = append(completeCalledFunctionParameterTypes, extraParameters...)

	completeCalledFunctionType := dectype.NewFunctionAtom(nil, completeCalledFunctionParameterTypes)

	functionCompatibleErr := dectype.CompatibleTypes(functionValueExpressionFunctionType, completeCalledFunctionType)
	if functionCompatibleErr != nil {
		return nil, decorated.NewFunctionCallTypeMismatch(functionCompatibleErr, call, functionValueExpressionFunctionType, completeCalledFunctionType)
	}

	resolvedFunctionArguments, _ := completeCalledFunctionType.ParameterAndReturn()

	errorPosLength := call.FunctionExpression().FetchPositionLength()
	if encounteredCallParametersType.ParameterCount() > len(completeCalledFunctionParameterTypes) {
		return nil, decorated.NewExtraFunctionArguments(errorPosLength, resolvedFunctionArguments, encounteredCallParametersType.FunctionParameterTypes())
	}

	return completeCalledFunctionType, nil
}

func decorateFunctionCallInternal(d DecorateStream, call *ast.FunctionCall, functionValueExpression decorated.Expression, decoratedEncounteredArgumentExpressions []decorated.Expression, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	var encounteredArgumentTypes []dtype.Type
	for _, encounteredArgumentExpression := range decoratedEncounteredArgumentExpressions {
		encounteredArgumentTypes = append(encounteredArgumentTypes, encounteredArgumentExpression.Type())
	}
	encounteredFunctionCallType := dectype.NewFunctionAtom(nil, encounteredArgumentTypes)

	originalFunctionValueType := functionValueExpression.Type()
	unaliasedType := dectype.UnaliasWithResolveInvoker(originalFunctionValueType)
	functionValueExpressionFunctionType, wasFunction := unaliasedType.(*dectype.FunctionAtom)
	if !wasFunction {
		functionTypeReference, wasFunctionTypeReference := unaliasedType.(*dectype.FunctionTypeReference)
		if !wasFunctionTypeReference {
			localTypeContext, wasLocalTypeContext := unaliasedType.(*dectype.LocalTypeNameContext)
			if !wasLocalTypeContext {
				return nil, decorated.NewExpectedFunctionTypeForCall(functionValueExpression)
			}

			concreteFunctionType, concreteErr := concretize.ConcreteArguments(localTypeContext, encounteredArgumentTypes)
			if concreteErr != nil {
				return nil, concreteErr
			}

			functionValueExpressionFunctionType, _ = concreteFunctionType.(*dectype.FunctionAtom)

		} else {
			functionValueExpressionFunctionType = functionTypeReference.FunctionAtom()
		}
	}

	completeCalledFunctionType, determineErr := determineEncounteredFunctionTypeAndArguments(d, call, functionValueExpressionFunctionType, encounteredFunctionCallType, context)
	if determineErr != nil {
		return nil, determineErr
	}

	/* Extra check here. Is it neccessary?
	expectedArgumentTypes := completeCalledFunctionType.FunctionParameterTypes()
	for index, encounteredArgumentType := range encounteredArgumentTypes {
		expectedArgumentType := expectedArgumentTypes[index]
		compatibleErr := dectype.CompatibleTypes(expectedArgumentType, encounteredArgumentType)
		if compatibleErr != nil {
			return nil, decorated.NewFunctionArgumentTypeMismatch(call.FetchPositionLength(), nil, nil, expectedArgumentType, encounteredArgumentType, fmt.Errorf("%v %v", completeCalledFunctionType, compatibleErr))
		}
	}
	*/

	isCurrying := len(decoratedEncounteredArgumentExpressions) < completeCalledFunctionType.ParameterCount()-1
	if isCurrying {
		providedArgumentCount := len(decoratedEncounteredArgumentExpressions)
		allFunctionTypes := functionValueExpressionFunctionType.FunctionParameterTypes()
		curryFunctionTypes := allFunctionTypes[providedArgumentCount:]
		curryFunctionType := dectype.NewFunctionAtom(nil, curryFunctionTypes)

		return decorated.NewCurryFunction(call, curryFunctionType, functionValueExpression, decoratedEncounteredArgumentExpressions), nil
	}

	return decorated.NewFunctionCall(call, functionValueExpression, completeCalledFunctionType, decoratedEncounteredArgumentExpressions), nil
}

func decorateFunctionCall(d DecorateStream, call *ast.FunctionCall, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	functionValueExpression, functionReferenceErr := getFunctionValueExpression(d, call, context)
	if functionReferenceErr != nil {
		return nil, functionReferenceErr
	}

	var decoratedEncounteredArgumentExpressions []decorated.Expression
	var decoratedEncounteredArgumentTypes []dtype.Type
	for _, rawExpression := range call.Arguments() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, rawExpression, context)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}
		decoratedEncounteredArgumentExpressions = append(decoratedEncounteredArgumentExpressions, decoratedExpression)
		decoratedEncounteredArgumentTypes = append(decoratedEncounteredArgumentTypes, decoratedExpression.Type())
	}

	localName, wasLocal := functionValueExpression.Type().(*dectype.LocalTypeNameContext)
	if wasLocal {
		_, wasPointingToFunctionAtom := functionValueExpression.Type().Next().(*dectype.FunctionAtom)
		if !wasPointingToFunctionAtom {
			_, wasPointingToFunctionAtomRef := functionValueExpression.Type().Next().(*dectype.FunctionTypeReference)
			if !wasPointingToFunctionAtomRef {
				return nil, decorated.NewInternalError(fmt.Errorf("unknown function type %v", functionValueExpression.Type().Next()))
			}
		}
		_, concreteErr := concretize.ConcreteArguments(localName, decoratedEncounteredArgumentTypes)
		if concreteErr != nil {
			return nil, concreteErr
		}
	}

	return decorateFunctionCallInternal(d, call, functionValueExpression, decoratedEncounteredArgumentExpressions, context)
}
