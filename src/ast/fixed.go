/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type FixedLiteral struct {
	Token token.NumberToken
	value int32
}

func (i *FixedLiteral) String() string {
	return fmt.Sprintf("#!%v", i.value)
}

func (i *FixedLiteral) Value() int32 {
	return i.value
}

func (i *FixedLiteral) PositionLength() token.PositionLength {
	return i.Token.FetchPositionLength()
}

func NewFixedLiteral(token token.NumberToken, v int32) *FixedLiteral {
	return &FixedLiteral{value: v, Token: token}
}

func (i *FixedLiteral) DebugString() string {
	return fmt.Sprintf("[fixed %v]", i.value)
}
