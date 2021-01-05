/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type LetAssignment struct {
	expression DecoratedExpression
	name       *ast.VariableIdentifier
}

func NewLetAssignment(name *ast.VariableIdentifier, expression DecoratedExpression) *LetAssignment {
	return &LetAssignment{name: name, expression: expression}
}

func (l *LetAssignment) String() string {
	return fmt.Sprintf("[letassign %v = %v]", l.name, l.expression)
}

func (l *LetAssignment) Name() *ast.VariableIdentifier {
	return l.name
}

func (l *LetAssignment) Expression() DecoratedExpression {
	return l.expression
}

type Let struct {
	assignments []*LetAssignment
	consequence DecoratedExpression
}

func NewLet(assignments []*LetAssignment, consequence DecoratedExpression) *Let {
	return &Let{assignments: assignments, consequence: consequence}
}

func (l *Let) Assignments() []*LetAssignment {
	return l.assignments
}

func (l *Let) Consequence() DecoratedExpression {
	return l.consequence
}

func (l *Let) Type() dtype.Type {
	return l.consequence.Type()
}

func (l *Let) String() string {
	return fmt.Sprintf("[let %v in %v]", l.assignments, l.consequence)
}

func (l *Let) FetchPositionAndLength() token.PositionLength {
	return l.consequence.FetchPositionAndLength()
}
