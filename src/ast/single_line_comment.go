/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type SingleLineComment struct {
	commentToken token.SingleLineCommentToken
}

func NewSingleLineComment(commentToken token.SingleLineCommentToken) *SingleLineComment {
	return &SingleLineComment{commentToken: commentToken}
}

func (i *SingleLineComment) Value() string {
	return i.commentToken.Raw()
}

func (i *SingleLineComment) String() string {
	return fmt.Sprintf("[singlelinecomment '%v']", i.commentToken.Text())
}

func (i *SingleLineComment) PositionLength() token.Range {
	return i.commentToken.Range
}

func (i *SingleLineComment) DebugString() string {
	return i.String()
}

func (i *SingleLineComment) FetchPositionLength() token.SourceFileReference {
	return i.commentToken.SourceFileReference
}
