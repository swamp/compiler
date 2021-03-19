/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

func CommentBlockToAst(comments token.CommentBlock) []*MultilineComment {
	var previousComment []*MultilineComment
	for _, comment := range comments.Comments {
		multiLine, wasMultiLine := comment.(token.MultiLineCommentToken)
		if wasMultiLine {
			astMultilineComment := NewMultilineComment(multiLine)
			previousComment = append(previousComment, astMultilineComment)
		} else {
			panic("not sure")
		}
	}
	return previousComment
}

type LocalComment struct {
	Singleline *SingleLineComment
	Multiline  *MultilineComment
}

type MultilineComment struct {
	commentToken token.MultiLineCommentToken
}

func NewMultilineComment(commentToken token.MultiLineCommentToken) *MultilineComment {
	return &MultilineComment{commentToken: commentToken}
}

func (i *MultilineComment) Value() string {
	return i.commentToken.Value()
}

func (i *MultilineComment) String() string {
	return fmt.Sprintf("[multilinecomment '%v']", i.commentToken.Value())
}

func (i *MultilineComment) PositionLength() token.Range {
	return i.commentToken.FetchPositionLength().Range
}

func (i *MultilineComment) DebugString() string {
	return i.String()
}

func (i *MultilineComment) Token() token.MultiLineCommentToken {
	return i.commentToken
}

func (i *MultilineComment) FetchPositionLength() token.SourceFileReference {
	return i.commentToken.FetchPositionLength()
}
