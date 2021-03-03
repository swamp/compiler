/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
)

func ParseAnnotation(p ParseStream, ident *ast.VariableIdentifier,
	precedingComments *ast.MultilineComment) (ast.Expression, parerr.ParseError) {
	typeParameterContext := ast.NewTypeParameterIdentifierContext(nil)
	t, tErr := parseTypeReferenceFunc(p, ident.Symbol().FetchIndentation(), typeParameterContext, precedingComments)
	if tErr != nil {
		return nil, tErr
	}
	return ast.NewAnnotation(ident, t, precedingComments), nil
}
