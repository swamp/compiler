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

func parseTypeSymbolWithOptionalModules(p ParseStream, x *ast.TypeIdentifier) (*ast.TypeIdentifier, parerr.ParseError) {
	var moduleNameParts []*ast.ModuleNamePart

	for p.maybeAccessor() {
		part := ast.NewModuleNamePart(x)
		moduleNameParts = append(moduleNameParts, part)
		var someErr parerr.ParseError
		x, someErr = p.readTypeIdentifier()
		if someErr != nil {
			return nil, someErr
		}
	}

	if len(moduleNameParts) > 0 {
		moduleReference := ast.NewModuleReference(moduleNameParts)
		x = ast.NewQualifiedTypeIdentifier(x.Symbol(), moduleReference)
	}
	return x, nil
}

func parseTypeTermReference(p ParseStream, keywordIndentation int,
	typeParameterContext *ast.TypeParameterIdentifierContext, precedingComments token.CommentBlock) (ast.Type, parerr.ParseError) {
	return internalParseTypeTermReference(p, keywordIndentation, typeParameterContext, true, precedingComments)
}

func parseTypeVariantParameter(p ParseStream, keywordIndentation int, typeParameterContext *ast.TypeParameterIdentifierContext) (ast.Type, parerr.ParseError) {
	return internalParseTypeTermReference(p, keywordIndentation, typeParameterContext, false, token.CommentBlock{})
}

func internalParseTypeTermReference(p ParseStream, keywordIndentation int,
	typeParameterContext *ast.TypeParameterIdentifierContext,
	checkTypeParam bool, precedingComments token.CommentBlock) (ast.Type, parerr.ParseError) {
	if p.maybeLeftParen() {
		t, tErr := parseTypeReference(p, keywordIndentation, typeParameterContext, precedingComments)
		if tErr != nil {
			return nil, tErr
		}
		if rightParenErr := p.eatRightParen(); rightParenErr != nil {
			return nil, rightParenErr
		}
		return t, nil
	} else if p.maybeLeftCurly() {
		return parseRecordType(p, typeParameterContext.AllTypeParameters(), keywordIndentation, precedingComments)
	} else if foundTypeIdentifier, wasTypeSymbol := p.wasTypeIdentifier(); wasTypeSymbol {
		x, xErr := parseTypeSymbolWithOptionalModules(p, foundTypeIdentifier)
		if xErr != nil {
			return nil, xErr
		}
		var typeParameters []ast.Type
		if checkTypeParam {
			var typeParameterIdentifiersErr parerr.ParseError
			typeParameters, typeParameterIdentifiersErr = readOptionalTypeParameters(p, keywordIndentation, typeParameterContext)
			if typeParameterIdentifiersErr != nil {
				return nil, typeParameterIdentifiersErr
			}
		}

		return ast.NewTypeReference(x, typeParameters), nil
	} else if ident, wasVariableIdentifier := p.wasVariableIdentifier(); wasVariableIdentifier {
		typeParameter := ast.NewTypeParameter(ident)
		return ast.NewLocalType(typeParameter), nil
	}

	parsePosition := p.positionLength()

	return nil, parerr.NewExpectedTypeReferenceError(parsePosition)
}
