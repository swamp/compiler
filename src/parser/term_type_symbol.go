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

func parseTypeSymbol(p ParseStream, startIndentation int, typeSymbol token.TypeSymbolToken) (ast.Expression, parerr.ParseError) {
	x := ast.NewTypeIdentifier(typeSymbol)

	var moduleNameParts []*ast.ModuleNamePart
	for p.maybeAccessor() {
		part := ast.NewModuleNamePart(x)
		moduleNameParts = append(moduleNameParts, part)
		if variable, wasVariable := p.wasVariableIdentifier(); wasVariable {
			moduleReference := ast.NewModuleReference(moduleNameParts)
			finalVar := ast.NewQualifiedVariableIdentifier(variable, moduleReference)
			return finalVar, nil
		}
		var someErr parerr.ParseError
		x, someErr = p.readTypeIdentifier()
		if someErr != nil {
			return nil, someErr
		}
	}
	if len(moduleNameParts) > 0 {
		moduleReference := ast.NewModuleReference(moduleNameParts)
		x = ast.NewQualifiedTypeIdentifier(x.Symbol(), moduleReference)
	}

	return x, nil
}
