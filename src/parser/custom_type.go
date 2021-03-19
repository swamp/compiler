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

func parseCustomType(p ParseStream, keywordType token.Keyword, precedingComments *ast.MultilineComment, keywordIndentation int) (ast.Expression, parerr.ParseError) {
	_, firstOneSpaceErr := p.eatOneSpace("space after TYPE keyword")
	if firstOneSpaceErr != nil {
		return nil, firstOneSpaceErr
	}

	tokenAlias, isAlias := p.maybeKeywordAlias()
	if isAlias {
		_, spaceAfterAliasErr := p.eatOneSpace("space after ALIAS")
		if spaceAfterAliasErr != nil {
			return nil, spaceAfterAliasErr
		}
	}

	nameOfType, nameOfTypeErr := p.readTypeIdentifier()
	if nameOfTypeErr != nil {
		return nil, nameOfTypeErr
	}

	_, spaceAfterTypeIdentifierErr := p.eatOneSpace("space after type identifier")
	if spaceAfterTypeIdentifierErr != nil {
		return nil, spaceAfterTypeIdentifierErr
	}

	var typeParameterIdentifiers []*ast.TypeParameter

	for !p.maybeAssign() {
		typeParameterIdent, typeParameterErr := p.readVariableIdentifier()
		if typeParameterErr != nil {
			return nil, parerr.NewTypeMustBeFollowedByTypeArgumentOrEqualError(typeParameterErr)
		}
		typeParameterIdentifier := ast.NewTypeParameter(typeParameterIdent)
		typeParameterIdentifiers = append(typeParameterIdentifiers, typeParameterIdentifier)

		if _, err := p.eatOneSpace("after typeParameterIdentifiers"); err != nil {
			return nil, err
		}
	}

	typeParameterContext := ast.NewTypeParameterIdentifierContext(typeParameterIdentifiers)

	if isAlias {
		return parseTypeAlias(p, keywordType, tokenAlias, keywordIndentation, nameOfType, typeParameterContext, precedingComments)
	}

	report, afterAssignSpacingErr := p.eatNewLineContinuationAllowComment(keywordIndentation)
	if afterAssignSpacingErr != nil {
		return nil, afterAssignSpacingErr
	}

	previousComment := report.Comments.LastComment()

	expectedIndentation := keywordIndentation + 1

	var fields []*ast.CustomTypeVariant

	index := 0

	for {
		variantIdentifier, variantIdentifierErr := p.readTypeIdentifier()
		if variantIdentifierErr != nil {
			return nil, variantIdentifierErr
		}

		variantTypes, variantTypesErr := parseCustomTypeVariantTypesUntilNewline(p, keywordIndentation, nil)
		if variantTypesErr != nil {
			return nil, variantTypesErr
		}

		field := ast.NewCustomTypeVariant(index, variantIdentifier, variantTypes, previousComment)
		fields = append(fields, field)

		continuedOneColumnBelow, report, foundPosLengthErr := p.eatNewLineContinuationOrDedent(expectedIndentation)
		previousComment = report.Comments.LastComment()

		if foundPosLengthErr != nil {
			return nil, foundPosLengthErr
		}

		if !continuedOneColumnBelow {
			break
		}

		continuationOrErr := p.eatOperatorUpdate()
		if continuationOrErr != nil {
			return nil, continuationOrErr
		}

		if _, eatSpaceErr := p.eatOneSpace("after operator update"); eatSpaceErr != nil {
			return nil, eatSpaceErr
		}

		index++
	}

	newCustomType := ast.NewCustomType(keywordType, nameOfType, fields, typeParameterIdentifiers, precedingComments)

	return newCustomType, nil
}
