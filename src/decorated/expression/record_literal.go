/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"sort"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type ByAssignmentName []*RecordLiteralAssignment

func (a ByAssignmentName) Len() int      { return len(a) }
func (a ByAssignmentName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByAssignmentName) Less(i, j int) bool {
	return a[i].fieldName.fieldName.Name() < a[j].fieldName.fieldName.Name()
}

type RecordLiteralField struct {
	fieldName    *ast.VariableIdentifier `debug:"true"`
	inAssignment *RecordLiteralAssignment
}

func NewRecordLiteralField(fieldName *ast.VariableIdentifier) *RecordLiteralField {
	return &RecordLiteralField{fieldName: fieldName}
}

func (n *RecordLiteralField) SetInAssignment(assignment *RecordLiteralAssignment) {
	n.inAssignment = assignment
}

func (n *RecordLiteralField) Ident() *ast.VariableIdentifier {
	return n.fieldName
}

func (n *RecordLiteralField) FetchPositionLength() token.SourceFileReference {
	return n.fieldName.FetchPositionLength()
}

func (n *RecordLiteralField) String() string {
	return "record field name"
}

func (n *RecordLiteralField) HumanReadable() string {
	return "Record field functionParameter"
}

func (n *RecordLiteralField) InAssignment() *RecordLiteralAssignment {
	return n.inAssignment
}

func (n *RecordLiteralField) Type() dtype.Type {
	return n.inAssignment.expression.Type()
}

type RecordLiteralAssignment struct {
	fieldName  *RecordLiteralField `debug:"true"`
	expression Expression          `debug:"true"`
	index      int
}

func NewRecordLiteralAssignment(index int, fieldName *RecordLiteralField,
	expression Expression) *RecordLiteralAssignment {
	assignment := &RecordLiteralAssignment{index: index, fieldName: fieldName, expression: expression}
	fieldName.SetInAssignment(assignment)

	return assignment
}

func (a *RecordLiteralAssignment) String() string {
	return fmt.Sprintf("%v = %v", a.index, a.expression)
}

func (a *RecordLiteralAssignment) Index() int {
	return a.index
}

func (a *RecordLiteralAssignment) FieldName() *RecordLiteralField {
	return a.fieldName
}

func (a *RecordLiteralAssignment) Expression() Expression {
	return a.expression
}

type RecordLiteral struct {
	t                       *dectype.RecordAtom
	sortedAssignments       []*RecordLiteralAssignment `debug:"true"`
	parseOrderedAssignments []*RecordLiteralAssignment
	recordTemplate          Expression
	inclusive               token.SourceFileReference
}

func NewRecordLiteral(t *dectype.RecordAtom, recordTemplate Expression,
	parseOrderedAssignments []*RecordLiteralAssignment) *RecordLiteral {
	lastFoundIndex := 0

	sortedAssignments := make([]*RecordLiteralAssignment, len(parseOrderedAssignments))
	copy(sortedAssignments, parseOrderedAssignments)
	sort.Sort(ByAssignmentName(sortedAssignments))

	for _, assignment := range sortedAssignments {
		if assignment.index < lastFoundIndex {
			panic("sortedAssignments are not sorted")
		}
		lastFoundIndex = assignment.index
	}

	first := parseOrderedAssignments[0].FieldName().FetchPositionLength()
	if recordTemplate != nil {
		first = recordTemplate.FetchPositionLength()
	}

	inclusive := token.MakeInclusiveSourceFileReference(first,
		parseOrderedAssignments[len(parseOrderedAssignments)-1].expression.FetchPositionLength())

	return &RecordLiteral{
		t: t, recordTemplate: recordTemplate,
		parseOrderedAssignments: parseOrderedAssignments,
		sortedAssignments:       sortedAssignments,
		inclusive:               inclusive,
	}
}

func (c *RecordLiteral) Type() dtype.Type {
	if c.recordTemplate != nil {
		return c.recordTemplate.Type()
	}
	return c.t
}

func (c *RecordLiteral) RecordType() *dectype.RecordAtom {
	return c.t
}

func (c *RecordLiteral) SortedAssignments() []*RecordLiteralAssignment {
	return c.sortedAssignments
}

func (c *RecordLiteral) ParseOrderedAssignments() []*RecordLiteralAssignment {
	return c.parseOrderedAssignments
}

func (c *RecordLiteral) RecordTemplate() Expression {
	return c.recordTemplate
}

func (c *RecordLiteral) String() string {
	return fmt.Sprintf("[RecordLiteral %v %v]", c.t, c.sortedAssignments)
}

func (c *RecordLiteral) FetchPositionLength() token.SourceFileReference {
	return c.inclusive
}
