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

func parseExternalFunction(p ParseStream, keyword token.ExternalFunctionToken) (ast.Expression, parerr.ParseError) {
	return ast.NewExternalFunction(keyword.ExternalFunction(), keyword.ParameterCount(), keyword.FetchPositionLength()), nil
}
