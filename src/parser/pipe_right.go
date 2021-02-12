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

func parsePipeRightExpression(p ParseStream, operatorToken token.OperatorToken, startIndentation int, precedence Precedence, left ast.Expression) (ast.Expression, parerr.ParseError) {
	p.maybeOneSpace()
	right, rightErr := p.parseExpression(precedence, startIndentation)
	if rightErr != nil {
		return nil, rightErr
	}

	rightCall, _ := right.(*ast.FunctionCall)
	if rightCall == nil {
		rightVar, _ := right.(*ast.VariableIdentifier)
		if rightVar == nil {
			return nil, parerr.NewRightPartOfPipeMustBeFunctionCallError(operatorToken)
		}
		rightCall = ast.NewFunctionCall(rightVar, nil)
	}

	args := rightCall.Arguments()
	args = append(args, left)
	rightCall.OverwriteArguments(args)

	return rightCall, nil
}
