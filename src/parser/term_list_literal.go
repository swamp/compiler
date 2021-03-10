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

func parseListLiteral(p ParseStream, startParen token.ParenToken, startIndentation int) (ast.Expression, parerr.ParseError) {
	var expressions []ast.Expression

	var wasRight bool
	var rightBracketToken token.ParenToken

	if rightBracketToken, wasRight = p.maybeRightBracket(); !wasRight {
		spaceReport, eatAfterErr := p.eatOneSpaceOrIndent("after left bracket [")
		if eatAfterErr != nil {
			return nil, eatAfterErr
		}
		subIndentation := spaceReport.ExactIndentation
		for {
			exp, expErr := p.parseExpressionNormal(subIndentation)
			if expErr != nil {
				return nil, expErr
			}
			expressions = append(expressions, exp)

			wasComma, _, commaErr := p.eatCommaSeparatorOrTermination(subIndentation, tokenize.NotAllowedAtAll)
			if commaErr != nil {
				return nil, commaErr
			}

			if !wasComma {
				var eatBracketErr parerr.ParseError
				if rightBracketToken, eatBracketErr = p.readRightBracket(); eatBracketErr != nil {
					return nil, eatBracketErr
				}
				break
			}
		}
	}

	return ast.NewListLiteral(startParen, rightBracketToken, expressions), nil
}
