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

func (p *Parser) parseTermUsingToken(someToken token.Token, startIndentation int) (ast.Expression, parerr.ParseError) {
	// ---------------------------------------------------------------
	// Term
	// ---------------------------------------------------------------
	switch t := someToken.(type) {
	case token.VariableSymbolToken:
		{
			return parseVariableSymbol(p.stream, t)
		}
	case token.TypeSymbolToken:
		{
			return parseTypeSymbol(p.stream, startIndentation, t)
		}
	case token.StringToken:
		{
			return parseStringLiteral(p.stream, t)
		}
	case token.CharacterToken:
		{
			return parseCharacterLiteral(p.stream, t)
		}
	case token.NumberToken:
		{
			if t.Type() == token.NumberFixed {
				return parseFixedLiteral(p.stream, t)
			}
			return parseIntegerLiteral(p.stream, t)
		}
	case token.ResourceName:
		{
			return parseResourceNameLiteral(p.stream, t)
		}
	case token.BooleanToken:
		{
			return parseBoolLiteral(p.stream, t)
		}
	case token.OperatorToken:
		{
			switch t.Type() {
			case token.OperatorBitwiseNot:
				{
					return parseUnary(p.stream, startIndentation, t)
				}
			case token.OperatorUnaryNot:
				{
					return parseUnary(p.stream, startIndentation, t)
				}
			case token.OperatorUnaryMinus:
				{
					return parseUnary(p.stream, startIndentation, t)
				}
			}
		}
	case token.GuardToken:
		{
			return parseGuard(p.stream, startIndentation)
		}
	case token.TypeId:
		{
			return parseTypeId(p.stream, t, startIndentation)
		}
	}

	// ---------------------------------------------------------------
	// Term - Block
	// ---------------------------------------------------------------
	someParenToken, wasParen := someToken.(token.ParenToken)
	if !wasParen {
		return nil, parerr.NewNotATermError(someToken)
	}

	switch someParenToken.Type() {
	case token.LeftParen:
		return parseParenExpression(p.stream, startIndentation, someParenToken)
	case token.LeftCurlyBrace:
		return parseRecordLiteral(p.stream, startIndentation, someParenToken)
	case token.LeftBracket:
		return parseListLiteral(p.stream, someParenToken, startIndentation)
	default:
		return nil, parerr.NewNotATermError(someToken)
	}
}

func (p *Parser) parseTerm(startIndentation int) (ast.Expression, parerr.ParseError) {
	someToken, someTokenErr := p.stream.readTermToken()
	if someTokenErr != nil {
		return nil, someTokenErr
	}

	return p.parseTermUsingToken(someToken, startIndentation)
}
