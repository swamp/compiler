/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// CharacterToken :
type CharacterToken struct {
	Range
	text rune
	raw  string
}

func NewCharacterToken(raw string, text rune, position Range) CharacterToken {
	return CharacterToken{raw: raw, text: text, Range: position}
}

func (s CharacterToken) Type() Type {
	return StringConstant
}

func (s CharacterToken) String() string {
	return fmt.Sprintf("[ch:%c]", s.text)
}

func (s CharacterToken) Raw() string {
	return s.raw
}

func (s CharacterToken) Character() rune {
	return s.text
}
