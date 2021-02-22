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
	identifier *VariableIdentifier
	expression Expression
}

func NewLetAssignment(identifier *VariableIdentifier, expression Expression) LetAssignment {
	return LetAssignment{identifier: identifier, expression: expression}
}

func (l LetAssignment) Identifier() *VariableIdentifier {
	return l.identifier
}

func (l LetAssignment) Expression() Expression {
	return l.expression
}

func (l LetAssignment) String() string {
	return fmt.Sprintf("[letassign %v = %v]", l.identifier, l.expression)
}

type Let struct {
	assignments []LetAssignment
	consequence Expression
}

func NewLet(assignments []LetAssignment, consequence Expression) *Let {
	return &Let{assignments: assignments, consequence: consequence}
}

func (i *Let) Assignments() []LetAssignment {
	return i.assignments
}

func (i *Let) Consequence() Expression {
	return i.consequence
}

func (i *Let) FetchPositionLength() token.Range {
	return i.consequence.FetchPositionLength()
}

func (i *Let) String() string {
	return fmt.Sprintf("[let: %v in %v]", i.assignments, i.consequence)
}

func (i *Let) DebugString() string {
	return fmt.Sprintf("[let]")
}
