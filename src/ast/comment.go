/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

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
	return i.commentToken.Text()
}

func (i *MultilineComment) String() string {
	return fmt.Sprintf("[multilinecomment '%v']", i.commentToken.Text())
}

func (i *MultilineComment) PositionLength() token.Range {
	return i.commentToken.Range
}

func (i *MultilineComment) DebugString() string {
	return i.String()
}

func (i *MultilineComment) Token() token.MultiLineCommentToken {
	return i.commentToken
}

func (i *MultilineComment) FetchPositionLength() token.SourceFileReference {
	return i.commentToken.SourceFileReference
}
