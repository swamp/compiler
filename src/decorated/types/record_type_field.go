/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type RecordFieldName struct {
	name *ast.VariableIdentifier
}

func NewRecordFieldName(identifier *ast.VariableIdentifier) *RecordFieldName {
	return &RecordFieldName{name: identifier}
}

func (r *RecordFieldName) Name() *ast.VariableIdentifier {
	return r.name
}

func (r *RecordFieldName) String() string {
	return r.name.String()
}

func (r *RecordFieldName) FetchPositionLength() token.SourceFileReference {
	return r.name.FetchPositionLength()
}

func (r *RecordFieldName) HumanReadable() string {
	return "type field name"
}

type (
	MemoryOffset uint32
	MemorySize   uint32
	MemoryAlign  uint32
)

type RecordField struct {
	index           int
	memoryOffset    MemoryOffset
	memorySize      MemorySize
	name            *RecordFieldName
	fieldType       dtype.Type
	recordTypeField *ast.RecordTypeField
}

func NewRecordField(fieldName *RecordFieldName, recordTypeField *ast.RecordTypeField, fieldType dtype.Type) *RecordField {
	_, wasPrimitive := fieldType.(*PrimitiveAtom)
	if wasPrimitive {
		panic(fmt.Errorf("use type reference, not primitive directly %v = %v %v %T", fieldName, recordTypeField, fieldType, fieldType))
	}

	return &RecordField{index: -1, name: fieldName, fieldType: fieldType, recordTypeField: recordTypeField}
}

func (s *RecordField) SetIndexBySorter(index int) {
	s.index = index
}

func (s *RecordField) MemoryOffset() MemoryOffset {
	return s.memoryOffset
}

func (s *RecordField) MemorySize() MemorySize {
	return s.memorySize
}

func (s *RecordField) Index() int {
	if s.index == -1 {
		panic("you can not read index if it isn't set properly")
	}
	return s.index
}

func (s *RecordField) Name() string {
	return s.name.Name().Name()
}

func (s *RecordField) VariableIdentifier() *ast.VariableIdentifier {
	return s.name.Name()
}

func (s *RecordField) AstRecordTypeField() *ast.RecordTypeField {
	return s.recordTypeField
}

func (s *RecordField) FieldName() *RecordFieldName {
	return s.name
}

func (s *RecordField) Type() dtype.Type {
	return s.fieldType
}

func (s *RecordField) String() string {
	return fmt.Sprintf("[Field %v %v (%v)]", s.name.Name(), s.fieldType, s.index)
}

func (s *RecordField) HumanReadable() string {
	return fmt.Sprintf("%v : %v", s.name.Name().Name(), s.fieldType.HumanReadable())
}
