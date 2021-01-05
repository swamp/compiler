/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
)

func parseArgumentExpression(p ParseStream, startIndentation int) (ast.Expression, parerr.ParseError) {
	e, eErr := p.parseTerm(startIndentation) // Must be terms
	return e, eErr
}

func parseFunctionCallArguments(p ParseStream, startIndentation int) ([]ast.Expression, parerr.ParseError) {
	var arguments []ast.Expression
	for i := 0; i < 99; i++ {
		e, eErr := parseArgumentExpression(p, startIndentation)
		if eErr != nil {
			return nil, eErr
		}
		if e == nil {
			break
		}
		arguments = append(arguments, e)

		wasEnd, _, _ := p.eatArgumentSpaceOrDetectEndOfArguments(startIndentation)
		if wasEnd {
			break
		}
	}

	return arguments, nil
}

