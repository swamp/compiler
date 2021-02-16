/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

type CommentBlock struct {
	Comments []CommentToken
}

func MakeCommentBlock(comments []CommentToken) CommentBlock {
	return CommentBlock{Comments: comments}
}

type CommentToken struct {
	PositionLength
	RawString        string
	CommentString    string
	ForDocumentation bool
}

func (s CommentToken) Type() Type {
	return CommentConstant
}

func (s CommentToken) Raw() string {
	return s.RawString
}

func (s CommentToken) Text() string {
	return s.CommentString
}

// CommentToken :
type MultiLineCommentToken struct {
	CommentToken
}

func NewMultiLineCommentToken(raw string, text string, forDocumentation bool, position PositionLength) MultiLineCommentToken {
	return MultiLineCommentToken{CommentToken: CommentToken{RawString: raw, CommentString: text, PositionLength: position, ForDocumentation: forDocumentation}}
}

func (s MultiLineCommentToken) String() string {
	return fmt.Sprintf("[comment:%s]", s.CommentString)
}
