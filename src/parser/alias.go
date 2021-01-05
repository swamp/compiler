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

func parseTypeAlias(p ParseStream, keywordIdentation int, nameOfAlias *ast.TypeIdentifier,
	typeParameterContext *ast.TypeParameterIdentifierContext, precedingComments token.CommentBlock) (ast.Expression, parerr.ParseError) {
	newIndentation, _, spaceAfterAssignAndBeforeActualReferenceErr := p.eatContinuationReturnIndentation(keywordIdentation)
	if spaceAfterAssignAndBeforeActualReferenceErr != nil {
		return nil, spaceAfterAssignAndBeforeActualReferenceErr
	}
	referencedType, referencedTypeErr := parseTypeTermReference(p, newIndentation, typeParameterContext, precedingComments)
	if referencedTypeErr != nil {
		return nil, referencedTypeErr
	}


	return ast.NewAliasStatement(nameOfAlias, referencedType), nil
}
