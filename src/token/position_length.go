/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
)

type DocumentURI string

type SourceFileReference struct {
	Range Range
	Uri   DocumentURI
}

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

func (p Range) Contains(pos Position) bool {
	return (pos.line > p.start.line && pos.line < p.end.line) || (pos.line == p.start.line && pos.column >= p.start.column) || (pos.line == p.end.line && pos.column <= p.end.column)
}

func (p Range) Position() Position {
	return p.start
}

func (p Range) Start() Position {
	return p.start
}

func (p Range) End() Position {
	return p.end
}

func (p Range) FetchPositionLength() Range {
	return p
}

func (p Range) FetchIndentation() int {
	return p.indentation
}

func (p Range) String() string {
	return fmt.Sprintf("[%v to %v (%v)] ", p.start, p.end, p.indentation)
}
