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
	name         *ast.VariableIdentifier
	variableType dtype.Type
	references   []*LetVariableReference
}

func (l *LetVariable) String() string {
	return fmt.Sprintf("[letvariable %v]", l.name)
}

func (l *LetVariable) Type() dtype.Type {
	return l.variableType
}

func (l *LetVariable) Name() *ast.VariableIdentifier {
	return l.name
}

func (l *LetVariable) AddReferee(ref *LetVariableReference) {
	l.references = append(l.references, ref)
}

func (l *LetVariable) References() []*LetVariableReference {
	return l.references
}

func NewLetVariable(name *ast.VariableIdentifier, variableType dtype.Type) *LetVariable {
	return &LetVariable{
		name:         name,
		variableType: variableType,
	}
}

func (l *LetVariable) FetchPositionLength() token.SourceFileReference {
	return l.name.FetchPositionLength()
}

type LetAssignment struct {
	expression  Expression
	letVariable *LetVariable
	inclusive   token.SourceFileReference
	references  []*LetVariableReference
}

func NewLetAssignment(name *ast.VariableIdentifier, expression Expression) *LetAssignment {
	letVar := NewLetVariable(name, expression.Type())
	inclusive := token.MakeInclusiveSourceFileReference(name.FetchPositionLength(), expression.FetchPositionLength())
	return &LetAssignment{letVariable: letVar, expression: expression, inclusive: inclusive}
}

func (l *LetAssignment) String() string {
	return fmt.Sprintf("[letassign %v = %v]", l.letVariable, l.expression)
}

func (l *LetAssignment) LetVariable() *LetVariable {
	return l.letVariable
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

func (l *LetAssignment) AddReferee(ref *LetVariableReference) {
	l.letVariable.AddReferee(ref)
	l.references = append(l.references, ref)
}

func (l *LetAssignment) References() []*LetVariableReference {
	return l.references
}

type Let struct {
	assignments []*LetAssignment
	consequence Expression
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
	return fmt.Sprintf("[let %v in %v]", l.assignments, l.consequence)
}

func (l *Let) FetchPositionLength() token.SourceFileReference {
	return l.inclusive
}
