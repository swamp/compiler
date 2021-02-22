/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

// LineDelimiterToken :
type LineDelimiterToken struct {
	SourceFileReference
}

func NewLineDelimiter(position SourceFileReference) LineDelimiterToken {
	return LineDelimiterToken{SourceFileReference: position}
}

func (s LineDelimiterToken) String() string {
	return "LF"
}

func (s LineDelimiterToken) Type() Type {
	return NewLine
}

func (s LineDelimiterToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
