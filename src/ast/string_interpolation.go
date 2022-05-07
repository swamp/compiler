/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type StringInterpolation struct {
	stringToken           token.StringToken
	expression            Expression
	referencedExpressions []Expression
}

func (i *StringInterpolation) String() string {
	return fmt.Sprintf("'%v'", i.stringToken)
}

func NewStringInterpolation(stringToken token.StringToken, expression Expression, referencedExpressions []Expression) *StringInterpolation {
	var lastExpression Expression
	for _, expr := range referencedExpressions {
		if lastExpression != nil {
			if !expr.FetchPositionLength().Range.IsAfter(lastExpression.FetchPositionLength().Range) {
				panic(fmt.Sprintf("not allowed %v %v", expr.FetchPositionLength().Range, lastExpression.FetchPositionLength().Range))
			}
		}
		lastExpression = expr
	}
	return &StringInterpolation{expression: expression, stringToken: stringToken, referencedExpressions: referencedExpressions}
}

func (i *StringInterpolation) FetchPositionLength() token.SourceFileReference {
	return i.stringToken.FetchPositionLength()
}

func (i *StringInterpolation) StringLiteral() token.StringToken {
	return i.stringToken
}

func (i *StringInterpolation) ReferencedExpressions() []Expression {
	return i.referencedExpressions
}

func (i *StringInterpolation) Expression() Expression {
	return i.expression
}

func (i *StringInterpolation) DebugString() string {
	return fmt.Sprintf("[StringInterpolation %v]", i.stringToken)
}
