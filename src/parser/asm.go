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

func parseAsm(p ParseStream, keyword token.AsmToken) (ast.Expression, parerr.ParseError) {
	asmString := keyword.Asm()
	return ast.NewAsm(asmString, keyword.FetchPositionLength()), nil
}
