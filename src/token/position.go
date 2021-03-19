/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

type Position struct {
	line        int
	column      int
	octetOffset int
}

func NewPositionTopLeft() Position {
	return Position{line: 0, column: 0, octetOffset: 0}
}

func MakePosition(line int, column int, octetOffset int) Position {
	return Position{line: line, column: column, octetOffset: octetOffset}
}

func (p Position) Line() int {
	return p.line
}

func (p Position) Column() int {
	return p.column
}

func (p Position) NextLine(octetOffset int) Position {
	return Position{line: p.line + 1, column: p.column, octetOffset: octetOffset}
}

func (p Position) FirstColumn(octetOffset int) Position {
	return Position{line: p.line, column: 0, octetOffset: octetOffset}
}

func (p Position) NewLine() Position {
	return Position{line: p.line + 1, column: 0, octetOffset: p.octetOffset + 1}
}

func (p Position) NextColumn() Position {
	return Position{line: p.line, column: p.column + 1, octetOffset: p.octetOffset + 1}
}

func (p Position) PreviousColumn() Position {
	return Position{line: p.line, column: p.column - 1, octetOffset: p.octetOffset - 1}
}

func (p Position) OctetOffset() int {
	return p.octetOffset
}

func (p Position) String() string {
	return fmt.Sprintf("[%d:%d]", p.line+1, p.column+1)
}

func (p Position) DebugString() string {
	return fmt.Sprintf("[%d:%d](%d)", p.line+1, p.column+1, p.octetOffset)
}
