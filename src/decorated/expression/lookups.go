/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type LookupField struct {
	reference *RecordTypeFieldReference
}

func (l LookupField) String() string {
	return fmt.Sprintf("[lookup %v]", l.reference.RecordTypeField())
}

func (l LookupField) Index() int {
	return l.reference.recordTypeField.Index()
}

func (l LookupField) Identifier() *ast.VariableIdentifier {
	return l.reference.ident
}

func (l LookupField) RecordTypeFieldReference() *RecordTypeFieldReference {
	return l.reference
}

func (l LookupField) MemoryOffset() dectype.MemoryOffset {
	return l.reference.recordTypeField.MemoryOffset()
}

func (l LookupField) MemorySize() dectype.MemorySize {
	return l.reference.recordTypeField.MemorySize()
}

func NewLookupField(reference *RecordTypeFieldReference) LookupField {
	return LookupField{reference: reference}
}

func (l LookupField) FetchPositionLength() token.SourceFileReference {
	return l.reference.FetchPositionLength()
}

type RecordLookups struct {
	ExpressionNode
	expressionToRecord Expression
	lookupFields       []LookupField
}

func NewRecordLookups(expressionToRecord Expression, lookupFields []LookupField) *RecordLookups {
	l := &RecordLookups{expressionToRecord: expressionToRecord, lookupFields: lookupFields}
	count := len(lookupFields)
	l.decoratedType = lookupFields[count-1].reference.recordTypeField.Type()
	return l
}

func (l *RecordLookups) Expression() Expression {
	return l.expressionToRecord
}

func (l *RecordLookups) LookupFields() []LookupField {
	return l.lookupFields
}

func (l *RecordLookups) String() string {
	return fmt.Sprintf("[lookups %v %v]", l.expressionToRecord, l.lookupFields)
}

func (l *RecordLookups) FetchPositionLength() token.SourceFileReference {
	return l.expressionToRecord.FetchPositionLength()
}
