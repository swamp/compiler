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
	precedingComments *ast.MultilineComment) ([]*ast.RecordTypeField, parerr.ParseError) {
	var fields []*ast.RecordTypeField
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

		if wasComma, report, wasErr = p.eatCommaSeparatorOrTermination(expectedIndentation, tokenize.OwnLine); wasErr != nil {
			return nil, wasErr
		}
		field := ast.NewRecordTypeField(index, symbolToken, userType, precedingComments)
		foundComments := ast.CommentBlockToAst(report.Comments)
		if len(foundComments) > 0 {
			precedingComments = foundComments[len(foundComments)-1]
		} else {
			precedingComments = nil
		}
		index++
		fields = append(fields, field)

		if !wasComma {
			break
		}
	}
	return fields, nil
}

func parseRecordType(p ParseStream, startCurly token.ParenToken, typeParameters []*ast.TypeParameter, keywordIndentation int,
	precedingComments *ast.MultilineComment) (ast.Type, parerr.ParseError) {
	if _, err := p.eatOneSpace("after record type left curly"); err != nil {
		return nil, err
	}

	fields, fieldsErr := parseRecordTypeFields(p, keywordIndentation, nil, precedingComments)
	if fieldsErr != nil {
		return nil, fieldsErr
	}

	var rightCurlyErr parerr.ParseError
	var rightCurly token.ParenToken

	if rightCurly, rightCurlyErr = p.readRightCurly(); rightCurlyErr != nil {
		return nil, rightCurlyErr
	}

	recordType := ast.NewRecordType(startCurly, rightCurly, fields, typeParameters, precedingComments)
	return recordType, nil
}
