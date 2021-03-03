package decorated

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type SingleLineComment struct {
	comment *ast.SingleLineComment
}

func NewSingleLineComment(comment *ast.SingleLineComment) *SingleLineComment {
	return &SingleLineComment{comment: comment}
}

func (m *SingleLineComment) String() string {
	return "singleline"
}

func (m *SingleLineComment) MarkdownString() string {
	return m.comment.Value()
}

func (m *SingleLineComment) StatementString() string {
	return "singleline"
}

func (m *SingleLineComment) FetchPositionLength() token.SourceFileReference {
	return m.comment.FetchPositionLength()
}
