package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

func parseCasePatternMatching(p ParseStream, test ast.Expression, keywordCase token.Keyword, keywordOf token.Keyword, startIndentation int, consequenceIndentation int, previousComment token.Comment) (*ast.CaseForPatternMatching, parerr.ParseError) {
	var consequences []*ast.CaseConsequencePatternMatching
	for {
		var prefix ast.Literal

		_, wasDefaultSymbol := p.wasDefaultSymbol()
		if wasDefaultSymbol {
			_, oneSpaceErr := p.eatOneSpace("space after default _")
			if oneSpaceErr != nil {
				return nil, oneSpaceErr
			}
			// fakeSymbol := token.NewTypeSymbolToken("_", defaultSymbolToken.FetchPositionLength(), 0)
			// prefix = ast.NewDefaultMarker(fakeSymbol)
		} else {
			var prefixErr tokenize.TokenError
			prefix, prefixErr = ParseLiteralOrConstant(p, startIndentation)
			if prefixErr != nil {
				return nil, parerr.NewExpectedCaseConsequenceSymbolError(prefixErr)
			}
			_, oneSpaceAfterType := p.eatOneSpace("space after case consequence pattern matching literal")
			if oneSpaceAfterType != nil {
				return nil, oneSpaceAfterType
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

		consequence := ast.NewCaseConsequenceForPatternMatching(len(consequences), prefix, expr, previousComment)
		consequences = append(consequences, consequence)

		foundTextInColumnBelow, report, posLengthErr := p.eatOneNewLineContinuationOrDedentAllowComment(consequenceIndentation)
		if posLengthErr != nil {
			return nil, posLengthErr
		}
		if !foundTextInColumnBelow {
			break
		}
		previousComment = report.Comments.LastComment()
	}

	a := ast.NewCaseForPatternMatching(keywordCase, keywordOf, test, consequences)
	return a, nil
}
