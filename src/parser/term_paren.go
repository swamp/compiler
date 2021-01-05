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
	p.maybeOneSpace()

	if rightErr := p.eatRightParen(); rightErr != nil {
		return nil, rightErr
	}

	return exp, nil
}
