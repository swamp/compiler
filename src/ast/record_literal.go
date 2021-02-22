/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"
	"sort"

	"github.com/swamp/compiler/src/token"
)

type ByAssignmentName []*RecordLiteralFieldAssignment

func (a ByAssignmentName) Len() int           { return len(a) }
func (a ByAssignmentName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAssignmentName) Less(i, j int) bool { return a[i].identifier.Name() < a[j].identifier.Name() }

type RecordLiteral struct {
	assignments       []*RecordLiteralFieldAssignment
	sortedAssignments []*RecordLiteralFieldAssignment
	templateRecord    *VariableIdentifier
	parenToken        token.ParenToken
}

func NewRecordLiteral(parenToken token.ParenToken, templateRecord *VariableIdentifier, assignments []*RecordLiteralFieldAssignment) *RecordLiteral {
	sortedAssignments := make([]*RecordLiteralFieldAssignment, len(assignments))
	copy(sortedAssignments, assignments)
	sort.Sort(ByAssignmentName(sortedAssignments))
	return &RecordLiteral{parenToken: parenToken, templateRecord: templateRecord, assignments: assignments, sortedAssignments: sortedAssignments}
}

func (i *RecordLiteral) String() string {
	if i.templateRecord != nil {
		return fmt.Sprintf("[record-literal: %v (%v)]", i.assignments, i.templateRecord)
	}
	return fmt.Sprintf("[record-literal: %v]", i.assignments)
}

func (i *RecordLiteral) DebugString() string {
	return fmt.Sprintf("[record-literal]")
}

func (i *RecordLiteral) FetchPositionLength() token.Range {
	return i.parenToken.FetchPositionLength()
}

func (i *RecordLiteral) SortedAssignments() []*RecordLiteralFieldAssignment {
	return i.sortedAssignments
}

func (i *RecordLiteral) ParseOrderedAssignments() []*RecordLiteralFieldAssignment {
	return i.assignments
}

func (i *RecordLiteral) TemplateRecord() *VariableIdentifier {
	return i.templateRecord
}
