/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package runestream

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func isIndentation(ch rune) bool {
	return ch == ' '
}

func isNewLine(ch rune) bool {
	return ch == '\n'
}

func isAllowedWhitespace(ch rune) bool {
	return isNewLine(ch) || isIndentation(ch)
}

// RuneReader :
type RuneReader struct {
	relativeFilename string
	octets           []byte
	index            int
}

func filterOutCr(octets []byte) []byte {
	var result []byte
	for _, ch := range octets {
		if ch == '\r' { // It is common on Windows for files to have CR before LF
			continue
		}
		result = append(result, ch)
	}

	return result
}

// NewRuneReader :
func NewRuneReader(r io.Reader, absoluteFilename string) (*RuneReader, error) {
	if len(absoluteFilename) == 0 {
		panic("relative filename can not be null")
	}
	if strings.Contains(absoluteFilename, "\\") {
		panic("backslash is not supported path:" + absoluteFilename)
	}

	octets, octetsErr := ioutil.ReadAll(r)
	if octetsErr != nil {
		return nil, fmt.Errorf("runereader: ReadAll %w", octetsErr)
	}

	octets = filterOutCr(octets)

	octets = append(octets, 0)

	return &RuneReader{octets: octets, relativeFilename: absoluteFilename, index: 0}, nil
}

func (s *RuneReader) Octets() []byte {
	return s.octets
}

func (s *RuneReader) RelativeFilename() string {
	return s.relativeFilename
}

func (s *RuneReader) Read() rune {
	if s.index >= len(s.octets) {
		panic("read too far")
	}

	ch := s.octets[s.index]

	s.index++

	return rune(ch)
}

func (s *RuneReader) Tell() int {
	return s.index
}

func (s *RuneReader) DetectCurrentLineLength() (int, int) {
	startPos := s.index - 1

	var detectedNonWhitespace bool

	indentationSpace := 0

	for pos := startPos; pos >= 0; pos-- {
		ch := rune(s.octets[pos])
		switch {
		case ch == '\n':
			return startPos - pos, indentationSpace
		case isAllowedWhitespace(ch):
			if detectedNonWhitespace {
				if isIndentation(ch) {
					indentationSpace++
				}
			}
		default:
			detectedNonWhitespace = true
			indentationSpace = 0
		}
	}

	return startPos + 1, 0
}

func (s *RuneReader) Unread() rune {
	if s.index == 0 {
		panic("problem in unread")
	}
	s.index--

	return rune(s.octets[s.index])
}
