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
	fieldName *ast.VariableIdentifier
}

func NewRecordLiteralField(fieldName *ast.VariableIdentifier) *RecordLiteralField {
	return &RecordLiteralField{fieldName: fieldName}
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
	return "Record field identifier"
}

type RecordLiteralAssignment struct {
	expression Expression
	index      int
	fieldName  *RecordLiteralField
}

func NewRecordLiteralAssignment(index int, fieldName *RecordLiteralField, expression Expression) *RecordLiteralAssignment {
	return &RecordLiteralAssignment{index: index, fieldName: fieldName, expression: expression}
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
	sortedAssignments       []*RecordLiteralAssignment
	parseOrderedAssignments []*RecordLiteralAssignment
	recordTemplate          Expression
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
	return &RecordLiteral{
		t: t, recordTemplate: recordTemplate,
		parseOrderedAssignments: parseOrderedAssignments,
		sortedAssignments:       sortedAssignments,
	}
}

func (c *RecordLiteral) Type() dtype.Type {
	if c.recordTemplate != nil {
		return c.recordTemplate.Type()
	}
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
	return fmt.Sprintf("[record-literal %v %v]", c.t, c.sortedAssignments)
}

func (c *RecordLiteral) FetchPositionLength() token.SourceFileReference {
	return token.SourceFileReference{}
}
