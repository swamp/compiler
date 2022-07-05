/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
	"log"
)

// StringToken :
type StringToken struct {
	SourceFileReference
	text        string
	raw         string
	stringLines []SameLineRange
}

func NewStringToken(raw string, text string, position SourceFileReference, stringLines []SameLineRange) StringToken {
	t := StringToken{raw: raw, text: text, SourceFileReference: position, stringLines: stringLines}
	strLength := len(text)

	for _, stringLine := range stringLines {
		if stringLine.LocalStringOffset() < 0 || (stringLine.LocalStringOffset() >= strLength && strLength != 0) {
			panic(fmt.Errorf("illegal string line %v %v '%v'", stringLine.LocalStringOffset(), strLength, text))
		}
		if stringLine.LocalStringOffset()+stringLine.Length > strLength {
			panic(fmt.Errorf("illegal string line totalOffset '%v' %v %v vs %v", text, stringLine.LocalStringOffset(), stringLine.Length, strLength))
		}
	}

	if len(text) > 0 {
		calculatedRange := t.CalculateRanges(0, len(text)-1)
		if !position.Range.ContainsSameLineRanges(calculatedRange) {
			t.DebugOutput()
			panic(fmt.Errorf("string literal do not contain '%v' %v \n%v\nvs\n%v (%v)", text, position, calculatedRange, position.Range, stringLines))
		}
	}

	return t
}

func (s StringToken) StringLines() []SameLineRange {
	return s.stringLines
}

func (s StringToken) DebugOutput() {
	log.Printf("--- Debug Output ---")
	for index, line := range s.stringLines {
		log.Printf("  pos:%v length:%v stringOffset:%v (len:%v) '%s'", line.Position, line.Length, line.LocalStringOffset(), len(s.text), s.text)
		log.Printf("  %v:  %v length:%v '%s'", index, line.Position, line.Length, s.text[line.LocalStringOffset():line.LocalStringOffset()+line.Length])
	}
}

func findSameLineRangeFromLocalOffset(stringLines []SameLineRange, stringOffset int) (int, SameLineRange) {
	if stringOffset < 0 {
		panic("wrong")
	}
	for index, stringLine := range stringLines {
		start := stringLine.LocalStringOffset()
		end := stringLine.LocalStringOffset() + stringLine.Length - 1
		if stringOffset >= start && stringOffset <= end {
			return index, stringLine
		}
	}
	panic(fmt.Errorf("not found offset:%v in '%v'", stringOffset, stringLines))
}

func (s StringToken) GetStringLinesFromLocalOctetOffset(start int, end int) []SameLineRange {
	if start > end {
		panic("not allowed")
	}

	if start == end {
		return nil
	}

	startIndex, startStringLine := findSameLineRangeFromLocalOffset(s.stringLines, start)
	endIndex, endStringLine := findSameLineRangeFromLocalOffset(s.stringLines, end-1)

	var convertedLines []SameLineRange

	remaining := start - startStringLine.LocalStringOffset()

	newColumn := startStringLine.Position.Column() + remaining
	firstPosition := Position{
		line:                            startStringLine.Position.Line(),
		column:                          newColumn,
		originalOctetOffsetInSourceFile: startStringLine.LocalOctetOffset + remaining,
	}

	var runLength int
	if startIndex == endIndex {
		runLength = end - start
	} else {
		if start > startStringLine.LocalStringOffset()+startStringLine.Length {
			panic(fmt.Errorf("illegal start position %v", start))
		}
		runLength = startStringLine.Length - (start - startStringLine.LocalStringOffset())
		if runLength == 0 {
			panic(fmt.Errorf("zero run length is illegal"))
		}
	}

	startLineRange := MakeSameLineRange(firstPosition, runLength, start)
	convertedLines = append(convertedLines, startLineRange)

	for i := startIndex + 1; i < endIndex; i++ {
		stringLine := s.stringLines[i]
		convertedLines = append(convertedLines, stringLine)
	}

	if startIndex != endIndex {
		remainingEnd := end - endStringLine.LocalStringOffset()
		newColumnEnd := endStringLine.Position.Column()
		lastPosition := Position{
			line:                            endStringLine.Position.Line(),
			column:                          newColumnEnd,
			originalOctetOffsetInSourceFile: endStringLine.LocalOctetOffset + endStringLine.Position.Column(),
		}

		endLineRange := MakeSameLineRange(lastPosition, remainingEnd, endStringLine.LocalStringOffset())
		convertedLines = append(convertedLines, endLineRange)
	}

	return convertedLines
}

func (s StringToken) CalculateRangesWithOffset(start int, end int, octetOffset int) []SameLineRange {
	ranges := s.CalculateRanges(start, end)
	var newRanges []SameLineRange
	for _, x := range ranges {
		x.LocalOctetOffset = x.LocalStringOffset() - octetOffset
		newRanges = append(newRanges, x)
	}
	return newRanges
}

func (s StringToken) CalculateRangesWithOffsetAndString(start int, end int, octetOffset int) ([]SameLineRange, string) {
	ranges := s.CalculateRangesWithOffset(start, end, octetOffset)
	return ranges, s.text[start:end]
}

func (s StringToken) CalculateRanges(start int, end int) []SameLineRange {
	if start > end {
		panic("not allowed")
	}

	if start == end {
		return nil
	}

	return s.GetStringLinesFromLocalOctetOffset(start, end)
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
