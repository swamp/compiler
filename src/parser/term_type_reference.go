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

func parseTypeId(p ParseStream, typeIdToken token.TypeId, startIndentation int) (ast.Expression, parerr.ParseError) {
	typeParameterContext := ast.NewTypeParameterIdentifierContext(nil)

	userType, userTypeErr := parseTypeReference(p, startIndentation, typeParameterContext, token.CommentBlock{})
	if userTypeErr != nil {
		return nil, userTypeErr
	}

	return ast.NewTypeId(typeIdToken, userType), nil
}
