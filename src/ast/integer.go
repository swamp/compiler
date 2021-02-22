/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type IntegerLiteral struct {
	Token token.NumberToken
	value int32
}

func (i *IntegerLiteral) String() string {
	return fmt.Sprintf("#%v", i.value)
}

func (i *IntegerLiteral) Value() int32 {
	return i.value
}

func (i *IntegerLiteral) FetchPositionLength() token.Range {
	return i.Token.FetchPositionLength()
}

func NewIntegerLiteral(token token.NumberToken, v int32) *IntegerLiteral {
	return &IntegerLiteral{value: v, Token: token}
}

func (i *IntegerLiteral) DebugString() string {
	return fmt.Sprintf("[integer %v]", i.value)
}
