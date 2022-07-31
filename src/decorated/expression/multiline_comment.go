/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type Comment interface {
	MarkdownString() string
}

type MultilineComment struct {
	comment *ast.MultilineComment
}

func NewMultilineComment(comment *ast.MultilineComment) *MultilineComment {
	return &MultilineComment{comment: comment}
}

func (m *MultilineComment) String() string {
	return "multiline"
}

func (m *MultilineComment) MarkdownString() string {
	return m.comment.Value()
}

func (m *MultilineComment) StatementString() string {
	return "multiline"
}

func (m *MultilineComment) DebugString() string {
	return "multiline"
}

func (m *MultilineComment) AstComment() *ast.MultilineComment {
	return m.comment
}

func (m *MultilineComment) FetchPositionLength() token.SourceFileReference {
	return m.comment.FetchPositionLength()
}
