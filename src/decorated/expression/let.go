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

type LetVariable struct {
	name         *ast.VariableIdentifier `debug:"true"`
	variableType dtype.Type              `debug:"true"`
	references   []*LetVariableReference
	comment      *ast.MultilineComment
}

func (l *LetVariable) String() string {
	return fmt.Sprintf("[LetVar %v]", l.name)
}

func (l *LetVariable) HumanReadable() string {
	return "Let variable"
}

func (l *LetVariable) Type() dtype.Type {
	return l.variableType
}

func (l *LetVariable) Name() *ast.VariableIdentifier {
	return l.name
}

func (l *LetVariable) IsIgnore() bool {
	return l.name.IsIgnore()
}

func (l *LetVariable) AddReferee(ref *LetVariableReference) {
	l.references = append(l.references, ref)
}

func (l *LetVariable) WasReferenced() bool {
	return len(l.references) > 0
}

func (l *LetVariable) References() []*LetVariableReference {
	return l.references
}

func (l *LetVariable) Comment() *ast.MultilineComment {
	return l.comment
}

func NewLetVariable(name *ast.VariableIdentifier, variableType dtype.Type, comment *ast.MultilineComment) *LetVariable {
	return &LetVariable{
		name:         name,
		variableType: variableType,
		comment:      comment,
	}
}

func (l *LetVariable) FetchPositionLength() token.SourceFileReference {
	return l.name.FetchPositionLength()
}

type LetAssignment struct {
	expression       Expression     `debug:"true"`
	letVariables     []*LetVariable `debug:"true"`
	astLetAssignment ast.LetAssignment
	inclusive        token.SourceFileReference
}

func NewLetAssignment(astLetAssignment ast.LetAssignment, letVariables []*LetVariable, expression Expression) *LetAssignment {
	inclusive := token.MakeInclusiveSourceFileReference(letVariables[0].FetchPositionLength(), expression.FetchPositionLength())
	return &LetAssignment{letVariables: letVariables, expression: expression, astLetAssignment: astLetAssignment, inclusive: inclusive}
}

func (l *LetAssignment) String() string {
	return fmt.Sprintf("[LetAssign %v = %v]", l.letVariables, l.expression)
}

func (l *LetAssignment) LetVariables() []*LetVariable {
	return l.letVariables
}

func (l *LetAssignment) WasRecordDestructuring() bool {
	return l.astLetAssignment.WasRecordDestructuring()
}

func (l *LetAssignment) Expression() Expression {
	return l.expression
}

func (l *LetAssignment) Type() dtype.Type {
	return l.expression.Type()
}

func (l *LetAssignment) FetchPositionLength() token.SourceFileReference {
	return l.inclusive
}

type Let struct {
	assignments []*LetAssignment `debug:"true"`
	consequence Expression       `debug:"true"`
	inclusive   token.SourceFileReference
	astLet      *ast.Let
}

func NewLet(let *ast.Let, assignments []*LetAssignment, consequence Expression) *Let {
	inclusive := token.MakeInclusiveSourceFileReference(let.FetchPositionLength(), consequence.FetchPositionLength())
	return &Let{assignments: assignments, consequence: consequence, inclusive: inclusive, astLet: let}
}

func (l *Let) Assignments() []*LetAssignment {
	return l.assignments
}

func (l *Let) Consequence() Expression {
	return l.consequence
}

func (l *Let) AstLet() *ast.Let {
	return l.astLet
}

func (l *Let) Type() dtype.Type {
	return l.consequence.Type()
}

func (l *Let) String() string {
	return fmt.Sprintf("[Let %v in %v]", l.assignments, l.consequence)
}

func (l *Let) FetchPositionLength() token.SourceFileReference {
	return l.inclusive
}
