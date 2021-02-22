/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// RuneToken :
type RuneToken struct {
	Range
	text rune
}

func NewRuneToken(text rune, position Range) RuneToken {
	return RuneToken{text: text, Range: position}
}

func (s RuneToken) Type() Type {
	return RuneConstant
}

func (s RuneToken) String() string {
	return fmt.Sprintf("[rune:%v]", s.text)
}
