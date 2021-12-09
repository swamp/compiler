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

func ParseScopedOrNormalVariable(p ParseStream) (ast.ScopedOrNormalVariableIdentifier, parerr.ParseError) {
	var moduleNameParts []*ast.ModuleNamePart

	for p.detectTypeIdentifierWithoutScope() {
		x, parseErr := p.readTypeIdentifier()
		if parseErr != nil {
			return nil, parseErr
		}

		for p.maybeAccessor() {
			part := ast.NewModuleNamePart(x)
			moduleNameParts = append(moduleNameParts, part)

			if variable, wasVariable := p.wasVariableIdentifier(); wasVariable {
				moduleReference := ast.NewModuleReference(moduleNameParts)
				finalVar := ast.NewQualifiedVariableIdentifierScoped(moduleReference, variable)
				return finalVar, nil
			}
		}
	}

	return p.readVariableIdentifier()
}

func ParseLiteralOrConstant(p ParseStream, indentation int) (ast.Expression, parerr.ParseError) {
	term, termErr := p.parseTerm(indentation)
	if termErr != nil {
		return nil, termErr
	}

	_, wasInteger := term.(*ast.IntegerLiteral)
	if wasInteger {
		return term, nil
	}

	return term, nil
}

func parseExpressionStartingWithTypeSymbol(p ParseStream, startIndentation int, typeSymbol token.TypeSymbolToken) (ast.Expression, parerr.ParseError) {
	x := ast.NewTypeIdentifier(typeSymbol)

	var moduleNameParts []*ast.ModuleNamePart
	for p.maybeAccessor() {
		part := ast.NewModuleNamePart(x)
		moduleNameParts = append(moduleNameParts, part)

		if variable, wasVariable := p.wasVariableIdentifier(); wasVariable {
			moduleReference := ast.NewModuleReference(moduleNameParts)
			finalVar := ast.NewQualifiedVariableIdentifierScoped(moduleReference, variable)
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
		return ast.NewQualifiedTypeIdentifierScoped(moduleReference, x), nil
	}

	return x, nil
}
