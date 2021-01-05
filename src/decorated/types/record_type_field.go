/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
)

type RecordField struct {
	index     int
	name      *ast.VariableIdentifier
	fieldType dtype.Type

}

func NewRecordField(name *ast.VariableIdentifier, fieldType dtype.Type) *RecordField {
	return &RecordField{index: -1, name: name, fieldType: fieldType}
}

func (s *RecordField) SetIndexBySorter(index int) {
	s.index = index
}

func (s *RecordField) Index() int {
	if s.index == -1 {
		panic("you can not read index if it isn't set properly")
	}
	return s.index
}

func (s *RecordField) Name() string {
	return s.name.Name()
}

func (s *RecordField) VariableIdentifier() *ast.VariableIdentifier {
	return s.name
}

func (s *RecordField) Type() dtype.Type {
	return s.fieldType
}

func (s *RecordField) String() string {
	return fmt.Sprintf("[record-type-field %v %v (%v)]", s.name.Name(), s.fieldType, s.index)
}

func (s *RecordField) ShortString() string {
	return fmt.Sprintf("[record-field %v %v]", s.name.Name(), s.fieldType.ShortString())
}

func (s *RecordField) HumanReadable() string {
	return fmt.Sprintf("%v:%v", s.name.Name(), s.fieldType.HumanReadable())
}

func (s *RecordField) DecoratedName() string {
	return fmt.Sprintf("%v:%v", s.name.Name(), s.fieldType.DecoratedName())
}
