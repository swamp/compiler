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

func parseLet(p ParseStream, t token.Keyword) (ast.Expression, parerr.ParseError) {
	keywordIndentation := t.FetchPositionLength().Range.FetchIndentation()

	var assignments []ast.LetAssignment

	expectedIndentation := keywordIndentation + 1

	if _, expectNewLineErr := p.eatNewLineContinuationAllowComment(keywordIndentation); expectNewLineErr != nil {
		return nil, expectNewLineErr
	}

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

		newLetAssignment := ast.NewLetAssignment(letVariableIdentifier, letExpr)
		assignments = append(assignments, newLetAssignment)

		expectingNewLetAssignment, _, posLengthErr := p.eatOneOrTwoNewLineContinuationOrDedent(expectedIndentation)
		if posLengthErr != nil {
			return nil, posLengthErr
		}

		if !expectingNewLetAssignment {
			expectedInErr := p.eatIn()
			if expectedInErr != nil {
				return nil, expectedInErr
			}
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
	a := ast.NewLet(assignments, consequence)

	return a, nil
}
