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

func parseParenExpression(p ParseStream, startIndentation int, parenToken token.ParenToken) (ast.Expression, parerr.ParseError) {
	p.maybeOneSpace()

	exp, expErr := p.parseExpressionNormal(startIndentation)
	if expErr != nil {
		return nil, expErr
	}

	if _, wasComma := p.maybeSpacingAndComma(startIndentation + 1); wasComma {
		return parseTuple(p, exp, startIndentation, parenToken)
	}

	p.maybeOneSpace()
	if _, rightErr := p.readRightParen(); rightErr != nil {
		return nil, rightErr
	}

	return exp, nil
}
