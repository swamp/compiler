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

func (p *Parser) parsePrefix(t token.Token, startIndentation int) (ast.Expression, parerr.ParseError) {
	// ---------------------------------------------------------------
	// Keywords
	// ---------------------------------------------------------------
	if keyword, isKeyword := t.(token.Keyword); isKeyword {
		switch keyword.Type() {
		case token.If:
			return parseIf(p.stream, keyword, startIndentation)
		case token.Let:
			return parseLet(p.stream, keyword, startIndentation)
		case token.Case:
			return parseCase(p.stream, keyword, startIndentation, p.previousComment)
		}
	}

	if t.Type() == token.Guard {
		return parseGuard(p.stream, t.(token.GuardToken), startIndentation, p.previousComment)
	}

	return nil, parerr.NewUnknownPrefixInExpression(t)
}
