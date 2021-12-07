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

func parseBinaryOperator(p ParseStream, startIndentation int, infixToken token.OperatorToken, precedence Precedence, left ast.Expression) (ast.Expression, parerr.ParseError) {
	newIndentation, _, eatErr := p.eatContinuationReturnIndentation(startIndentation)
	if eatErr != nil {
		return nil, parerr.NewExpectedOneSpaceAfterBinaryOperator(eatErr)
	}
	right, rightErr := p.parseExpression(precedence, newIndentation)
	if rightErr != nil {
		return nil, rightErr
	}
	expression := ast.NewBinaryOperator(infixToken, infixToken, left, right)
	return expression, nil
}
