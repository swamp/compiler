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

func parseRecordTypeFields(p ParseStream, expectedIndentation int,
	parameterIdentifierContext *ast.TypeParameterIdentifierContext,
	precedingComments token.CommentBlock) ([]*ast.RecordField, parerr.ParseError) {
	var fields []*ast.RecordField
	index := 0
	for {
		symbolToken, symbolTokenErr := p.readVariableIdentifier()
		if symbolTokenErr != nil {
			return nil, symbolTokenErr
		}

		if _, spaceErr := p.eatOneSpace("before colon"); spaceErr != nil {
			return nil, spaceErr
		}

		if err := p.eatColon(); err != nil {
			return nil, err
		}

		if _, spaceErr := p.eatOneSpace("after colon"); spaceErr != nil {
			return nil, parerr.NewOneSpaceAfterRecordTypeColon(spaceErr)
		}

		userType, userTypeErr := parseTypeReference(p, expectedIndentation, parameterIdentifierContext, precedingComments)
		if userTypeErr != nil {
			return nil, userTypeErr
		}

		var report token.IndentationReport
		var wasErr parerr.ParseError
		var wasComma bool

		if wasComma, report, wasErr = p.eatCommaSeparatorOrTermination(expectedIndentation, tokenize.SameLine); wasErr != nil {
			return nil, wasErr
		}
		precedingComments := report.Comments
		field := ast.NewRecordTypeField(index, symbolToken, userType, precedingComments)
		index++
		fields = append(fields, field)

		if !wasComma {
			if err := p.readRightCurly(); err != nil {
				return nil, err
			}
			break
		}
	}
	return fields, nil
}

func parseRecordType(p ParseStream, startCurly token.ParenToken, typeParameters []*ast.TypeParameter, keywordIndentation int,
	precedingComments token.CommentBlock) (ast.Type, parerr.ParseError) {
	if _, err := p.eatOneSpace("after record type left curly"); err != nil {
		return nil, err
	}

	fields, fieldsErr := parseRecordTypeFields(p, keywordIndentation, nil, precedingComments)
	if fieldsErr != nil {
		return nil, fieldsErr
	}

	recordType := ast.NewRecordType(fields, typeParameters)
	return recordType, nil
}
