/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package lspservice

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/token"
)

type DocumentVersion uint

type LineInfo struct {
	OffsetInPayload uint
	LineLength      uint
}

func (l LineInfo) OffsetFromColumn(column uint) (uint, error) {
	if column > l.LineLength {
		return 0, fmt.Errorf("we don't have column %v max %v", column, l.LineLength)
	}

	return l.OffsetInPayload + column, nil
}

type InMemoryDocument struct {
	payload     string
	version     DocumentVersion
	lineOffsets []LineInfo
}

func NewInMemoryDocument(payload string) *InMemoryDocument {
	firstLine := LineInfo{
		OffsetInPayload: 0,
		LineLength:      0,
	}

	c := &InMemoryDocument{
		payload:     "",
		version:     1,
		lineOffsets: []LineInfo{firstLine},
	}

	c.Overwrite(payload)

	return c
}

func (c *InMemoryDocument) Overwrite(payload string) {
	c.payload = payload
	var lineOffset []LineInfo

	firstLine := LineInfo{
		OffsetInPayload: 0,
		LineLength:      0,
	}

	lastLineInfo := &firstLine
	for offsetInString, ch := range payload {
		if ch == 10 {
			lastLineInfo.LineLength = uint(offsetInString) - lastLineInfo.OffsetInPayload
			lineOffset = append(lineOffset, *lastLineInfo)
			lineInfo := LineInfo{
				OffsetInPayload: uint(offsetInString) + 1,
			}
			lastLineInfo = &lineInfo
		}
	}

	lastLineInfo.LineLength = uint(len(c.payload)) - lastLineInfo.OffsetInPayload
	lineOffset = append(lineOffset, *lastLineInfo)

	c.lineOffsets = lineOffset
}

func (c *InMemoryDocument) DebugLines() {
	for lineNumber, lineOffset := range c.lineOffsets {
		log.Printf(" line %v length %v  (%v)", lineNumber, lineOffset.LineLength, lineOffset.OffsetInPayload)
	}
}

func (c *InMemoryDocument) MakeChange(editRange token.Range, text string) error {
	lineCount := len(c.lineOffsets)
	if editRange.Start().Line() > lineCount {
		return fmt.Errorf("dont have line %v, max is %v", editRange.Start().Line(), lineCount)
	}
	startInfo := c.lineOffsets[editRange.Start().Line()]
	startOffset, startOffsetErr := startInfo.OffsetFromColumn(uint(editRange.Start().Column()))
	if startOffsetErr != nil {
		return startOffsetErr
	}

	endInfo := c.lineOffsets[editRange.End().Line()]
	endOffset, endOffsetErr := endInfo.OffsetFromColumn(uint(editRange.End().Column()))
	if endOffsetErr != nil {
		return endOffsetErr
	}
	if startOffset > endOffset {
		return fmt.Errorf("offsets are in wrong order %v %v", startOffset, endOffset)
	}

	keepBefore := c.payload[:startOffset]
	keepAfter := c.payload[endOffset:]

	newPayload := keepBefore + text + keepAfter

	log.Printf("newPayload '%v'", newPayload)
	c.Overwrite(newPayload)

	c.version++
	return nil
}

func (c *InMemoryDocument) UpdateVersion(newVersion DocumentVersion) error {
	if c.version+1 != newVersion {
		return fmt.Errorf("problem with version. expected %v but got %v", c.version+1, newVersion)
	}

	c.version = newVersion
	return nil
}
