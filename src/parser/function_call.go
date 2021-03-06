/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
)

func parseFunctionCall(p ParseStream, startIndentation int, expressionResolveToFunction ast.Expression) (ast.Expression, parerr.ParseError) {
	_, isJustATypeIdentifier := expressionResolveToFunction.(ast.TypeIdentifierNormalOrScoped)
	arguments, argumentsErr := parseFunctionCallArguments(p, startIndentation)
	if argumentsErr != nil {
		return nil, argumentsErr
	}

	if isJustATypeIdentifier {
		scoped, wasScoped := expressionResolveToFunction.(*ast.TypeIdentifierScoped)
		var someRef ast.TypeReferenceScopedOrNormal
		if wasScoped {
			someRef = ast.NewScopedTypeReference(scoped, nil)
		} else {
			someRef = ast.NewTypeReference(expressionResolveToFunction.(*ast.TypeIdentifier), nil)
		}
		return ast.NewConstructorCall(someRef, arguments), nil
	}

	call := ast.NewFunctionCall(expressionResolveToFunction, arguments)

	return call, nil
}
