/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

type PositionLength struct {
	position    Position
	runeWidth   int
	indentation int
}

func NewPositionLength(position Position, runeCount int, indentation int) PositionLength {
	return PositionLength{position: position, runeWidth: runeCount, indentation: indentation}
}

func (p PositionLength) RuneWidth() int {
	return p.runeWidth
}

func (p PositionLength) Position() Position {
	return p.position
}

func (p PositionLength) FetchPositionLength() PositionLength {
	return p
}

func (p PositionLength) FetchIndentation() int {
	return p.indentation
}

func (p PositionLength) String() string {
	return fmt.Sprintf("[%v (%v)] ", p.position, p.indentation)
}
