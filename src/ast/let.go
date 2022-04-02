/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type LetAssignment struct {
	identifiers            []*VariableIdentifier
	expression             Expression
	inclusive              token.SourceFileReference
	comment                *MultilineComment
	wasRecordDestructuring bool
}

func NewLetAssignment(wasRecordDestructuring bool, identifiers []*VariableIdentifier, expression Expression, comment *MultilineComment) LetAssignment {
	inclusive := token.MakeInclusiveSourceFileReference(identifiers[0].FetchPositionLength(), expression.FetchPositionLength())
	return LetAssignment{wasRecordDestructuring: wasRecordDestructuring, identifiers: identifiers, expression: expression, inclusive: inclusive, comment: comment}
}

func (l LetAssignment) Identifiers() []*VariableIdentifier {
	return l.identifiers
}

func (l LetAssignment) WasRecordDestructuring() bool {
	return l.wasRecordDestructuring
}

func (l LetAssignment) Expression() Expression {
	return l.expression
}

func (l LetAssignment) String() string {
	return fmt.Sprintf("[letassign %v = %v]", l.identifiers, l.expression)
}

func (l LetAssignment) CommentBlock() *MultilineComment {
	return l.comment
}

func (l LetAssignment) FetchPositionLength() token.SourceFileReference {
	return l.inclusive
}

type Let struct {
	assignments         []LetAssignment
	consequence         Expression
	keyword             token.Keyword
	inKeyword           token.Keyword
	sourceFileReference token.SourceFileReference
}

func NewLet(keyword token.Keyword, inKeyword token.Keyword, assignments []LetAssignment, consequence Expression) *Let {
	sourceFileReference := token.MakeInclusiveSourceFileReference(keyword.FetchPositionLength(), consequence.FetchPositionLength())
	return &Let{keyword: keyword, inKeyword: inKeyword, assignments: assignments, consequence: consequence, sourceFileReference: sourceFileReference}
}

func (i *Let) Keyword() token.Keyword {
	return i.keyword
}

func (i *Let) InKeyword() token.Keyword {
	return i.inKeyword
}

func (i *Let) Assignments() []LetAssignment {
	return i.assignments
}

func (i *Let) Consequence() Expression {
	return i.consequence
}

func (i *Let) FetchPositionLength() token.SourceFileReference {
	return i.sourceFileReference
}

func (i *Let) String() string {
	return fmt.Sprintf("[let: %v in %v]", i.assignments, i.consequence)
}

func (i *Let) DebugString() string {
	return "[let]"
}
