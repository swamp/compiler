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

func parseCustomType(p ParseStream, keyword token.VariableSymbolToken, precedingComments token.CommentBlock) (ast.Expression, parerr.ParseError) {
	keywordIndentation := keyword.FetchIndentation()
	_, firstOneSpaceErr := p.eatOneSpace("space after TYPE keyword")
	if firstOneSpaceErr != nil {
		return nil, firstOneSpaceErr
	}

	isAlias := p.maybeKeywordAlias()
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

	wasNewLineBeforeAssign := false

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
		return parseTypeAlias(p, keywordIndentation, nameOfType, typeParameterContext, precedingComments)
	}

	if !wasNewLineBeforeAssign {
		_, afterAssignSpacingErr := p.eatNewLineContinuation(keywordIndentation)
		if afterAssignSpacingErr != nil {
			return nil, afterAssignSpacingErr
		}
	}

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

		field := ast.NewCustomTypeVariant(index, variantIdentifier, variantTypes)
		fields = append(fields, field)

		continuedOneColumnBelow, _, foundPosLengthErr := p.maybeNewLineContinuation(expectedIndentation)
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

	var customType ast.Type
	newCustomType := ast.NewCustomType(nameOfType, fields, typeParameterIdentifiers)
	customType = newCustomType

	statement := ast.NewCustomTypeStatement(nameOfType, customType, precedingComments)
	return statement, nil
}
