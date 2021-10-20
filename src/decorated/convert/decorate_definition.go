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
