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

func parseVariableSymbol(p ParseStream, variableToken token.VariableSymbolToken) (ast.Expression, parerr.ParseError) {
	var lookups []*ast.VariableIdentifier
	ident := ast.NewVariableIdentifier(variableToken)
	for p.maybeAccessor() {
		lookup, lookupErr := p.readVariableIdentifier()
		if lookupErr != nil {
			return nil, lookupErr
		}
		lookups = append(lookups, lookup)
	}
	if len(lookups) > 0 {
		return ast.NewLookups(ident, lookups), nil
	}
	return ident, nil
}
