/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

func parseTuple(p ParseStream, startExpression ast.Expression, startIndentation int, start token.ParenToken) (ast.Expression, parerr.ParseError) {
	var expressions []ast.Expression
	expressions = append(expressions, startExpression)
	var lastParen token.ParenToken
	for {
		p.maybeOneSpace()

		exp, expErr := p.parseExpressionNormal(startIndentation)
		if expErr != nil {
			return nil, expErr
		}

		expressions = append(expressions, exp)

		wasComma, _, separatorErr := p.eatCommaSeparatorOrTermination(startIndentation+1, tokenize.NotAllowedAtAll)
		if separatorErr != nil {
			return nil, separatorErr
		}
		if !wasComma {
			p.maybeOneSpace()
			foundLastParen, rightErr := p.readRightParen()
			if rightErr != nil {
				return nil, rightErr
			}
			lastParen = foundLastParen
			break
		}
	}

	tupleLiteral := ast.NewTupleLiteral(start, lastParen, expressions)

	return tupleLiteral, nil
}
