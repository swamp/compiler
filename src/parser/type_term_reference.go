/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"
	"github.com/swamp/compiler/src/ast"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"reflect"
)

func parseTypeSymbolWithOptionalModules(p ParseStream, x *ast.TypeIdentifier) (ast.TypeIdentifierNormalOrScoped, parerr.ParseError) {
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
		return ast.NewQualifiedTypeIdentifierScoped(moduleReference, x), nil
	}

	return x, nil
}

func parseTypeTermReference(p ParseStream, keywordIndentation int,
	typeParameterContext ast.LocalTypeNameDefinitionContextDynamic, precedingComments *ast.MultilineComment) (ast.Type, parerr.ParseError) {
	return internalParseTypeTermReference(p, keywordIndentation, typeParameterContext, true, precedingComments)
}

func parseTypeVariantParameter(p ParseStream, keywordIndentation int, typeParameterContext *ast.LocalTypeNameDefinitionContext) (ast.Type, parerr.ParseError) {
	return internalParseTypeTermReference(p, keywordIndentation, typeParameterContext, false, nil)
}

func internalParseTypeTermReference(p ParseStream, keywordIndentation int,
	typeParameterContext ast.LocalTypeNameDefinitionContextDynamic,
	checkTypeParam bool, precedingComments *ast.MultilineComment) (ast.Type, parerr.ParseError) {
	if reflect.ValueOf(typeParameterContext).IsNil() {
		panic(fmt.Errorf("can not be nil"))
	}
	if leftParen, wasLeftParen := p.maybeLeftParen(); wasLeftParen {
		t, tErr := parseTypeReference(p, keywordIndentation, typeParameterContext, precedingComments)
		if tErr != nil {
			return nil, tErr
		}
		if _, wasComma := p.maybeComma(); wasComma {
			if _, err := p.eatOneSpace("afterComma"); err != nil {
				return nil, err
			}
			return parseTupleTypeReference(p, keywordIndentation, leftParen, typeParameterContext, precedingComments, t)
		}
		if _, rightParenErr := p.readRightParen(); rightParenErr != nil {
			return nil, rightParenErr
		}
		return t, nil
	} else if leftCurly, wasLeftCurly := p.maybeLeftCurly(); wasLeftCurly {
		return parseRecordType(p, leftCurly, keywordIndentation, nil, typeParameterContext)
	} else if foundTypeIdentifier, wasTypeSymbol := p.wasTypeIdentifier(); wasTypeSymbol {
		x, xErr := parseTypeSymbolWithOptionalModules(p, foundTypeIdentifier)
		if xErr != nil {
			return nil, xErr
		}
		if foundTypeIdentifier.Name() == "Unmanaged" {
			leftAngleBracket, leftErr := p.readLeftAngleBracket()
			if leftErr != nil {
				return nil, leftErr
			}

			nativeLanguageTypeName, typeErr := p.readTypeIdentifier()
			if typeErr != nil {
				return nil, typeErr
			}

			rightAngleBracket, rightErr := p.readRightAngleBracket()
			if rightErr != nil {
				return nil, leftErr
			}
			return ast.NewUnmanagedType(leftAngleBracket, rightAngleBracket, nativeLanguageTypeName, foundTypeIdentifier, nil), nil
		}
		var typeParameters []ast.Type
		if checkTypeParam {
			var typeParameterIdentifiersErr parerr.ParseError
			typeParameters, typeParameterIdentifiersErr = readOptionalTypeParameters(p, keywordIndentation, typeParameterContext)
			if typeParameterIdentifiersErr != nil {
				return nil, typeParameterIdentifiersErr
			}
		}
		scoped, isScoped := x.(*ast.TypeIdentifierScoped)
		if isScoped {
			return ast.NewScopedTypeReference(scoped, typeParameters), nil
		}
		return ast.NewTypeReference(x.(*ast.TypeIdentifier), typeParameters), nil
	} else if ident, wasVariableIdentifier := p.wasVariableIdentifier(); wasVariableIdentifier {
		typeParameter, refErr := typeParameterContext.GetOrCreateReferenceFromName(ast.NewLocalTypeName(ident))
		if refErr != nil {
			return nil, decorated.NewInternalError(refErr)
		}
		return typeParameter, nil
	} else if asterisk, wasAsterisk := p.maybeAsterisk(); wasAsterisk {
		return ast.NewAnyMatchingType(asterisk), nil
	}

	parsePosition := p.positionLength()

	return nil, parerr.NewExpectedTypeReferenceError(parsePosition)
}
