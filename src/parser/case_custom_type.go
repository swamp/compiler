package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

func parseCaseForCustomType(p ParseStream, test ast.Expression, keywordCase token.Keyword, keywordOf token.Keyword, startIndentation int, consequenceIndentation int) (*ast.CaseForCustomType, parerr.ParseError) {
	var consequences []*ast.CaseConsequenceForCustomType
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
			_, oneSpaceAfterType := p.eatOneSpace("space after case consequence literal")
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

		consequence := ast.NewCaseConsequenceForCustomType(prefix, parameters, expr)
		consequences = append(consequences, consequence)

		foundTextInColumnBelow, _, posLengthErr := p.eatOneNewLineContinuationOrDedentAllowComment(consequenceIndentation)
		if posLengthErr != nil {
			return nil, posLengthErr
		}
		if !foundTextInColumnBelow {
			break
		}
	}

	a := ast.NewCaseForCustomType(keywordCase, keywordOf, test, consequences)
	return a, nil
}
