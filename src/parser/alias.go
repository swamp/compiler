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
	precedingComments *ast.MultilineComment) (*ast.Alias, parerr.ParseError) {

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

	return ast.NewAlias(keywordType, keywordAlias, nameOfAlias, referencedType, precedingComments), nil
}
