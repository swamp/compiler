/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type CharacterLiteral struct {
	Token token.CharacterToken
	value rune
}

func (i *CharacterLiteral) String() string {
	return fmt.Sprintf("'%c'", i.value)
}

func (i *CharacterLiteral) Value() rune {
	return i.value
}

func NewCharacterConstant(t token.CharacterToken, v rune) *CharacterLiteral {
	return &CharacterLiteral{value: v, Token: t}
}

func (i *CharacterLiteral) FetchPositionLength() token.Range {
	return i.Token.FetchPositionLength()
}

func (i *CharacterLiteral) CharacterToken() token.CharacterToken {
	return i.Token
}

func (i *CharacterLiteral) DebugString() string {
	return fmt.Sprintf("[Char %v]", i.value)
}
