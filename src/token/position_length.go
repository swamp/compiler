/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
)

type DocumentURI string

type SourceFileDocument struct {
	Uri DocumentURI
}

type SourceFileReference struct {
	Range    Range
	Document *SourceFileDocument
}

func (s SourceFileReference) ToReferenceString() string {
	return fmt.Sprintf("%v:%d:%d:", s.Document.Uri, s.Range.start.line+1, s.Range.start.column+1)
}

func MakeInclusiveSourceFileReference(start SourceFileReference, end SourceFileReference) SourceFileReference {
	/*
		if start.Document == nil {
			panic("document can not be nil")
		}
		if start.Document != end.Document {
			panic("source file reference can not span files")
		}

	*/
	tokenRange := MakeInclusiveRange(start.Range, end.Range)
	return SourceFileReference{
		Range:    tokenRange,
		Document: nil,
	}
}

type SourceFileReferenceProvider interface {
	FetchPositionLength() SourceFileReference
}

func MakeInclusiveSourceFileReferenceSlice(references []SourceFileReferenceProvider) SourceFileReference {
	if len(references) < 1 {
		panic("MakeInclusiveSourceFileReferenceSlice can not be empty")
	}

	first := references[0]
	last := references[len(references)-1]
	return MakeInclusiveSourceFileReference(first.FetchPositionLength(), last.FetchPositionLength())
}

type Range struct {
	start       Position
	end         Position
	indentation int
}

func MakeRange(start Position, end Position) Range {
	return Range{start: start, end: end, indentation: -1}
}

func (p Range) SmallerThan(other Range) bool {
	diffLineOther := other.end.line - other.start.line
	diffLine := p.end.line - other.start.line
	if diffLine > diffLineOther {
		return false
	}

	if diffLine == diffLineOther {
		diffColOther := other.end.column - other.start.column
		diffCol := p.end.column - p.start.column

		return diffCol < diffColOther
	}

	return true
}

func (p Range) SingleLineLength() int {
	if p.start.line != p.end.line {
		return -1
	}

	return p.end.column - p.start.column + 1
}

func (p Range) IsAfter(other Range) bool {
	return (p.start.line > other.end.line) || ((p.start.line == other.end.line) && p.start.column > other.end.column)
}

func MakeInclusiveRange(start Range, end Range) Range {
	return Range{
		start:       start.Start(),
		end:         end.End(),
		indentation: start.FetchIndentation(),
	}
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
	if pos.line < p.start.line || pos.line > p.end.line {
		return false
	}

	if pos.line > p.start.line && pos.line < p.end.line {
		return true
	}

	if p.start.line == p.end.line {
		return pos.column >= p.start.column && pos.column <= p.end.column
	}

	if pos.line == p.start.line {
		return pos.column >= p.start.column
	}
	if pos.line == p.end.line {
		return pos.column <= p.end.column
	}
	panic("what happened")
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

func (p Range) FetchIndentation() int {
	return p.indentation
}

func (p Range) String() string {
	return fmt.Sprintf("[%v to %v (%v)] ", p.start, p.end, p.indentation)
}
