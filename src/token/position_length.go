/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"
)

type DocumentURI string

func MakeDocumentURI(s string) DocumentURI {
	if !strings.HasPrefix(s, "file://") {
		panic("illegal uri")
	}
	return DocumentURI(s)
}

func MakeDocumentURIFromLocalPath(s string) DocumentURI {
	if strings.HasPrefix(s, "file://") {
		panic("illegal localPath")
	}

	if isWindowsDrivePrefix(s) {
		s = "/" + strings.ToUpper(string(s[0])) + s[1:]
	}

	u := url.URL{
		Scheme: "file",
		Path:   s,
	}

	return DocumentURI(u.String())
}

func isWindowsDrivePrefix(path string) bool {
	if len(path) < 3 {
		return false
	}
	return unicode.IsLetter(rune(path[0])) && path[1] == ':'
}

func isWindowsDriveURI(uri string) bool {
	if len(uri) < 4 {
		return false
	}
	return uri[0] == '/' && unicode.IsLetter(rune(uri[1])) && uri[2] == ':'
}

func (s DocumentURI) ToLocalFilePath() (string, error) {
	fullUrl, urlErr := url.ParseRequestURI(string(s))
	if urlErr != nil {
		return "", urlErr
	}

	if fullUrl.Scheme != "file" {
		panic("must have file prefix")
	}

	pathOnly := fullUrl.Path

	if isWindowsDriveURI(pathOnly) {
		pathOnly = strings.ToUpper(string(fullUrl.Path[1])) + fullUrl.Path[2:]
	}

	return pathOnly, nil
}

type SourceFileDocument struct {
	Uri DocumentURI
}

func (d *SourceFileDocument) EqualTo(uri DocumentURI) bool {
	return d.Uri == uri
}

func (d *SourceFileDocument) String() string {
	return fmt.Sprintf("document %v", d.Uri)
}

type SourceFileReference struct {
	Range    Range
	Document *SourceFileDocument
}

func MakeSourceFileDocument(uri string) *SourceFileDocument {
	return MakeSourceFileDocumentFromURI(MakeDocumentURI(uri))
}

func MakeSourceFileDocumentFromLocalPath(localPath string) *SourceFileDocument {
	return MakeSourceFileDocumentFromURI(MakeDocumentURIFromLocalPath(localPath))
}

func MakeSourceFileDocumentFromURI(uri DocumentURI) *SourceFileDocument {
	if strings.Contains(string(uri), "document") {
		panic("stop")
	}
	return &SourceFileDocument{
		uri,
	}
}

func (s SourceFileReference) ToReferenceString() string {
	if s.Document == nil {
		return fmt.Sprintf("document is nil in sourcefilereference")
	}
	return fmt.Sprintf("%v:%d:%d:", s.Document.Uri, s.Range.start.line+1, s.Range.start.column+1)
}

func (s SourceFileReference) ToCompleteReferenceString() string {
	var uri DocumentURI
	if s.Document != nil {
		uri = s.Document.Uri
	}

	localPath, err := uri.ToLocalFilePath()
	if err != nil {
		localPath = string(uri)
	}
	if localPath == "" {
		localPath = string(uri)
	}
	return fmt.Sprintf("%v:%d:%d - %d:%d:", localPath, s.Range.start.line+1, s.Range.start.column+1, s.Range.end.line+1, s.Range.end.column+1)
}

func (s SourceFileReference) ToStandardReferenceString() string {
	var uri DocumentURI
	if s.Document != nil {
		uri = s.Document.Uri
	}

	localPath, err := uri.ToLocalFilePath()
	if err != nil {
		localPath = string(uri)
	}
	if localPath == "" {
		localPath = string(uri)
	}
	return fmt.Sprintf("%v:%d:%d:", localPath, s.Range.start.line+1, s.Range.start.column+1)
}

func (s SourceFileReference) String() string {
	return s.ToReferenceString()
}

func MakeSourceFileReferenceFromString(uri string, tokenRange Range) SourceFileReference {
	return SourceFileReference{
		Range:    tokenRange,
		Document: MakeSourceFileDocument(uri),
	}
}

func MakeSourceFileReference(uri *SourceFileDocument, tokenRange Range) SourceFileReference {
	return SourceFileReference{
		Range:    tokenRange,
		Document: uri,
	}
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
		Document: start.Document,
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

type SameLineRange struct {
	start  Position
	length int
}

func MakeSameLineRange(start Position, length int) SameLineRange {
	if length == 0 {
		panic("how is this possible")
	}
	return SameLineRange{
		start:  start,
		length: length,
	}
}

func (s SameLineRange) String() string {
	return fmt.Sprintf("[sameline pos:%v length:%v]", s.start, s.length)
}

type Range struct {
	start Position
	end   Position
}

func MakeRange(start Position, end Position) Range {
	return Range{start: start, end: end}
}

func (p Range) SmallerThan(other Range) bool {
	diffLineOther := other.end.line - other.start.line
	diffLine := p.end.line - p.start.line
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

func (p Range) IsAfter(other Range) bool {
	return (p.start.line > other.end.line) || ((p.start.line == other.end.line) && p.start.column > other.end.column)
}

func (p Range) IsEqual(other Range) bool {
	return p.start.line == other.start.line && p.start.column == other.start.column && p.end.line == other.end.line && p.end.column == other.end.column
}

func MakeInclusiveRange(start Range, end Range) Range {
	return Range{
		start: start.Start(),
		end:   end.End(),
	}
}

func NewPositionLength(start Position, octetCountIncludingWhitespace int) Range {
	return Range{start: start, end: Position{
		line:        start.line,
		column:      start.column + octetCountIncludingWhitespace - 1,
		octetOffset: start.octetOffset + octetCountIncludingWhitespace,
	}}
}

func NewPositionLengthFromEndPosition(start Position, endPosition Position) Range {
	if endPosition.octetOffset < 0 {
		panic("octet offset is wrong")
	}
	return Range{start: start, end: endPosition}
}

func (p Range) RuneWidth() int {
	return p.end.octetOffset - p.start.octetOffset + 1
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

func (p Range) ContainsRange(other Range) bool {
	return other.Start().IsOnOrAfter(other.Start()) && p.end.IsOnOrAfter(other.End())
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

func (p Range) OctetCount() int {
	return p.end.octetOffset - p.start.octetOffset + 1
}

func (p Range) String() string {
	return fmt.Sprintf("[%v to %v (%v)] ", p.start, p.end, p.OctetCount())
}
