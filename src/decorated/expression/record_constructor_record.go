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
	//arguments         []DecoratedExpression
	recordType     *dectype.RecordAtom
	record         *RecordLiteral
	typeIdentifier *ast.TypeIdentifier
}

func NewRecordConstructorRecord(typeIdentifier *ast.TypeIdentifier, recordType *dectype.RecordAtom, record *RecordLiteral) *RecordConstructorRecord {
	return &RecordConstructorRecord{typeIdentifier: typeIdentifier, record: record, recordType: recordType}
}

func (c *RecordConstructorRecord) Type() dtype.Type {
	return c.recordType
}

func (c *RecordConstructorRecord) Expression() DecoratedExpression {
	return c.record
}

func (c *RecordConstructorRecord) String() string {
	return fmt.Sprintf("[record-constructor-record %v %v]", c.typeIdentifier, c.record)
}

func (c *RecordConstructorRecord) FetchPositionAndLength() token.PositionLength {
	return c.typeIdentifier.Symbol().FetchPositionLength()
}
