/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
	"log"
)

type StringLine struct {
	Position     Position
	Length       int
	StringOffset int
}

// StringToken :
type StringToken struct {
	SourceFileReference
	text        string
	raw         string
	stringLines []StringLine
}

func NewStringToken(raw string, text string, position SourceFileReference, stringLines []StringLine) StringToken {
	t := StringToken{raw: raw, text: text, SourceFileReference: position, stringLines: stringLines}
	strLength := len(text)

	for _, stringLine := range stringLines {
		if stringLine.StringOffset < 0 || (stringLine.StringOffset >= strLength && strLength != 0) {
			panic(fmt.Errorf("illegal string line %v %v", stringLine.StringOffset, strLength))
		} // : stringLine.StringOffset+stringLine.Length]
		if stringLine.StringOffset+stringLine.Length > strLength {
			panic(fmt.Errorf("illegal string line"))
		}
	}
	return t
}

func (s StringToken) DebugOutput() {
	for index, line := range s.stringLines {
		log.Printf("pos:%v length:%v stringOffset:%v (len:%v) '%s'", line.Position, line.Length, line.StringOffset, len(s.text), s.text)
		log.Printf("%v:  %v length:%v %s", index, line.Position, line.Length, s.text[line.StringOffset:line.StringOffset+line.Length])
	}
}

func findStringLine(stringLines []StringLine, stringOffset int) StringLine {
	if stringOffset < 0 {
		panic("wrong")
	}
	for _, stringLine := range stringLines {
		start := stringLine.StringOffset
		end := stringLine.StringOffset + stringLine.Length - 1
		if stringOffset >= start && stringOffset <= end {
			return stringLine
		}
	}
	panic(fmt.Errorf("not found offset:%v in '%v'", stringOffset, stringLines))
}

func (s StringToken) GetPosition(start int) Position {
	stringLine := findStringLine(s.stringLines, start)
	remaining := start - stringLine.StringOffset

	newColumn := stringLine.Position.Column() + remaining

	/*
		line := s.SourceFileReference.Range.start.line
		column := s.SourceFileReference.Range.start.column + 2 // assumes string interpolation starts with %" or $"
		for i := 0; i < start; i++ {
			ch := s.text[i:i]
			if ch == "\n" {
				line++
				column = 0
			} else {
				column++
			}
		}

	*/
	return Position{
		line:        stringLine.Position.Line(),
		column:      newColumn,
		octetOffset: start,
	}
}

// CalculateRange should return stringLines TODO:
func (s StringToken) CalculateRange(start int, end int) Range {
	if start > end {
		panic("not allowed")
	}
	if start == end {
		start := s.GetPosition(start)
		newRange := MakeRange(start, start)
		return newRange
	}

	startPos := s.GetPosition(start)
	endPos := s.GetPosition(end - 1)

	newRange := MakeRange(startPos, endPos)

	return newRange
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

func (s StringToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
