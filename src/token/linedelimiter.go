/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

// LineDelimiterToken :
type LineDelimiterToken struct {
	PositionLength
}

func NewLineDelimiter(position PositionLength) LineDelimiterToken {
	return LineDelimiterToken{PositionLength: position}
}

func (s LineDelimiterToken) String() string {
	return "LF"
}

func (s LineDelimiterToken) Type() Type {
	return NewLine
}
