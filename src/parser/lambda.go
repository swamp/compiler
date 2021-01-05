/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
)

func parseLambda(p ParseStream, lambdaToken token.LambdaToken, startIndentation int) (ast.Expression, parerr.ParseError) {
	var parameters []*ast.VariableIdentifier
	for {
		ident, identErr := p.readVariableIdentifier()
		if identErr != nil {
			return nil, identErr
		}
		_, spaceAfterLambdaParameterErr := p.eatOneSpace("space after lambda identifier")
		if spaceAfterLambdaParameterErr != nil {
			return nil, spaceAfterLambdaParameterErr
		}

		parameters = append(parameters, ident)
		if p.maybeRightArrow() {
			break
		}
	}

	_, spaceAfterLambdaParametersErr := p.eatOneSpace("space after lambda identifiers")
	if spaceAfterLambdaParametersErr != nil {
		return nil, spaceAfterLambdaParametersErr
	}

	if len(parameters) == 0 {
		return nil, parerr.NewMustHaveAtLeastOneParameterError(token.PositionLength{})
	}
	expr, exprErr := p.parseExpressionNormal(startIndentation)
	if exprErr != nil {
		return nil, exprErr
	}

	return ast.NewLambdaFunctionValue(lambdaToken, parameters, expr), nil
}
