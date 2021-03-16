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

func parseGuard(p ParseStream, guardToken token.GuardToken, startIndentation int) (ast.Expression, parerr.ParseError) {
	var items []ast.GuardItem
	var defaultConsequence *ast.GuardDefault

	latestGuardToken := guardToken
	for {
		if _, err := p.eatOneSpace("space after guard"); err != nil {
			return nil, err
		}
		var condition ast.Expression

		if _, wasToken := p.wasDefaultSymbol(); wasToken {
			condition = nil
		} else {
			var leftErr parerr.ParseError

			condition, leftErr = p.parseExpressionNormal(startIndentation)
			if leftErr != nil {
				return nil, leftErr
			}
		}
		if _, err := p.eatOneSpace("space after expression and before ->"); err != nil {
			return nil, err
		}
		if err := p.eatRightArrow(); err != nil {
			return nil, err
		}
		subIndentation, _, err := p.eatContinuationReturnIndentation(startIndentation)
		if err != nil {
			return nil, err
		}
		consequence, consequenceErr := p.parseExpressionNormal(subIndentation)
		if consequenceErr != nil {
			return nil, consequenceErr
		}

		if condition == nil {
			defaultConsequence = &ast.GuardDefault{GuardItemBasic: ast.GuardItemBasic{
				Consequence: consequence,
				GuardToken:  latestGuardToken,
			}}
		}

		if condition != nil {
			item := ast.GuardItem{Condition: condition, GuardItemBasic: ast.GuardItemBasic{GuardToken: latestGuardToken, Consequence: consequence}}
			items = append(items, item)
		}

		wasContinuation, _, continuationErr := p.eatOneNewLineContinuationOrDedentAllowComment(startIndentation)
		if continuationErr != nil {
			return nil, continuationErr
		}
		if wasContinuation {
			if defaultConsequence != nil {
				return nil, parerr.NewExpectedDefaultLastError(consequence)
			}
			var guardErr parerr.ParseError
			if latestGuardToken, guardErr = p.readGuardPipe(); guardErr != nil {
				return nil, guardErr
			}
		} else {
			if defaultConsequence == nil {
				return nil, parerr.NewMustHaveDefaultInConditionsError(consequence)
			}
			break
		}
	}

	/*
		if p.maybeUpdate() {

		}

	*/

	expression := ast.NewGuardExpression(items, defaultConsequence)
	return expression, nil
}
