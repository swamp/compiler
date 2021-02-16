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

func parseCharacterLiteral(p ParseStream, characterToken token.CharacterToken) (ast.Expression, parerr.ParseError) {
	lit := ast.NewCharacterConstant(characterToken, characterToken.Character())

	return lit, nil
}
