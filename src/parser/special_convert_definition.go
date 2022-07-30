/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
)

func isConstant(expression ast.Expression) bool {
	switch expression.(type) {
	case *ast.IntegerLiteral:
		return true
	case *ast.StringLiteral:
		return true
	case *ast.CharacterLiteral:
		return true
	case *ast.TypeId:
		return true
	case *ast.ResourceNameLiteral:
		return true
	case *ast.RecordLiteral:
		return true
	case *ast.ListLiteral:
		return true
	case *ast.ArrayLiteral:
		return true
	case *ast.FixedLiteral:
		return true
	}

	return false
}

func parseDefinition(p ParseStream, ident *ast.VariableIdentifier,
	commentBlock *ast.MultilineComment) (ast.Expression, parerr.ParseError) {
	var parameters []*ast.VariableIdentifier
	keywordIndentation := ident.Symbol().FetchIndentation()
	for {

		if p.maybeAssign() {
			break
		}

		variable, variableErr := p.readVariableIdentifier()
		if variableErr != nil {
			return nil, variableErr
		}
		parameters = append(parameters, variable)

		_, skipAfterIdentifierErr := p.eatOneSpace("space after skip identifier in definition")
		if skipAfterIdentifierErr != nil {
			return nil, skipAfterIdentifierErr
		}
	}
	newIndentation, _, indentationErr := p.eatContinuationReturnIndentationAllowComment(keywordIndentation)
	if indentationErr != nil {
		return nil, indentationErr
	}
	expressionIndentation := newIndentation
	expr, exprErr := p.parseExpressionNormal(expressionIndentation)
	if exprErr != nil {
		return nil, exprErr
	}

	if len(parameters) == 0 && isConstant(expr) {
		return ast.NewConstantDefinition(ident, expr, commentBlock), nil
	}

	newFunction := ast.NewFunctionValue(ident.Symbol(), parameters, expr, commentBlock)

	return ast.NewFunctionValueNamedDefinition(ident, newFunction), nil
}
