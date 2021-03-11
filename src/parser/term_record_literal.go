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

func parseRecordLiteral(p ParseStream, indentation int, t token.ParenToken) (ast.Expression, parerr.ParseError) {
	var recordFieldAssignments []*ast.RecordLiteralFieldAssignment

	var templateRecordIdentifier ast.Expression

	if _, spaceAfterCurlyErr := p.eatOneSpace("space after left curly"); spaceAfterCurlyErr != nil {
		return nil, spaceAfterCurlyErr
	}

	varIdent, wasNormalRecordLiteral, _, varIdentErr := p.readVariableIdentifierAssignOrUpdate(indentation)
	if varIdentErr != nil {
		return nil, varIdentErr
	}
	if !wasNormalRecordLiteral {
		templateRecordIdentifier = varIdent
	}
	subIndentation := indentation
	for {
		if wasNormalRecordLiteral {
			exp, expErr := p.parseExpressionNormal(subIndentation)
			if expErr != nil {
				return nil, expErr
			}
			assignment := ast.NewRecordLiteralFieldAssignment(varIdent, exp)
			recordFieldAssignments = append(recordFieldAssignments, assignment)
			wasComma, _, commaErr := p.eatCommaSeparatorOrTermination(indentation, tokenize.SameLine)
			if commaErr != nil {
				return nil, commaErr
			}

			if !wasComma {
				if _, eatRightErr := p.readRightCurly(); eatRightErr != nil {
					return nil, eatRightErr
				}
				break
			}
		}
		varIdent, varIdentErr = p.readVariableIdentifier()
		if varIdentErr != nil {
			return nil, varIdentErr
		}

		if _, spaceBeforeErr := p.eatOneSpace("after identifier and before ="); spaceBeforeErr != nil {
			return nil, spaceBeforeErr
		}

		if newAssignErr := p.eatAssign(); newAssignErr != nil {
			return nil, newAssignErr
		}

		newIndentation, _, spaceAfterAssignErr := p.eatContinuationReturnIndentation(subIndentation)
		if spaceAfterAssignErr != nil {
			return nil, parerr.NewExpectedOneSpaceAfterAssign(spaceAfterAssignErr)
		}
		subIndentation = newIndentation

		wasNormalRecordLiteral = true
	}

	return ast.NewRecordLiteral(t, templateRecordIdentifier, recordFieldAssignments), nil
}
