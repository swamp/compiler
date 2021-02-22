/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

type Range struct {
	start       Position
	end         Position
	indentation int
}

func NewPositionLength(start Position, runeCount int, indentation int) Range {
	return Range{start: start, end: Position{
		line:   start.line,
		column: start.column + runeCount - 1,
	}, indentation: indentation}
}

func (p Range) RuneWidth() int {
	return p.end.column - p.start.column + 1
}

func (p Range) Position() Position {
	return p.start
}

func (p Range) Start() Position {
	return p.start
}

func (p Range) End() Position {
	return p.start
}

func (p Range) FetchPositionLength() Range {
	return p
}

func (p Range) FetchIndentation() int {
	return p.indentation
}

func (p Range) String() string {
	return fmt.Sprintf("[%v (%v)] ", p.start, p.indentation)
}
