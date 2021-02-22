/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type StringLiteral struct {
	Token token.StringToken
	value string
}

func (i *StringLiteral) String() string {
	return fmt.Sprintf("'%v'", i.value)
}

func (i *StringLiteral) Value() string {
	return i.value
}

func NewStringConstant(t token.StringToken, v string) *StringLiteral {
	return &StringLiteral{value: v, Token: t}
}

func (i *StringLiteral) FetchPositionLength() token.Range {
	return i.Token.FetchPositionLength()
}

func (i *StringLiteral) StringToken() token.StringToken {
	return i.Token
}

func (i *StringLiteral) DebugString() string {
	return fmt.Sprintf("[String %v]", i.value)
}
