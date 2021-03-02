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
	stringToken token.StringToken
	expression  Expression
}

func (i *StringInterpolation) String() string {
	return fmt.Sprintf("'%v'", i.stringToken)
}

func NewStringInterpolation(stringToken token.StringToken, expression Expression) *StringInterpolation {
	return &StringInterpolation{expression: expression, stringToken: stringToken}
}

func (i *StringInterpolation) FetchPositionLength() token.SourceFileReference {
	return i.stringToken.FetchPositionLength()
}

func (i *StringInterpolation) StringLiteral() token.StringToken {
	return i.stringToken
}

func (i *StringInterpolation) Expression() Expression {
	return i.expression
}

func (i *StringInterpolation) DebugString() string {
	return fmt.Sprintf("[StringInterpolation %v]", i.stringToken)
}
