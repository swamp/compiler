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
)

/*
func CheckForNoLint(commentBlock token.Comment) string {
	for _, comment := range commentBlock.Comments {
		if strings.HasPrefix(comment.CommentString, "nolint:") {
			return strings.TrimSpace(strings.TrimPrefix(comment.CommentString, "nolint:"))
		}
	}
	return ""
}
*/

func createVariableContextFromParameters(context *VariableContext, parameters []*decorated.FunctionParameterDefinition) *VariableContext {
	newVariableContext := context.MakeVariableContext()

	for _, parameter := range parameters {
		namedDecoratedExpression := decorated.NewNamedDecoratedExpression(parameter.Identifier().Name(), nil, parameter)
		newVariableContext.Add(parameter.Identifier(), namedDecoratedExpression)
	}

	return newVariableContext
}

func checkParameterCount(forcedFunctionType *dectype.FunctionAtom, potentialFunc *ast.FunctionValue) decshared.DecoratedError {
	foundAnyMatching, minimalCount := dectype.HasAnyMatchingTypes(forcedFunctionType.FunctionParameterTypes())
	if foundAnyMatching {
		if (len(potentialFunc.Parameters()) + 1) < minimalCount {
			return decorated.NewWrongNumberOfArgumentsInFunctionValue(potentialFunc, forcedFunctionType, forcedFunctionType.FunctionParameterTypes())
		}
		return nil
	}

	parameterTypes, _ := forcedFunctionType.ParameterAndReturn()
	if len(parameterTypes) != len(potentialFunc.Parameters()) {
		return decorated.NewWrongNumberOfArgumentsInFunctionValue(potentialFunc, forcedFunctionType, parameterTypes)
	}

	return nil
}

func DefineExpressionInPreparedFunctionValue(d DecorateStream, targetFunctionValue *decorated.FunctionValue, context *VariableContext) decshared.DecoratedError {
	annotation := targetFunctionValue.Annotation()

	var decoratedExpression decorated.Expression
	if !annotation.Annotation().IsSomeKindOfExternal() {
		subVariableContext := createVariableContextFromParameters(context, targetFunctionValue.Parameters())
		functionValueExpression := targetFunctionValue.AstFunctionValue().Expression()
		convertedDecoratedExpression, decoratedExpressionErr := DecorateExpression(d, functionValueExpression, subVariableContext)
		if decoratedExpressionErr != nil {
			return decoratedExpressionErr
		}

		decoratedExpression = convertedDecoratedExpression

		decoratedExpressionType := decoratedExpression.Type()
		if decoratedExpressionType == nil {
			fmt.Printf("%v %T\n", decoratedExpressionType, decoratedExpressionType)
		}

		compatibleErr := dectype.CompatibleTypes(targetFunctionValue.ForcedFunctionType().ReturnType(), decoratedExpressionType)
		if compatibleErr != nil {
			return decorated.NewUnMatchingFunctionReturnTypesInFunctionValue(targetFunctionValue.AstFunctionValue(),
				functionValueExpression, targetFunctionValue.Type(), decoratedExpression.Type(), compatibleErr)
		}

		for _, param := range targetFunctionValue.Parameters() {
			if !param.WasReferenced() && !param.Identifier().IsIgnore() {
				unusedErr := decorated.NewUnusedParameter(param, targetFunctionValue)
				d.AddDecoratedError(unusedErr)
			}
		}

		checkForNoLint := "a" // CheckForNoLint(comments)
		if checkForNoLint != "unused" {
		} else {
			// log.Printf("info: skipping %v\n", potentialFunc.DebugFunctionIdentifier().Name())
		}
	} else {
		decoratedExpression = annotation
	}

	targetFunctionValue.DefineExpression(decoratedExpression)

	return nil
}
