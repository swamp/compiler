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

func decorateNamedFunctionValue(d DecorateStream, context *VariableContext, nameIdent *ast.VariableIdentifier, functionValue *ast.FunctionValue, expectedType dtype.Type, annotation *decorated.AnnotationStatement, localCommentBlock *ast.MultilineComment) (*decorated.NamedFunctionValue, decshared.DecoratedError) {
	name := nameIdent.Name()
	localName := name
	verboseFlag := false
	if verboseFlag {
		fmt.Printf("######### RootDefinition: %v = %v\n", localName, functionValue)
	}
	if expectedType == nil {
		err := fmt.Errorf("expected type can not be nil:%v %v", localName, functionValue)
		return nil, decorated.NewInternalError(err)
	}

	var decoratedExpression decorated.Expression

	foundFunctionType := DerefFunctionType(annotation.Type())
	if foundFunctionType == nil {
		return nil, decorated.NewExpectedFunctionType(expectedType, functionValue)
	}

	decoratedFunction, decoratedFunctionErr := DecorateFunctionValue(d, annotation, functionValue, foundFunctionType, nameIdent, context, localCommentBlock)
	if decoratedFunctionErr != nil {
		return nil, decoratedFunctionErr
	}

	verboseFlag = false
	if verboseFlag {
		fmt.Printf(">>>>>>>>>>>>>> %v = %v\n", localName, decoratedExpression)
	}
	d.AddDefinition(nameIdent, decoratedFunction)

	return decorated.NewNamedFunctionValue(nameIdent, decoratedFunction), nil
}
