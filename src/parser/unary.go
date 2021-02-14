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

func parseUnary(p ParseStream, startIndentation int, infixToken token.OperatorToken) (ast.Expression, parerr.ParseError) {
	left, leftErr := p.parseTerm(startIndentation)
	if leftErr != nil {
		return nil, leftErr
	}
	intValue, isInt := left.(*ast.IntegerLiteral)
	if isInt {
		return ast.NewIntegerLiteral(intValue.Token, -intValue.Value()), nil
	}
	expression := ast.NewUnaryExpression(infixToken, infixToken, left)
	return expression, nil
}
