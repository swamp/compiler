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

	var typeParameterIdentifiers []*ast.LocalTypeName

	for !p.maybeAssign() {
		typeParameterIdent, typeParameterErr := p.readVariableIdentifier()
		if typeParameterErr != nil {
			return nil, parerr.NewTypeMustBeFollowedByTypeArgumentOrEqualError(typeParameterErr)
		}
		typeParameterIdentifiers = append(typeParameterIdentifiers, ast.NewLocalTypeName(typeParameterIdent))

		_, spaceErr := p.eatOneSpace("after generic parameter")
		if spaceErr != nil {
			return nil, spaceErr
		}
	}

	continuationIndentation, report, afterAssignSpacingErr := p.eatContinuationReturnIndentationAllowComment(keywordIndentation)
	if afterAssignSpacingErr != nil {
		return nil, afterAssignSpacingErr
	}

	typeParameterContext := ast.NewLocalTypeNameContext(typeParameterIdentifiers, nil)

	if isAlias {
		return parseTypeAlias(p, keywordType, tokenAlias, continuationIndentation, nameOfType, typeParameterContext, precedingComments)
	}

	previousComment := report.Comments.LastComment()

	expectedIndentation := continuationIndentation

	var fields []*ast.CustomTypeVariant

	index := 0

	for {
		variantIdentifier, variantIdentifierErr := p.readTypeIdentifier()
		if variantIdentifierErr != nil {
			return nil, variantIdentifierErr
		}

		variantTypes, variantTypesErr := parseCustomTypeVariantTypesUntilNewline(p, keywordIndentation, typeParameterContext)
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

	newCustomType := ast.NewCustomType(keywordType, nameOfType, fields, precedingComments)

	var typeToReturn ast.Type

	typeToReturn = newCustomType
	if !typeParameterContext.IsEmpty() {
		typeParameterContext.SetNextType(newCustomType)
		typeToReturn = typeParameterContext
	}

	return ast.NewCustomTypeNamedDefinition(typeToReturn), nil
}
