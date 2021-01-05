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

func parseCase(p ParseStream, keyword token.Keyword, startIndentation int) (ast.Expression, parerr.ParseError) {
	_, firstSpaceErr := p.eatOneSpace("space after CASE")
	if firstSpaceErr != nil {
		return nil, firstSpaceErr
	}
	test, testErr := p.parseExpressionNormal(startIndentation)
	if testErr != nil {
		return nil, testErr
	}
	_, afterTestSpaceErr := p.eatOneSpace("space after case EXPRESSION")
	if afterTestSpaceErr != nil {
		return nil, afterTestSpaceErr
	}
	ofErr := p.eatOf()
	if ofErr != nil {
		return nil, ofErr
	}
	consequenceIndentation := startIndentation + 1
	_, secondSpaceErr := p.eatNewLineContinuation(startIndentation)
	if secondSpaceErr != nil {
		return nil, secondSpaceErr
	}

	var consequences []*ast.CaseConsequence
	for {
		var prefix *ast.TypeIdentifier
		var parameters []*ast.VariableIdentifier

		defaultSymbolToken, wasDefaultSymbol := p.wasDefaultSymbol()
		if wasDefaultSymbol {
			_, oneSpaceErr := p.eatOneSpace("space after default _")
			if oneSpaceErr != nil {
				return nil, oneSpaceErr
			}
			fakeSymbol := token.NewTypeSymbolToken("_", defaultSymbolToken.FetchPositionLength(), 0)
			prefix = ast.NewTypeIdentifier(fakeSymbol)
		} else {
			var prefixErr tokenize.TokenError
			prefix, prefixErr = p.readTypeIdentifier()
			if prefixErr != nil {
				return nil, parerr.NewExpectedCaseConsequenceSymbolError(prefixErr)
			}
			_, oneSpaceAfterType := p.eatOneSpace("space after case consequence type identifier")
			if oneSpaceAfterType != nil {
				return nil, oneSpaceAfterType
			}
			for {
				ident, wasIdent := p.wasVariableIdentifier()
				if !wasIdent {
					break
				}
				parameters = append(parameters, ident)
				_, oneSpaceErr := p.eatOneSpace("space after CASE consequence parameter")
				if oneSpaceErr != nil {
					return nil, oneSpaceErr
				}
			}
		}

		if arrowRightErr := p.eatRightArrow(); arrowRightErr != nil {
			return nil, parerr.NewCaseConsequenceExpectedVariableOrRightArrow(arrowRightErr)
		}

		detectedIndentation, _, oneSpaceAfterArrowErr := p.eatContinuationReturnIndentation(consequenceIndentation)
		if oneSpaceAfterArrowErr != nil {
			return nil, oneSpaceAfterArrowErr
		}
		wasIndented := detectedIndentation != consequenceIndentation
		expressionIndentation := consequenceIndentation
		if wasIndented {
			expressionIndentation = consequenceIndentation + 1
		}

		expr, exprErr := p.parseExpressionNormal(expressionIndentation)
		if exprErr != nil {
			return nil, exprErr
		}

		consequence := ast.NewCaseConsequence(prefix, parameters, expr)
		consequences = append(consequences, consequence)

		foundTextInColumnBelow, _, posLengthErr := p.maybeNewLineContinuationWithExtraEmptyLine(consequenceIndentation)
		if posLengthErr != nil {
			return nil, posLengthErr
		}
		if !foundTextInColumnBelow {
			break
		}
	}

	a := ast.NewCase(keyword, test, consequences)
	return a, nil
}
