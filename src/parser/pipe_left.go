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
	right, rightErr := p.parseExpression(precedence, startIndentation)
	if rightErr != nil {
		return nil, rightErr
	}
	leftCall, _ := left.(*ast.FunctionCall)
	if leftCall == nil {
		leftVar, _ := left.(*ast.VariableIdentifier)
		if leftVar == nil {
			return nil, parerr.NewLeftPartOfPipeMustBeFunctionCallError(operatorToken)
		}
		leftCall = ast.NewFunctionCall(leftVar, nil)
	}

	rightCall, _ := right.(*ast.FunctionCall)
	if rightCall == nil {
		return nil, parerr.NewRightPartOfPipeMustBeFunctionCallError(operatorToken)
	}
	p.maybeOneSpace()
	if p.maybePipeLeft() {
		p.maybeOneSpace()
		expressionToAppend, expressionToAppendErr := p.parseExpression(precedence, startIndentation)
		if expressionToAppendErr != nil {
			return nil, expressionToAppendErr
		}
		innerArgs := rightCall.Arguments()
		innerArgs = append(innerArgs, expressionToAppend)
		rightCall.OverwriteArguments(innerArgs)
	}
	args := leftCall.Arguments()
	args = append(args, rightCall)
	leftCall.OverwriteArguments(args)

	return leftCall, nil
}
