/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func isConstant(someReference decorated.Expression) (*decorated.FunctionReference, bool) {
	functionReference, isFunctionReference := someReference.(*decorated.FunctionReference)
	if !isFunctionReference {
		return nil, false
	}

	functionValue := functionReference.FunctionValue()
	hasParameters := len(functionValue.Parameters()) != 0
	if hasParameters {
		return nil, false
	}

	switch functionValue.Expression().(type) {
	case *decorated.IntegerLiteral:
		return functionReference, true
	case *decorated.StringLiteral:
		return functionReference, true
	case *decorated.CharacterLiteral:
		return functionReference, true
	case *decorated.TypeIdLiteral:
		return functionReference, true
	case *decorated.ResourceNameLiteral:
		return functionReference, true
	case *decorated.RecordLiteral:
		return functionReference, true
	case *decorated.ListLiteral:
		return functionReference, true
	case *decorated.FixedLiteral:
		return functionReference, true
	}

	return nil, false
}

func decorateIdentifier(d DecorateStream, ident *ast.VariableIdentifier, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	expression, expressionErr := context.ResolveVariable(ident)
	if expressionErr != nil {
		return nil, decorated.NewUnknownVariable(ident)
	}

	return expression, nil
}

func decorateIdentifierScoped(d DecorateStream, ident *ast.VariableIdentifierScoped, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	def := context.FindScopedNamedDecoratedExpression(ident)
	if def == nil {
		return nil, decorated.NewUnknownVariable(ident.AstVariableReference())
	}

	return ReferenceFromVariable(ident, def.Expression(), def.ModuleDefinition().OwnedByModule())
}
