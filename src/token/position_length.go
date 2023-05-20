/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"unicode"
)

type DocumentURI string

func MakeDocumentURI(s string) DocumentURI {
	if !strings.HasPrefix(s, "file://") {
		panic("illegal uri")
	}

	if strings.TrimSpace(s) == "" {
		panic("not allowed to be zero string")
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

func NewInternalSourceFileDocument() *SourceFileDocument {
	return MakeSourceFileDocumentFromLocalPath("")
}

type SourceFileDocument struct {
	Uri DocumentURI
}

func (d *SourceFileDocument) EqualTo(uri DocumentURI) bool {
	return d.Uri == uri
}

func (d *SourceFileDocument) String() string {
	return fmt.Sprintf("[document '%v']", d.Uri)
}

type SourceFileReference struct {
	Range    Range
	Document *SourceFileDocument
}

func (s SourceFileReference) Verify() bool {
	r := s.Range

	if r.Position().line == 0 && r.Position().column == 0 {
		log.Printf("suspicious")
		return false
	}

	return true
}

func NewInternalSourceFileReference() SourceFileReference {
	return SourceFileReference{
		Range:    MakeRange(NewPositionTopLeft(), NewPositionTopLeft()),
		Document: NewInternalSourceFileDocument(),
	}
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
		panic(fmt.Errorf("document is nil in sourcefilereference"))
	}
	return fmt.Sprintf("%v:%d:%d:", s.Document.Uri, s.Range.start.line+1, s.Range.start.column+1)
}

func (s SourceFileReference) ToStartAndEndReferenceString() string {
	if s.Document == nil {
		panic(fmt.Errorf("document is nil in sourcefilereference %T"))
	}
	return fmt.Sprintf("%v:%d:%d (%d:%d)", s.Document.Uri, s.Range.start.line+1, s.Range.start.column+1, s.Range.end.line+1, s.Range.end.column+1)
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

func MakeSourceFileReference(uri *SourceFileDocument, tokenRange Range) SourceFileReference {
	return SourceFileReference{
		Range:    tokenRange,
		Document: uri,
	}
}

func MakeInclusiveSourceFileReference(start SourceFileReference, end SourceFileReference) SourceFileReference {
	tokenRange := MakeInclusiveRange(start.Range, end.Range)
	if start.Document == nil {
		panic("start document can not be nil")
	}

	if end.Document == nil {
		panic("end document can not be nil")
	}

	if !start.Document.EqualTo(end.Document.Uri) {
		//panic(fmt.Sprintf("start and end must come from same document. '%v' vs '%v'", start.Document, end.Document))
	}

	return SourceFileReference{
		Range:    tokenRange,
		Document: start.Document,
	}
}

func MakeInclusiveSourceFileReferenceFlipIfNeeded(start SourceFileReference, end SourceFileReference) SourceFileReference {
	if start.Range.IsAfter(end.Range) {
		return MakeInclusiveSourceFileReference(end, start)
	}
	return MakeInclusiveSourceFileReference(start, end)
}

func MakeInclusiveSourceFileReferenceWithoutCheck(start SourceFileReference, end SourceFileReference) SourceFileReference {
	tokenRange := MakeInclusiveRangeWithoutCheck(start.Range, end.Range)
	if start.Document == nil {
		panic("start document can not be nil")
	}

	if end.Document == nil {
		panic("end document can not be nil")
	}

	if !start.Document.EqualTo(end.Document.Uri) {
		//panic(fmt.Sprintf("start and end must come from same document. '%v' vs '%v'", start.Document, end.Document))
	}

	return SourceFileReference{
		Range:    tokenRange,
		Document: start.Document,
	}
}

type SourceFileReferenceProvider interface {
	FetchPositionLength() SourceFileReference
}

type SameLineRange struct {
	Position         Position
	Length           int
	LocalOctetOffset int
}

func (s SameLineRange) String() string {
	return fmt.Sprintf("[stringline %v->%v (local offset: %v)]", s.Position, s.Length, s.LocalOctetOffset)
}

func (s SameLineRange) LocalStringOffset() int {
	return s.LocalOctetOffset
}

func (s SameLineRange) EndLocalStringOffset() int {
	return s.LocalOctetOffset + s.Length - 1
}

func MakeSameLineRange(start Position, length int, localOffset int) SameLineRange {
	if length == 0 {
		panic("how is this possible")
	}
	return SameLineRange{
		Position:         start,
		Length:           length,
		LocalOctetOffset: localOffset,
	}
}

func RangeFromSameLineRanges(ranges []SameLineRange) Range {
	if len(ranges) == 0 {
		return Range{}
	}
	first := ranges[0]
	last := ranges[len(ranges)-1]

	return MakeRange(first.Position, last.Position.AddColumn(last.Length))
}

func RangeFromSingleSameLineRange(singleRange SameLineRange) Range {
	return MakeRange(singleRange.Position, singleRange.Position.AddColumn(singleRange.Length))
}

type Range struct {
	start Position
	end   Position
}

func MakeRange(start Position, end Position) Range {
	if !end.IsOnOrAfter(start) {
		panic(fmt.Errorf("wrong range to create %v : %v", start, end))
	}
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
	if !end.End().IsOnOrAfter(start.Start()) {
		panic(fmt.Errorf("wrong inclusive range %v : %v", start.Start(), end.End()))
	}
	return Range{
		start: start.Start(),
		end:   end.End(),
	}
}

func MakeInclusiveRangeWithoutCheck(start Range, end Range) Range {
	return Range{
		start: start.Start(),
		end:   end.End(),
	}
}

func NewPositionLength(start Position, octetCountIncludingWhitespace int) Range {
	return Range{start: start, end: Position{
		line:                            start.line,
		column:                          start.column + octetCountIncludingWhitespace - 1,
		originalOctetOffsetInSourceFile: start.originalOctetOffsetInSourceFile + octetCountIncludingWhitespace,
	}}
}

func NewPositionLengthFromEndPosition(start Position, endPosition Position) Range {
	if endPosition.originalOctetOffsetInSourceFile < 0 {
		panic("octet offset is wrong")
	}
	return Range{start: start, end: endPosition}
}

func (p Range) RuneWidth() int {
	return p.end.originalOctetOffsetInSourceFile - p.start.originalOctetOffsetInSourceFile + 1
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
	return other.Start().IsOnOrAfter(p.Start()) && p.end.IsOnOrAfter(other.End())
}

func (p Range) ContainsSameLineRanges(other []SameLineRange) bool {
	return true
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
	return p.end.originalOctetOffsetInSourceFile - p.start.originalOctetOffsetInSourceFile + 1
}

func (p Range) String() string {
	return fmt.Sprintf("[%v to %v (%v)] ", p.start.DebugString(), p.end.DebugString(), p.OctetCount())
}
