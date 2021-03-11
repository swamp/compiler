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

func parseArrayLiteral(p ParseStream, startParen token.ParenToken, startIndentation int) (ast.Expression, parerr.ParseError) {
	var expressions []ast.Expression

	var endArrayBracket token.ParenToken

	var wasEndBracket bool

	if endArrayBracket, wasEndBracket = p.maybeRightArrayBracket(); !wasEndBracket {
		newIndentation, _, eatAfterErr := p.eatContinuationReturnIndentation(startIndentation)
		if eatAfterErr != nil {
			return nil, eatAfterErr
		}
		startIndentation = newIndentation
		for {
			exp, expErr := p.parseExpressionNormal(startIndentation)
			if expErr != nil {
				return nil, expErr
			}
			expressions = append(expressions, exp)

			wasComma, _, commaErr := p.eatCommaSeparatorOrTermination(startIndentation, tokenize.NotAllowedAtAll)
			if commaErr != nil {
				return nil, commaErr
			}

			if !wasComma {
				var eatBracketErr parerr.ParseError

				if endArrayBracket, eatBracketErr = p.readRightArrayBracket(); eatBracketErr != nil {
					return nil, eatBracketErr
				}
				break
			}
		}
	}

	return ast.NewArrayLiteral(startParen, endArrayBracket, expressions), nil
}
