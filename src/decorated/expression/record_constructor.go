/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type RecordConstructor struct {
	arguments            []*RecordLiteralAssignment
	parseOrderArguments  []Expression
	recordAliasReference *dectype.AliasReference
	recordType           *dectype.RecordAtom
}

// typeIdentifier ast.TypeReferenceScopedOrNormal
func NewRecordConstructor(recordAliasReference *dectype.AliasReference, recordType *dectype.RecordAtom, arguments []*RecordLiteralAssignment, parseOrderArguments []Expression) *RecordConstructor {
	if recordAliasReference == nil {
		panic("can not be nil")
	}
	return &RecordConstructor{recordAliasReference: recordAliasReference, arguments: arguments, parseOrderArguments: parseOrderArguments, recordType: recordType}
}

func (c *RecordConstructor) SortedAssignments() []*RecordLiteralAssignment {
	return c.arguments
}

func (c *RecordConstructor) NamedTypeReference() *dectype.NamedDefinitionTypeReference {
	return c.recordAliasReference.NamedTypeRef()
}

func (c *RecordConstructor) ParseOrderArguments() []Expression {
	return c.parseOrderArguments
}

func (c *RecordConstructor) Type() dtype.Type {
	return c.recordType
}

func (c *RecordConstructor) RecordType() *dectype.RecordAtom {
	return c.recordType
}

func (c *RecordConstructor) String() string {
	return fmt.Sprintf("[record-constructor %v %v]", c.recordAliasReference, c.arguments)
}

func (c *RecordConstructor) HumanReadable() string {
	return "Record Constructor"
}

func (c *RecordConstructor) FetchPositionLength() token.SourceFileReference {
	return c.recordAliasReference.FetchPositionLength()
}
