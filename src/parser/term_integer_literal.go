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

func parseIntegerLiteral(p ParseStream, number token.NumberToken) (*ast.IntegerLiteral, parerr.ParseError) {
	return ast.NewIntegerLiteral(number, number.Value()), nil
}

func parseFixedLiteral(p ParseStream, number token.NumberToken) (*ast.FixedLiteral, parerr.ParseError) {
	return ast.NewFixedLiteral(number, number.Value()), nil
}
