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
)

func internalDecorateExpression(d DecorateStream, e ast.Expression, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	if e == nil {
		panic(fmt.Sprintf("expression is nil %v", context))
	}

	switch v := e.(type) {
	case *ast.Let:
		return decorateLet(d, v, context)
	case *ast.IfExpression:
		return decorateIf(d, v, context)
	case *ast.GuardExpression:
		return decorateGuard(d, v, context)
	case *ast.CaseForCustomType:
		return decorateCaseCustomType(d, v, context)
	case *ast.CaseForPatternMatching:
		return decorateCasePatternMatching(d, v, context)
	case *ast.VariableIdentifier:
		return decorateIdentifier(d, v, context)
	case *ast.VariableIdentifierScoped:
		return decorateIdentifierScoped(d, v, context)
	case *ast.IntegerLiteral:
		return decorateInteger(d, v)
	case *ast.FixedLiteral:
		return decorateFixed(d, v)
	case *ast.ResourceNameLiteral:
		return decorateResourceName(d, v)
	case *ast.StringLiteral:
		return decorateString(d, v)
	case *ast.StringInterpolation:
		return decorateStringInterpolation(d, v, context)
	case *ast.CharacterLiteral:
		return decorateCharacter(d, v)
	case *ast.TypeId:
		return decorateTypeId(d, v)
	case *ast.BooleanLiteral:
		return decorateBoolean(d, v)
	case *ast.ListLiteral:
		return decorateListLiteral(d, v, context)
	case *ast.TupleLiteral:
		return decorateTupleLiteral(d, v, context)
	case *ast.ArrayLiteral:
		return decorateArrayLiteral(d, v, context)
	case *ast.UnaryExpression:
		return decorateUnary(d, v, context)
	case *ast.FunctionCall:
		return decorateFunctionCall(d, v, context)
	case *ast.ConstructorCall:
		return decorateConstructorCall(d, v, context)
	case *ast.RecordLiteral:
		return decorateRecordLiteral(d, v, context)
	case *ast.Lookups:
		return decorateRecordLookups(d, v, context)
	case *ast.Asm:
		return decorateAsm(d, v)
	case *ast.BinaryOperator:
		return decorateBinaryOperator(d, v, context)
	default:
		return nil, decorated.NewInternalError(fmt.Errorf("don't know how to decorate %v %T %v", e, e, e.FetchPositionLength()))
	}
}

func DecorateExpression(d DecorateStream, e ast.Expression, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	expr, exprErr := internalDecorateExpression(d, e, context)
	if exprErr != nil {
		d.AddDecoratedError(exprErr)
		return nil, exprErr
	}

	if expr == nil {
		return nil, decorated.NewInternalError(fmt.Errorf("expr is nil:%v", e))
	}

	return expr, nil
}
