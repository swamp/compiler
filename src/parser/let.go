/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
)

func parseLet(p ParseStream, t token.Keyword) (ast.Expression, parerr.ParseError) {
	keywordIndentation := t.FetchPositionLength().Range.FetchIndentation()

	var assignments []ast.LetAssignment

	expectedIndentation := keywordIndentation + 1

	lastReport, expectNewLineErr := p.eatNewLineContinuationAllowComment(keywordIndentation)

	if expectNewLineErr != nil {
		return nil, expectNewLineErr
	}

	var inKeyword token.Keyword
	for {
		letVariableIdentifier, letVariableIdentifierErr := p.readVariableIdentifier()
		if letVariableIdentifierErr != nil {
			return nil, letVariableIdentifierErr
		}

		if _, spaceAfterIdentifierErr := p.eatOneSpace("after variable identifier"); spaceAfterIdentifierErr != nil {
			return nil, parerr.NewExpectedOneSpaceAfterVariableAndBeforeAssign(spaceAfterIdentifierErr)
		}

		if expectedAssignErr := p.eatAssign(); expectedAssignErr != nil {
			return nil, expectedAssignErr
		}

		detectedIndent, _, beforeExpressionSpaceErr := p.eatContinuationReturnIndentation(expectedIndentation)
		if beforeExpressionSpaceErr != nil {
			return nil, beforeExpressionSpaceErr
		}
		wasBlock := detectedIndent != expectedIndentation
		expressionIndentation := expectedIndentation
		if wasBlock {
			expressionIndentation = expressionIndentation + 1
		}

		letExpr, assignmentErr := p.parseExpressionNormal(expressionIndentation)
		if assignmentErr != nil {
			return nil, assignmentErr
		}

		var astMultilineComment *ast.MultilineComment
		if len(lastReport.Comments.Comments) > 0 {
			comment := lastReport.Comments.Comments[len(lastReport.Comments.Comments)-1]
			astMultilineComment = ast.NewMultilineComment(token.NewMultiLineCommentToken(comment.RawString, comment.CommentString, comment.ForDocumentation, comment.SourceFileReference))
		}
		newLetAssignment := ast.NewLetAssignment(letVariableIdentifier, letExpr, astMultilineComment)
		assignments = append(assignments, newLetAssignment)

		expectingNewLetAssignment, nextReport, posLengthErr := p.eatOneOrTwoNewLineContinuationOrDedent(expectedIndentation)
		if posLengthErr != nil {
			return nil, posLengthErr
		}
		lastReport = nextReport

		if !expectingNewLetAssignment {
			inKeywordIdentifier, expectedInErr := p.readKeyword()
			if expectedInErr != nil {
				return nil, expectedInErr
			}
			if inKeywordIdentifier.Type() != token.In {
				return nil, parerr.NewInternalError(inKeywordIdentifier.FetchPositionLength(), fmt.Errorf("expected IN keyword here"))
			}
			inKeyword = inKeywordIdentifier
			break
		}
	}
	expectedInConsequencesIndentation := keywordIndentation
	_, spaceBeforeConsequenceErr := p.eatNewLineContinuationAllowComment(expectedInConsequencesIndentation - 1)
	if spaceBeforeConsequenceErr != nil {
		return nil, parerr.NewLetInConsequenceOnSameColumn(spaceBeforeConsequenceErr)
	}

	consequence, consequenceErr := p.parseExpressionNormal(keywordIndentation)
	if consequenceErr != nil {
		return nil, consequenceErr
	}
	a := ast.NewLet(t, inKeyword, assignments, consequence)

	return a, nil
}
