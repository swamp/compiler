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

func parseTypeAlias(p ParseStream, keywordType token.Keyword, keywordAlias token.Keyword, keywordIdentation int, nameOfAlias *ast.TypeIdentifier,
	typeParameterContext *ast.TypeParameterIdentifierContext,
	precedingComments *ast.MultilineComment) (*ast.Alias, parerr.ParseError) {
	newIndentation, _, spaceAfterAssignAndBeforeActualReferenceErr := p.eatContinuationReturnIndentation(keywordIdentation)
	if spaceAfterAssignAndBeforeActualReferenceErr != nil {
		return nil, spaceAfterAssignAndBeforeActualReferenceErr
	}

	referencedType, referencedTypeErr := parseTypeTermReference(p, newIndentation, typeParameterContext, precedingComments)
	if referencedTypeErr != nil {
		return nil, referencedTypeErr
	}

	return ast.NewAlias(keywordType, keywordAlias, nameOfAlias, referencedType, precedingComments), nil
}
