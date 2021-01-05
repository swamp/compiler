/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// StringToken :
type StringToken struct {
	PositionLength
	text string
	raw  string
}

func NewStringToken(raw string, text string, position PositionLength) StringToken {
	return StringToken{raw: raw, text: text, PositionLength: position}
}

func (s StringToken) Type() Type {
	return StringConstant
}

func (s StringToken) String() string {
	return fmt.Sprintf("[str:%s]", s.text)
}

func (s StringToken) Raw() string {
	return s.raw
}

func (s StringToken) Text() string {
	return s.text
}
