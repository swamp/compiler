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

func parseCase(p ParseStream, keyword token.Keyword, startIndentation int, previousComment token.Comment) (ast.Expression, parerr.ParseError) {
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
	ofToken, ofErr := p.readOf()
	if ofErr != nil {
		return nil, ofErr
	}
	consequenceIndentation := startIndentation + 1
	report, secondSpaceErr := p.eatNewLineContinuationAllowComment(startIndentation)
	if secondSpaceErr != nil {
		return nil, secondSpaceErr
	}

	subPreviousComent := report.Comments.LastComment()

	if p.detectTypeIdentifierWithoutScope() {
		return parseCaseForCustomType(p, test, keyword, ofToken, startIndentation, consequenceIndentation, subPreviousComent)
	}

	return parseCasePatternMatching(p, test, keyword, ofToken, startIndentation, consequenceIndentation, subPreviousComent)
}
