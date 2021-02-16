/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
)

func parseDefinition(p ParseStream, ident *ast.VariableIdentifier) (ast.Expression, parerr.ParseError) {
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
	_, indentationErr := p.eatNewLineContinuationAllowComment(keywordIndentation)
	if indentationErr != nil {
		return nil, indentationErr
	}
	expressionIndentation := keywordIndentation + 1
	expr, exprErr := p.parseExpressionNormal(expressionIndentation)
	if exprErr != nil {
		return nil, exprErr
	}

	newFunction := ast.NewFunctionValue(ident.Symbol(), parameters, expr)

	return ast.NewDefinitionAssignment(ident, newFunction), nil
}
