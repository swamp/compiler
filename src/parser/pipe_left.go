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

func parsePipeLeftExpression(p ParseStream, operatorToken token.OperatorToken, startIndentation int, precedence Precedence, left ast.Expression) (ast.Expression, parerr.ParseError) {
	_, spaceErr := p.eatOneSpace("space after pipe left")
	if spaceErr != nil {
		return nil, spaceErr
	}
	right, rightErr := p.parseExpressionNormal(startIndentation)
	if rightErr != nil {
		return nil, rightErr
	}

	leftCall, _ := left.(ast.FunctionCaller)
	if leftCall == nil {
		leftVar, _ := left.(*ast.VariableIdentifier)
		if leftVar == nil {
			return nil, parerr.NewLeftPartOfPipeMustBeFunctionCallError(operatorToken)
		}
		leftCall = ast.NewFunctionCall(leftVar, nil)
	}

	rightCall, _ := right.(ast.FunctionCaller)
	if rightCall == nil {
		return nil, parerr.NewRightPartOfPipeMustBeFunctionCallError(operatorToken)
	}

	args := leftCall.Arguments()
	args = append(args, rightCall)
	leftCall.OverwriteArguments(args)

	return leftCall, nil
}
