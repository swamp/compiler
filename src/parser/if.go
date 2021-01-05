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

func parseIf(p ParseStream, keyword token.Keyword, keywordIndentation int) (ast.Expression, parerr.ParseError) {
	_, spaceAfterIfErr := p.eatOneSpace("space after IF")
	if spaceAfterIfErr != nil {
		return nil, spaceAfterIfErr
	}

	condition, conditionErr := p.parseExpressionNormal(keywordIndentation)
	if conditionErr != nil {
		return nil, conditionErr
	}

	_, spaceAfterExpression := p.eatOneSpace("space after IF expression")
	if spaceAfterExpression != nil {
		return nil, spaceAfterExpression
	}
	if thenErr := p.eatThen(); thenErr != nil {
		return nil, thenErr
	}

	foundIndentationAfterIf, _, spaceAfterIfErr := p.eatContinuationReturnIndentation(keywordIndentation)
	if spaceAfterIfErr != nil {
		return nil, spaceAfterIfErr
	}

	isIndentedBlock := foundIndentationAfterIf != keywordIndentation
	consequence, consequenceErr := p.parseExpressionNormal(keywordIndentation)
	if consequenceErr != nil {
		return nil, consequenceErr
	}

	_, spaceAfterConsequenceErr := p.eatBlockSpacingOneExtraLine(isIndentedBlock, keywordIndentation)
	if spaceAfterConsequenceErr != nil {
		return nil, spaceAfterConsequenceErr
	}

	if elseErr := p.eatElse(); elseErr != nil {
		return nil, parerr.NewExpectedElseKeyword(elseErr)
	}

	foundIndentationAfterElse, _, spaceAfterElseErr := p.eatContinuationReturnIndentation(keywordIndentation)
	if spaceAfterElseErr != nil {
		return nil, spaceAfterElseErr
	}

	alternative, alternativeErr := p.parseExpressionNormal(foundIndentationAfterElse)
	if alternativeErr != nil {
		return nil, parerr.NewMissingElseExpression(alternativeErr)
	}

	expression := ast.NewIfExpression(condition, consequence, alternative)

	return expression, nil
}
