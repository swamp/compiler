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

func parseTypeAlias(p ParseStream, keywordType token.Keyword, keywordAlias token.Keyword, keywordIndentation int, nameOfAlias *ast.TypeIdentifier,
	typeParameterContext *ast.LocalTypeNameDefinitionContext,
	precedingComments *ast.MultilineComment) (ast.Expression, parerr.ParseError) {

	referencedType, referencedTypeErr := parseTypeTermReference(p, keywordIndentation, typeParameterContext, precedingComments)
	if referencedTypeErr != nil {
		return nil, referencedTypeErr
	}

	/*
		typeToUse := referencedType
		if !typeParameterContext.IsEmpty() {
			typeParameterContext.SetNextType(referencedType)
			typeToUse = typeParameterContext
		}

	*/
	unusedNames := typeParameterContext.NotReferencedNames()
	if len(unusedNames) != 0 {
		return nil, ast.NewExtraTypeNameParametersError(unusedNames, referencedType)
	}

	alias := ast.NewAlias(keywordType, keywordAlias, nameOfAlias, referencedType, precedingComments)

	if !typeParameterContext.IsEmpty() {
		typeParameterContext.SetNextType(alias)
		return typeParameterContext, nil
	}

	return alias, nil
}
