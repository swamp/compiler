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
)

func DerefFunctionType(expectedFunctionType dtype.Type) *dectype.FunctionAtom {
	switch info := expectedFunctionType.(type) {
	case *dectype.FunctionAtom:
		return info
	case *dectype.FunctionTypeReference:
		return info.FunctionAtom()
	}

	return nil
}

func DerefFunctionTypeLike(functionTypeLike dectype.FunctionTypeLike) *dectype.FunctionAtom {
	switch info := functionTypeLike.(type) {
	case *dectype.FunctionAtom:
		return info
	case *dectype.FunctionTypeReference:
		return info.FunctionAtom()
	}
	return nil
}

func decorateConstant(d DecorateStream, nameIdent *ast.VariableIdentifier, astConstant *ast.ConstantDefinition, context *VariableContext, localCommentBlock *ast.MultilineComment) (*decorated.Constant, decshared.DecoratedError) {
	decoratedExpression, decoratedExpressionErr := DecorateExpression(d, astConstant.Expression(), context)
	if decoratedExpressionErr != nil {
		return nil, decoratedExpressionErr
	}

	return decorated.NewConstant(nameIdent, astConstant, decoratedExpression, localCommentBlock), nil
}

func defineExpressionInNamedFunctionValueHelper(d DecorateStream, context *VariableContext,
	targetPreparedFunctionValue *decorated.FunctionValue) decshared.DecoratedError {
	decoratedFunctionErr := DefineExpressionInPreparedFunctionValue(d, targetPreparedFunctionValue, context)
	if decoratedFunctionErr != nil {
		return decoratedFunctionErr
	}

	return nil
}

func defineExpressionInNamedFunctionValue(d DecorateStream, targetPreparedFunctionValue *decorated.FunctionValue, context *VariableContext, nameIdent *ast.VariableIdentifier,
	functionValue *ast.FunctionValue, expectedType dtype.Type) decshared.DecoratedError {
	name := nameIdent.Name()
	localName := name
	verboseFlag := false
	if verboseFlag {
		fmt.Printf("######### RootDefinition: %v = %v\n", localName, functionValue)
	}
	if expectedType == nil {
		err := fmt.Errorf("expected type can not be nil:%v %v", localName, functionValue)
		return decorated.NewInternalError(err)
	}

	err := defineExpressionInNamedFunctionValueHelper(d, context, targetPreparedFunctionValue)
	if err != nil {
		return err
	}

	return nil
}
