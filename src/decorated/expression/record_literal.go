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

func (a ByAssignmentName) Len() int           { return len(a) }
func (a ByAssignmentName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAssignmentName) Less(i, j int) bool { return a[i].fieldName.Name() < a[j].fieldName.Name() }

type RecordLiteralAssignment struct {
	expression Expression
	index      int
	fieldName  *ast.VariableIdentifier
}

func NewRecordLiteralAssignment(index int, fieldName *ast.VariableIdentifier, expression Expression) *RecordLiteralAssignment {
	return &RecordLiteralAssignment{index: index, fieldName: fieldName, expression: expression}
}

func (a *RecordLiteralAssignment) String() string {
	return fmt.Sprintf("%v = %v", a.index, a.expression)
}

func (a *RecordLiteralAssignment) Index() int {
	return a.index
}

func (a *RecordLiteralAssignment) FieldName() *ast.VariableIdentifier {
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
