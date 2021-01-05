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

func parseListLiteral(p ParseStream, startParen token.ParenToken, startIndentation int) (ast.Expression, parerr.ParseError) {
	var expressions []ast.Expression

	if !p.maybeRightBracket() {
		if _, eatAfterErr := p.eatOneSpace("after left bracket ["); eatAfterErr != nil {
			return nil, eatAfterErr
		}
		for {
			exp, expErr := p.parseExpressionNormal(startIndentation)
			if expErr != nil {
				return nil, expErr
			}
			expressions = append(expressions, exp)

			wasComma, _, commaErr := p.eatCommaSeparatorOrTermination(startIndentation, false)
			if commaErr != nil {
				return nil, commaErr
			}

			if !wasComma {
				if eatBracketErr := p.eatRightBracket(); eatBracketErr != nil {
					return nil, eatBracketErr
				}
				break
			}
		}
	}

	return ast.NewListLiteral(startParen, expressions), nil
}
