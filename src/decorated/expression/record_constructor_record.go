/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type RecordConstructorRecord struct {
	recordType           *dectype.RecordAtom
	record               *RecordLiteral
	recordAliasReference *dectype.AliasReference
	astConstructorCall   *ast.ConstructorCall
}

func NewRecordConstructorRecord(astConstructorCall *ast.ConstructorCall, recordAliasReference *dectype.AliasReference, recordType *dectype.RecordAtom, record *RecordLiteral) *RecordConstructorRecord {
	return &RecordConstructorRecord{astConstructorCall: astConstructorCall, recordAliasReference: recordAliasReference, record: record, recordType: recordType}
}

func (c *RecordConstructorRecord) Type() dtype.Type {
	return c.recordType
}

func (c *RecordConstructorRecord) Expression() Expression {
	return c.record
}

func (c *RecordConstructorRecord) NamedTypeReference() *dectype.NamedDefinitionTypeReference {
	return c.recordAliasReference.NameReference()
}

func (c *RecordConstructorRecord) String() string {
	return fmt.Sprintf("[record-constructor-record %v %v]", c.astConstructorCall, c.record)
}

func (c *RecordConstructorRecord) FetchPositionLength() token.SourceFileReference {
	return c.astConstructorCall.FetchPositionLength()
}
