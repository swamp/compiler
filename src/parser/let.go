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

func parseMultipleIdentifiers(p ParseStream) ([]*ast.VariableIdentifier, parerr.ParseError) {
	var identifiers []*ast.VariableIdentifier

	for {
		letVariableIdentifier, letVariableIdentifierErr := p.readVariableIdentifier()
		if letVariableIdentifierErr != nil {
			return nil, letVariableIdentifierErr
		}
		identifiers = append(identifiers, letVariableIdentifier)

		if _, wasComma := p.maybeComma(); !wasComma {
			break
		}
		if _, err := p.eatOneSpace("afterComma"); err != nil {
			return nil, err
		}
	}

	if _, spaceAfterIdentifierErr := p.eatOneSpace("after variable identifier"); spaceAfterIdentifierErr != nil {
		return nil, parerr.NewExpectedOneSpaceAfterVariableAndBeforeAssign(spaceAfterIdentifierErr)
	}

	return identifiers, nil
}

func parseRecordDestructuring(p ParseStream, keywordIndentation int) ([]*ast.VariableIdentifier, parerr.ParseError) {
	if _, err := p.eatOneSpace("after destructuring {"); err != nil {
		return nil, err
	}

	var identifiers []*ast.VariableIdentifier
	for {
		letVariableIdentifier, letVariableIdentifierErr := p.readVariableIdentifier()
		if letVariableIdentifierErr != nil {
			return nil, letVariableIdentifierErr
		}
		identifiers = append(identifiers, letVariableIdentifier)

		if _, wasComma := p.maybeComma(); !wasComma {
			break
		}
		if _, err := p.eatOneSpace("afterComma"); err != nil {
			return nil, err
		}
	}

	if _, spaceAfterIdentifierErr := p.eatOneSpace("after variable identifier"); spaceAfterIdentifierErr != nil {
		return nil, parerr.NewExpectedOneSpaceAfterVariableAndBeforeAssign(spaceAfterIdentifierErr)
	}
	if _, endCurlyErr := p.readRightCurly(); endCurlyErr != nil {
		return nil, endCurlyErr
	}
	if _, spaceAfterCurlyErr := p.eatOneSpace("after } destructuring"); spaceAfterCurlyErr != nil {
		return nil, parerr.NewExpectedOneSpaceAfterVariableAndBeforeAssign(spaceAfterCurlyErr)
	}
	return identifiers, nil
}

func parseLet(p ParseStream, t token.Keyword, keywordIndentation int) (ast.Expression, parerr.ParseError) {
	var assignments []ast.LetAssignment

	expectedIndentation := keywordIndentation + 1

	lastReport, expectNewLineErr := p.eatNewLineContinuationAllowComment(keywordIndentation)

	if expectNewLineErr != nil {
		return nil, expectNewLineErr
	}

	var inKeyword token.Keyword
	var identifiers []*ast.VariableIdentifier
	var identifiersErr parerr.ParseError
	wasCurlyDestructuring := false
	for {
		if _, wasCurly := p.maybeLeftCurly(); wasCurly {
			wasCurlyDestructuring = true
			identifiers, identifiersErr = parseRecordDestructuring(p, keywordIndentation)
		} else {
			wasCurlyDestructuring = false
			identifiers, identifiersErr = parseMultipleIdentifiers(p)
			if identifiersErr != nil {
				return nil, identifiersErr
			}

			for _, checkIdentifier := range identifiers {
				for _, existingAssignment := range assignments {
					for _, existingIdentifier := range existingAssignment.Identifiers() {
						if checkIdentifier.IsIgnore() || existingIdentifier.IsIgnore() {
							continue
						}
						if checkIdentifier.Name() == existingIdentifier.Name() {
							return nil, parerr.NewExpectedUniqueLetIdentifier(checkIdentifier.FetchPositionLength())
						}
					}
				}
			}
		}

		if expectedAssignErr := p.eatAssign(); expectedAssignErr != nil {
			return nil, expectedAssignErr
		}

		detectedIndent, _, beforeExpressionSpaceErr := p.eatContinuationReturnIndentation(expectedIndentation)
		if beforeExpressionSpaceErr != nil {
			return nil, beforeExpressionSpaceErr
		}
		wasBlock := detectedIndent > expectedIndentation
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
			multiline, wasMultiline := comment.(token.MultiLineCommentToken)
			if wasMultiline {
				astMultilineComment = ast.NewMultilineComment(multiline)
			}
		}
		newLetAssignment := ast.NewLetAssignment(wasCurlyDestructuring, identifiers, letExpr, astMultilineComment)
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
				return nil, parerr.NewExpectedInKeyword(inKeywordIdentifier.FetchPositionLength())
			}
			inKeyword = inKeywordIdentifier
			break
		}
	}
	expectedInConsequencesIndentation := keywordIndentation
	report, spaceBeforeConsequenceErr := p.eatNewLineContinuationAllowComment(expectedInConsequencesIndentation - 1)
	if spaceBeforeConsequenceErr != nil {
		return nil, parerr.NewLetInConsequenceOnSameColumn(spaceBeforeConsequenceErr)
	}

	consequence, consequenceErr := p.parseExpressionNormalWithComment(keywordIndentation, report.Comments.LastComment())
	if consequenceErr != nil {
		return nil, consequenceErr
	}
	a := ast.NewLet(t, inKeyword, assignments, consequence)

	return a, nil
}
