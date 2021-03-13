/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func CheckForNoLint(commentBlock token.CommentBlock) string {
	for _, comment := range commentBlock.Comments {
		if strings.HasPrefix(comment.CommentString, "nolint:") {
			return strings.TrimSpace(strings.TrimPrefix(comment.CommentString, "nolint:"))
		}
	}
	return ""
}

func createVariableContextFromParameters(context *VariableContext, parameters []*decorated.FunctionParameterDefinition, forcedFunctionType *dectype.FunctionAtom, functionName *ast.VariableIdentifier) *VariableContext {
	newVariableContext := context.MakeVariableContext()

	for _, parameter := range parameters {
		namedDecoratedExpression := decorated.NewNamedDecoratedExpression(parameter.Identifier().Name(), nil, parameter)
		// namedDecoratedExpression.SetReferenced()
		newVariableContext.Add(parameter.Identifier(), namedDecoratedExpression)
	}

	self := ast.NewVariableIdentifier(token.NewVariableSymbolToken("__self", token.SourceFileReference{}, 0))
	selfDef := decorated.NewFunctionParameterDefinition(self, forcedFunctionType)
	namedDecoratedExpression := decorated.NewNamedDecoratedExpression(functionName.Name(), nil, selfDef)
	namedDecoratedExpression.SetReferenced()
	newVariableContext.Add(self, namedDecoratedExpression)

	return newVariableContext
}

func DecorateFunctionValue(d DecorateStream, annotation *decorated.AnnotationStatement, potentialFunc *ast.FunctionValue, forcedFunctionTypeLike dectype.FunctionTypeLike,
	functionName *ast.VariableIdentifier, context *VariableContext, comments *ast.MultilineComment) (*decorated.FunctionValue, decshared.DecoratedError) {
	forcedFunctionType := DerefFunctionTypeLike(forcedFunctionTypeLike)
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
	functionValueExpression := potentialFunc.Expression()
	decoratedExpression, decoratedExpressionErr := DecorateExpression(d, functionValueExpression, subVariableContext)
	if decoratedExpressionErr != nil {
		return nil, decoratedExpressionErr
	}

	decoratedExpressionType := decoratedExpression.Type()
	if decoratedExpressionType == nil {
		fmt.Printf("%v %T\n", decoratedExpressionType, decoratedExpressionType)
	}

	compatibleErr := dectype.CompatibleTypes(expectedReturnType, decoratedExpressionType)
	if compatibleErr != nil {
		return nil, decorated.NewUnMatchingFunctionReturnTypesInFunctionValue(potentialFunc, functionValueExpression, expectedReturnType, decoratedExpression.Type(), compatibleErr)
	}

	checkForNoLint := "a" // CheckForNoLint(comments)
	if checkForNoLint != "unused" {
		for _, functionVariable := range subVariableContext.InternalLookups() {
			if !functionVariable.WasReferenced() {
				_, isAssemblerFunction := potentialFunc.Expression().(*ast.Asm)
				if !isAssemblerFunction {
					// err := decorated.NewUnusedVariable(functionVariable, potentialFunc)
					// TODO: ADD THIS AS ERROR
					// d.AddDecoratedError(err)
				}
			}
		}
	} else {
		// log.Printf("info: skipping %v\n", potentialFunc.DebugFunctionIdentifier().Name())
	}

	return decorated.NewFunctionValue(annotation, potentialFunc, forcedFunctionTypeLike, parameters, decoratedExpression, comments), nil
}
