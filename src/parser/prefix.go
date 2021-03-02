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
	if externalFunctionToken, isExternalFunction := t.(token.ExternalFunctionToken); isExternalFunction {
		return parseExternalFunction(p.stream, externalFunctionToken)
	}

	if asmToken, isAsm := t.(token.AsmToken); isAsm {
		return parseAsm(p.stream, asmToken)
	}

	// ---------------------------------------------------------------
	// Keywords
	// ---------------------------------------------------------------
	if keyword, isKeyword := t.(token.Keyword); isKeyword {
		switch keyword.Type() {
		case token.If:
			return parseIf(p.stream, keyword, startIndentation)
		case token.Let:
			return parseLet(p.stream, keyword)
		case token.Case:
			return parseCase(p.stream, keyword, startIndentation)
		}
	}

	if t.Type() == token.Guard {
		return parseGuard(p.stream, t.(token.GuardToken), startIndentation)
	}

	return nil, parerr.NewUnknownPrefixInExpression(t)
}
