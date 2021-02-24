/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func isConstant(expression decorated.DecoratedExpression) (decorated.DecoratedExpression, bool) {
	functionValue, isFunctionValue := expression.(*decorated.FunctionValue)
	if isFunctionValue {
		hasParameters := len(functionValue.Parameters()) != 0
		if hasParameters {
			return nil, false
		}
		return functionValue.Expression(), true
	}

	return nil, false
}

func decorateIdentifier(d DecorateStream, ident *ast.VariableIdentifier, context *VariableContext) (decorated.DecoratedExpression, decshared.DecoratedError) {
	expression := context.ResolveVariable(ident)
	if expression == nil {
		return nil, decorated.NewUnknownVariable(ident)
	}

	if constantExpression, wasConstant := isConstant(expression); wasConstant {
		switch t := constantExpression.(type) {
		case *decorated.IntegerLiteral:
			return t, nil
		case *decorated.StringLiteral:
			return t, nil
		case *decorated.CharacterLiteral:
			return t, nil
		case *decorated.TypeIdLiteral:
			return t, nil
		case *decorated.ResourceNameLiteral:
			return t, nil
		case *decorated.RecordLiteral:
			return t, nil
		case *decorated.ListLiteral:
			return t, nil
		case *decorated.FixedLiteral:
			return t, nil
		}
	}

	log.Printf("found variable: %T '%v'\n", expression, expression)

	return expression, nil
}
