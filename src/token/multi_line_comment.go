/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
	"strings"
)

type Comment interface {
	Raw() string
	Value() string
	FetchPositionLength() SourceFileReference
}

type CommentBlock struct {
	Comments []Comment
}

func (c CommentBlock) LastComment() Comment {
	if len(c.Comments) == 0 {
		return nil
	}

	return c.Comments[len(c.Comments)-1]
}

func MakeCommentBlock(comments []Comment) CommentBlock {
	return CommentBlock{Comments: comments}
}

type MultiLineCommentPart struct {
	SourceFileReference
	RawString     string
	CommentString string
}

type CommentToken struct {
	SourceFileReference
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

func (s CommentToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}

// MultiLineCommentToken :
type MultiLineCommentToken struct {
	parts            []MultiLineCommentPart
	ForDocumentation bool
}

func NewMultiLineCommentToken(parts []MultiLineCommentPart, forDocumentation bool) MultiLineCommentToken {
	return MultiLineCommentToken{parts: parts, ForDocumentation: forDocumentation}
}

func (s MultiLineCommentToken) String() string {
	return fmt.Sprintf("[comment:%s]", s.parts)
}

func (s MultiLineCommentToken) Raw() string {
	return fmt.Sprintf("[comment:%s]", s.parts)
}

func (s MultiLineCommentToken) Type() Type {
	return CommentConstant
}

func (s MultiLineCommentToken) Parts() []MultiLineCommentPart {
	return s.parts
}

func (s MultiLineCommentToken) Value() string {
	str := ""
	for _, part := range s.parts {
		upcomingString := strings.TrimSpace(part.CommentString)
		if len(str) == 0 && len(upcomingString) == 0 {
			continue
		}
		if len(str) > 0 {
			str += "\n"
		}
		str += part.CommentString
	}

	return str
}

func (s MultiLineCommentToken) FetchPositionLength() SourceFileReference {
	return MakeInclusiveSourceFileReference(s.parts[0].SourceFileReference, s.parts[len(s.parts)-1].SourceFileReference)
}
