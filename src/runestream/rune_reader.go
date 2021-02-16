/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package runestream

import (
	"io"
	"io/ioutil"
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

// NewRuneReader :
func NewRuneReader(r io.Reader, relativeFilename string) (*RuneReader, error) {
	if len(relativeFilename) == 0 {
		panic("relative filename can not be null")
	}
	octets, octetsErr := ioutil.ReadAll(r)
	if octetsErr != nil {
		return nil, octetsErr
	}
	octets = append(octets, 0)
	return &RuneReader{octets: octets, relativeFilename: relativeFilename}, nil
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

func (s *RuneReader) DetectCurrentColumn() (int, int) {
	startPos := s.index - 1

	var detectedNonWhitespace bool

	indentationSpace := 0

	for pos := startPos; pos >= 0; pos-- {
		ch := rune(s.octets[pos])
		if ch == '\n' {
			return startPos - pos, indentationSpace
		} else if isAllowedWhitespace(ch) {
			if detectedNonWhitespace {
				if isIndentation(ch) {
					indentationSpace++
				}
			}
		} else {
			detectedNonWhitespace = true
			indentationSpace = 0
		}
	}

	return startPos, 0
}

func (s *RuneReader) Unread() rune {
	if s.index == 0 {
		panic("problem in unread")
	}
	s.index--

	return rune(s.octets[s.index])
}
