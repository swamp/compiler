/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

type IndentationReport struct {
	PreviousIndentationSpaces int
	PreviousCloseIndentation  int
	PreviousExactIndentation  int
	IndentationSpaces         int
	SpacesUntilMaybeNewline   int
	ExactIndentation          int
	CloseIndentation          int
	NewLineCount              int
	EndOfFile                 bool
	StartPos                  PositionToken
	PositionLength            PositionLength
	Comments                  CommentBlock
	TrailingSpacesFound       bool
}
