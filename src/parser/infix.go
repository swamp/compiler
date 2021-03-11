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

func (p *Parser) parseInfix(left ast.Expression, startIndentation int) (ast.Expression, parerr.ParseError) {
	if left == nil {
		panic("not allowed to parse infix with nil")
	}

	if _, spaceErr := p.stream.skipMaybeSpaceAndSameIndentationOrContinuation(); spaceErr != nil {
		return nil, spaceErr
	}
	someToken, operatorErr := p.stream.readOperatorToken()
	if operatorErr != nil {
		return nil, operatorErr
	}
	operator, wasOperator := someToken.(token.OperatorToken)
	if !wasOperator {
		panic("wasnt operator")
	}

	precedence := p.stream.getPrecedenceFromToken(operator)

	return parseBinaryOperator(p.stream, startIndentation, operator, precedence, left)
}
