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

type RecordConstructorFromRecord struct {
	recordType           *dectype.RecordAtom
	record               *RecordLiteral
	recordAliasReference *dectype.AliasReference
	astConstructorCall   *ast.ConstructorCall
}

func NewRecordConstructorFromRecord(astConstructorCall *ast.ConstructorCall, recordAliasReference *dectype.AliasReference, recordType *dectype.RecordAtom, record *RecordLiteral) *RecordConstructorFromRecord {
	return &RecordConstructorFromRecord{astConstructorCall: astConstructorCall, recordAliasReference: recordAliasReference, record: record, recordType: recordType}
}

func (c *RecordConstructorFromRecord) Type() dtype.Type {
	return c.recordAliasReference
}

func (c *RecordConstructorFromRecord) Expression() Expression {
	return c.record
}

func (c *RecordConstructorFromRecord) HumanReadable() string {
	return "Record Constructor"
}

func (c *RecordConstructorFromRecord) NamedTypeReference() *dectype.NamedDefinitionTypeReference {
	return c.recordAliasReference.NameReference()
}

func (c *RecordConstructorFromRecord) String() string {
	return fmt.Sprintf("[RecordConstructorRecord %v %v]", c.astConstructorCall, c.record)
}

func (c *RecordConstructorFromRecord) FetchPositionLength() token.SourceFileReference {
	return c.astConstructorCall.FetchPositionLength()
}
