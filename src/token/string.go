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
		} // : stringLine.LocalStringOffset+stringLine.Length]
		if stringLine.LocalStringOffset()+stringLine.Length > strLength {
			panic(fmt.Errorf("illegal string line totalOffset '%v' %v %v vs %v", text, stringLine.LocalStringOffset(), stringLine.Length, strLength))
		}
	}

	log.Printf("NewStringToken '%v'", text)
	calculatedRange := t.CalculateRanges(0, len(text)-1)
	if !position.Range.ContainsSameLineRanges(calculatedRange) {
		t.DebugOutput()
		panic(fmt.Errorf("string literal do not contain '%v' %v \n%v\nvs\n%v (%v)", text, position, calculatedRange, position.Range, stringLines))
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
			log.Printf("found stringOffset: %v to index %v to %v", stringOffset, index, stringLine)
			return index, stringLine
		}
	}
	panic(fmt.Errorf("not found offset:%v in '%v'", stringOffset, stringLines))
}

func (s StringToken) GetStringLinesFromLocalOctetOffset(start int, end int) []SameLineRange {
	log.Printf("'%v' %v:%v (length:%v)", s.text, start, end, len(s.text))
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

	log.Printf("found '%v' startIndex %v endIndex %v, skip first line %v", s.text, startIndex, endIndex, remaining)
	newColumn := startStringLine.Position.Column() + remaining
	firstPosition := Position{
		line:                            startStringLine.Position.Line(),
		column:                          newColumn,
		originalOctetOffsetInSourceFile: start,
	}

	var runLength int
	if end <= startStringLine.LocalStringOffset()+startStringLine.Length {
		runLength = end - start
	} else {
		runLength = startStringLine.Length
	}

	startLineRange := MakeSameLineRange(firstPosition, runLength, start)
	convertedLines = append(convertedLines, startLineRange)

	for i := startIndex + 1; i < endIndex; i++ {
		stringLine := s.stringLines[i]
		//sameLineRange := MakeSameLineRange(stringLine.Position, stringLine.Length, stringLine.LocalStringOffset)
		convertedLines = append(convertedLines, stringLine)
	}

	if startIndex != endIndex {
		remainingEnd := end - endStringLine.LocalStringOffset()
		newColumnEnd := endStringLine.Position.Column()
		lastPosition := Position{
			line:                            endStringLine.Position.Line(),
			column:                          newColumnEnd,
			originalOctetOffsetInSourceFile: end,
		}

		endLineRange := MakeSameLineRange(lastPosition, remainingEnd, end)
		convertedLines = append(convertedLines, endLineRange)
	}

	log.Printf("resulting stringLines %v", convertedLines)

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
