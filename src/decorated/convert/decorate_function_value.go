/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func createVariableContextFromParameters(context *VariableContext, parameters []*decorated.FunctionParameterDefinition, forcedFunctionType *dectype.FunctionAtom, functionName *ast.VariableIdentifier) *VariableContext {
	newVariableContext := context.MakeVariableContext()

	for _, parameter := range parameters {
		namedDecoratedExpression := decorated.NewNamedDecoratedExpression("___", nil, parameter)
		newVariableContext.Add(parameter.Identifier(), namedDecoratedExpression)
	}


	self := ast.NewVariableIdentifier(token.NewVariableSymbolToken("__self", token.PositionLength{}, 0))
	selfDef := decorated.NewFunctionParameterDefinition(self, forcedFunctionType)
	namedDecoratedExpression := decorated.NewNamedDecoratedExpression(functionName.Name(), nil, selfDef)
	newVariableContext.Add(self, namedDecoratedExpression)

	return newVariableContext
}

func DecorateFunctionValue(d DecorateStream, potentialFunc *ast.FunctionValue, forcedFunctionType *dectype.FunctionAtom,
	functionName *ast.VariableIdentifier, context *VariableContext, comments []ast.LocalComment) (decorated.DecoratedExpression, decshared.DecoratedError) {
	if forcedFunctionType == nil {
		return nil, decorated.NewInternalError(fmt.Errorf("I have no forced function type %v", potentialFunc))
	}

	parameterTypes, expectedReturnType := forcedFunctionType.ParameterAndReturn()
	if len(parameterTypes) != len(potentialFunc.Parameters()) {
		return nil, decorated.NewWrongNumberOfArgumentsInFunctionValue(potentialFunc, forcedFunctionType, parameterTypes)
	}

	functionParameterTypes, _ := forcedFunctionType.ParameterAndReturn()
	identifiers := potentialFunc.Parameters()
	var parameters []*decorated.FunctionParameterDefinition
	for index, parameterType := range functionParameterTypes {
		identifier := identifiers[index]
		argDef := decorated.NewFunctionParameterDefinition(identifier, parameterType)
		parameters = append(parameters, argDef)
	}

	subVariableContext := createVariableContextFromParameters(context, parameters, forcedFunctionType, functionName)
	expression := potentialFunc.Expression()
	decoratedExpression, decoratedExpressionErr := DecorateExpression(d, expression, subVariableContext)
	if decoratedExpressionErr != nil {
		return nil, decoratedExpressionErr
	}
	decoratedExpressionType := decoratedExpression.Type()
	if decoratedExpressionType == nil {
		fmt.Printf("%v %T\n", decoratedExpressionType, decoratedExpressionType)
	}
	compatibleErr := dectype.CompatibleTypes(expectedReturnType, decoratedExpressionType)
	if compatibleErr != nil {
		return nil, decorated.NewUnMatchingFunctionReturnTypesInFunctionValue(potentialFunc, expression, expectedReturnType, decoratedExpression.Type(), compatibleErr)
	}

	return decorated.NewFunctionValue(forcedFunctionType, parameters, decoratedExpression), nil
}




