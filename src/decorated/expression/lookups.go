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
	structField      *dectype.RecordField
	lookupIdentifier *ast.VariableIdentifier
}

func (l LookupField) String() string {
	return fmt.Sprintf("[lookup %v]", l.structField)
}

func (l LookupField) Index() int {
	return l.structField.Index()
}

func (l LookupField) Identifier() *ast.VariableIdentifier {
	return l.lookupIdentifier
}

func NewLookupField(lookupIdentifier *ast.VariableIdentifier, structField *dectype.RecordField) LookupField {
	return LookupField{structField: structField, lookupIdentifier: lookupIdentifier}
}

func (l LookupField) FetchPositionLength() token.SourceFileReference {
	return l.lookupIdentifier.FetchPositionLength()
}

/*
type LookupVariable struct {
	name       *ast.VariableIdentifier
	lookupType dtype.Type
}

func (l LookupVariable) Identifier() *ast.VariableIdentifier {
	return l.name
}

func (l LookupVariable) String() string {
	return fmt.Sprintf("[lookupvar %v (%v)]", l.name, l.lookupType)
}

func (l LookupVariable) DecoratedExpression() dtype.Type {
	return l.lookupType
}

*/

type RecordLookups struct {
	ExpressionNode
	expressionToRecord Expression
	lookupFields       []LookupField
}

func NewRecordLookups(expressionToRecord Expression, lookupFields []LookupField) *RecordLookups {
	l := &RecordLookups{expressionToRecord: expressionToRecord, lookupFields: lookupFields}
	count := len(lookupFields)
	l.decoratedType = lookupFields[count-1].structField.Type()
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
	return token.SourceFileReference{}
}
