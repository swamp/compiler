/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type Asm struct {
	tokens string
	posLen token.Range
}

func NewAsm(tokens string, posLen token.Range) *Asm {
	return &Asm{tokens: tokens, posLen: posLen}
}

func (d *Asm) Asm() string {
	return d.tokens
}

func (d *Asm) FetchPositionLength() token.Range {
	return d.posLen
}

func (d *Asm) String() string {
	return fmt.Sprintf("[asm: %v]", d.Asm())
}

func (d *Asm) DebugString() string {
	return d.String()
}
