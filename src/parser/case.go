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

	if p.detectTypeIdentifier() {
		return parseCaseCustomType(p, test, keyword, startIndentation, consequenceIndentation)
	}

	return parseCasePatternMatching(p, test, keyword, startIndentation, consequenceIndentation)
}
