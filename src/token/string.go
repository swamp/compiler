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
			panic(fmt.Errorf("illegal string line totalOffset '%v' %v vs %v", text, stringLine.StringOffset+stringLine.Length, strLength))
		}
	}

	log.Printf("NewStringToken '%v'", text)
	calculatedRange := t.CalculateRange(0, len(text))
	if !position.Range.ContainsRange(calculatedRange) {
		t.DebugOutput()
		panic(fmt.Errorf("string literal do not contain '%v' %v \n%v\nvs\n%v (%v)", text, position, calculatedRange, position.Range, stringLines))
	}

	return t
}

func (s StringToken) DebugOutput() {
	log.Printf("--- Debug Output ---")
	for index, line := range s.stringLines {
		log.Printf("  pos:%v length:%v stringOffset:%v (len:%v) '%s'", line.Position, line.Length, line.StringOffset, len(s.text), s.text)
		log.Printf("  %v:  %v length:%v '%s'", index, line.Position, line.Length, s.text[line.StringOffset:line.StringOffset+line.Length])
	}
}

func findStringLine(stringLines []StringLine, stringOffset int) (int, StringLine) {
	if stringOffset < 0 {
		panic("wrong")
	}
	for index, stringLine := range stringLines {
		start := stringLine.StringOffset
		end := stringLine.StringOffset + stringLine.Length - 1
		if stringOffset >= start && stringOffset <= end {
			return index, stringLine
		}
	}
	panic(fmt.Errorf("not found offset:%v in '%v'", stringOffset, stringLines))
}

func (s StringToken) GetStringLinesFromPosition(start int, end int) []SameLineRange {
	log.Printf("'%v' %v:%v (length:%v)", s.text, start, end, len(s.text))
	startIndex, startStringLine := findStringLine(s.stringLines, start)
	endIndex, endStringLine := findStringLine(s.stringLines, end-1)

	if start == end {
		return nil
	}

	var convertedLines []SameLineRange

	remaining := start - startStringLine.StringOffset

	log.Printf("found '%v' startIndex %v endIndex %v, skip first line %v", s.text, startIndex, endIndex, remaining)
	newColumn := startStringLine.Position.Column() + remaining
	firstPosition := Position{
		line:        startStringLine.Position.Line(),
		column:      newColumn,
		octetOffset: start,
	}

	var runLength int
	if end <= startStringLine.StringOffset+startStringLine.Length {
		runLength = end - start
	} else {
		runLength = startStringLine.Length
	}

	startLineRange := MakeSameLineRange(firstPosition, runLength)
	convertedLines = append(convertedLines, startLineRange)

	for i := startIndex + 1; i < endIndex; i++ {
		stringLine := s.stringLines[i]
		sameLineRange := MakeSameLineRange(stringLine.Position, stringLine.Length)
		convertedLines = append(convertedLines, sameLineRange)
	}

	if startIndex != endIndex {
		remainingEnd := end - endStringLine.StringOffset
		newColumnEnd := endStringLine.Position.Column()
		lastPosition := Position{
			line:        endStringLine.Position.Line(),
			column:      newColumnEnd,
			octetOffset: end,
		}

		endLineRange := MakeSameLineRange(lastPosition, remainingEnd)
		convertedLines = append(convertedLines, endLineRange)
	}

	log.Printf("resulting stringLines %v", convertedLines)

	return convertedLines
}

// CalculateRange should return stringLines TODO:
func (s StringToken) CalculateRange(start int, end int) Range {
	if start > end {
		panic("not allowed")
	}

	if start == end {
		return Range{}
	}
	ranges := s.GetStringLinesFromPosition(start, end)

	length := len(ranges)

	startPos := ranges[0].start
	lastPos := ranges[length-1]

	endPos := MakePosition(lastPos.start.line, lastPos.start.column+lastPos.length-1, lastPos.start.octetOffset)

	log.Printf("calculate '%v' %v:%v  %v %v", s.text[start:end], start, end, startPos, endPos)

	return MakeRange(startPos, endPos)
}

func (s StringToken) CalculateRanges(start int, end int) []SameLineRange {
	if start > end {
		panic("not allowed")
	}

	if start == end {
		return nil
	}

	return s.GetStringLinesFromPosition(start, end)
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
