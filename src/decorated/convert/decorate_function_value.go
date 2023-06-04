/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"log"
	"strings"

	"github.com/swamp/compiler/src/decorated/debug"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"

	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
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

func createVariableContextFromParameters(context *VariableContext,
	parameters []*decorated.FunctionParameterDefinition) *VariableContext {
	newVariableContext := context.MakeVariableContext()

	for _, parameter := range parameters {
		if parameter.Parameter().Identifier() == nil {
			continue
		}
		namedDecoratedExpression := decorated.NewNamedDecoratedExpression(
			parameter.Parameter().Identifier().Name(), nil, parameter,
		)
		newVariableContext.Add(parameter.Parameter().Identifier(), namedDecoratedExpression)
	}

	return newVariableContext
}

func DefineExpressionInPreparedFunctionValue(d DecorateStream, targetFunctionNamedValue *decorated.NamedFunctionValue,
	context *VariableContext) decshared.DecoratedError {
	var decoratedExpression decorated.Expression
	targetFunctionValue := targetFunctionNamedValue.Value()
	subVariableContext := createVariableContextFromParameters(context, targetFunctionValue.Parameters())
	functionValueExpression := targetFunctionValue.AstFunctionValue().Expression()
	convertedDecoratedExpression, decoratedExpressionErr := DecorateExpression(
		d, functionValueExpression, subVariableContext,
	)
	if decoratedExpressionErr != nil {
		return decoratedExpressionErr
	}

	decoratedExpression = convertedDecoratedExpression

	decoratedExpressionType := decoratedExpression.Type()
	if decoratedExpressionType == nil {
		log.Printf("%v %T\n", decoratedExpressionType, decoratedExpressionType)
	}

	log.Printf("decorated Expression %s turns out to %T %s", decoratedExpression, decoratedExpressionType,
		decoratedExpressionType)

	targetFunctionNamedValue.DefineExpression(decoratedExpression)

	if targetFunctionValue.IsSomeKindOfExternal() {
		return nil
	}

	//	log.Printf("checking return type of %v : %v %T %T", targetFunctionNamedValue.FunctionName(), targetFunctionValue, targetFunctionValue.UnaliasedDeclaredFunctionType().ReturnType(), decoratedExpressionType)
	//log.Printf("checking return type of '%v'\n  %v\n  %v", targetFunctionNamedValue.FunctionName(), targetFunctionValue.UnaliasedDeclaredFunctionType().ReturnType(), decoratedExpressionType)

	log.Printf("calculated dec expr %s", debug.TreeString(decoratedExpression))

	var constructingTypes []dtype.Type
	for i := 0; i < len(targetFunctionNamedValue.Value().Parameters()); i += 1 {
		constructingTypes = append(constructingTypes, dectype.NewAnyType())
	}
	constructingTypes = append(constructingTypes, decoratedExpressionType)
	encounteredFunctionType := dectype.NewFunctionAtom(nil, constructingTypes)

	expectedFunctionType := targetFunctionNamedValue.Value().Type()

	targetFunctionParameters := dectype.NonResolvedFunctionParameters(targetFunctionValue.Type())
	targetFunctionDeclaredReturn := targetFunctionParameters[len(targetFunctionParameters)-1]

	log.Printf("checking decorated expression inside function with expected function type:\n%s\nvs\n%s",
		debug.TreeString(expectedFunctionType),
		debug.TreeString(encounteredFunctionType))
	compatibleErr := dectype.CompatibleTypes(
		expectedFunctionType, encounteredFunctionType,
	)

	if compatibleErr != nil {
		var writer strings.Builder
		fmt.Fprintf(&writer, "\n\ntargetType: %v", targetFunctionNamedValue.FunctionName())
		debug.Tree(targetFunctionDeclaredReturn, &writer)

		fmt.Fprintf(&writer, "\n\nencounteredType:")
		debug.Tree(decoratedExpression, &writer)

		log.Println(writer.String())

		return decorated.NewUnMatchingFunctionReturnTypesInFunctionValue(
			targetFunctionValue.AstFunctionValue(),
			functionValueExpression, targetFunctionValue.Type(), decoratedExpression.Type(), compatibleErr,
		)
	}

	if !targetFunctionValue.IsSomeKindOfExternal() {
		for _, param := range targetFunctionValue.Parameters() {
			if !param.WasReferenced() && !param.Parameter().IsIgnore() {
				unusedErr := decorated.NewUnusedParameter(param, targetFunctionValue)
				d.AddDecoratedError(unusedErr)
			}
		}
	}

	/*
		checkForNoLint := "a" // CheckForNoLint(comments)
		if checkForNoLint != "unused" {
		} else {
			// log.Printf("info: skipping %v\n", potentialFunc.DebugFunctionIdentifier().Name())
		}

	*/

	return nil
}
