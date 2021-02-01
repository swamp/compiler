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

func parseRecordLiteral(p ParseStream, indentation int, t token.ParenToken) (ast.Expression, parerr.ParseError) {
	var recordFieldAssignments []*ast.RecordLiteralFieldAssignment

	var templateRecordIdentifier *ast.VariableIdentifier

	if _, spaceAfterCurlyErr := p.eatOneSpace("space after left curly"); spaceAfterCurlyErr != nil {
		return nil, spaceAfterCurlyErr
	}

	varIdent, wasAssign, varIdentErr := p.readVariableIdentifierAssignOrUpdate()
	if varIdentErr != nil {
		return nil, varIdentErr
	}
	if !wasAssign {
		templateRecordIdentifier = varIdent
	}
	for {
		if wasAssign {
			exp, expErr := p.parseExpressionNormal(indentation)
			if expErr != nil {
				return nil, expErr
			}
			assignment := ast.NewRecordLiteralFieldAssignment(varIdent, exp)
			recordFieldAssignments = append(recordFieldAssignments, assignment)
			wasComma, _, commaErr := p.eatCommaSeparatorOrTermination(indentation, false)
			if commaErr != nil {
				return nil, commaErr
			}

			if !wasComma {
				if eatRightErr := p.eatRightCurly(); eatRightErr != nil {
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

		_, spaceAfterAssignErr := p.eatOneSpace("space after assign (=) in record literal")
		if spaceAfterAssignErr != nil {
			return nil, parerr.NewExpectedOneSpaceAfterAssign(spaceAfterAssignErr)
		}
		wasAssign = true
	}

	return ast.NewRecordLiteral(t, templateRecordIdentifier, recordFieldAssignments), nil
}
