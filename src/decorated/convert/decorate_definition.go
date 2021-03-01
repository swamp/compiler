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

func decorateDefinition(d DecorateStream, context *VariableContext, nameIdent *ast.VariableIdentifier, expr ast.Expression, expectedType dtype.Type, annotation *decorated.Annotation, localCommentBlock token.CommentBlock) (decorated.Expression, decshared.DecoratedError) {
	name := nameIdent.Name()
	localName := name
	verboseFlag := false
	if verboseFlag {
		fmt.Printf("######### RootDefinition: %v = %v\n", localName, expr)
	}
	if expectedType == nil {
		err := fmt.Errorf("expected type can not be nil:%v %v", localName, expr)
		return nil, decorated.NewInternalError(err)
	}

	var decoratedExpression decorated.Expression

	switch e := expr.(type) {
	case *ast.FunctionValue:
		foundFunctionType := DerefFunctionType(expectedType)
		if foundFunctionType == nil {
			return nil, decorated.NewExpectedFunctionType(expectedType, expr)
		}
		decoratedFunction, decoratedFunctionErr := DecorateFunctionValue(d, annotation, e, foundFunctionType, nameIdent, context, localCommentBlock)
		if decoratedFunctionErr != nil {
			return nil, decoratedFunctionErr
		}
		d.InternalAddNode(decoratedFunction)
		decoratedExpression = decoratedFunction
	default:
		return nil, decorated.NewInternalError(fmt.Errorf("unknown root definition:%v %T", e, e))
	}

	verboseFlag = false
	if verboseFlag {
		fmt.Printf(">>>>>>>>>>>>>> %v = %v\n", localName, decoratedExpression)
	}
	d.AddDefinition(nameIdent, decoratedExpression)
	return decoratedExpression, nil
}
