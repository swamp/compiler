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
	arguments, argumentsErr := parseFunctionCallArguments(p, startIndentation)
	if argumentsErr != nil {
		return nil, argumentsErr
	}


	typeIdentifier, isJustATypeIdentifier := expressionResolveToFunction.(*ast.TypeIdentifier)
	if isJustATypeIdentifier {
		return ast.NewConstructorCall(typeIdentifier, arguments), nil
	}
	call := ast.NewFunctionCall(expressionResolveToFunction, arguments)

	return call, nil
}
