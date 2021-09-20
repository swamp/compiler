/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
)

type TupleTypeField struct {
	index        int
	fieldType    dtype.Type
	memoryOffset MemoryOffset
	memorySize   MemorySize
	memoryAlign  MemoryAlign
}

func NewTupleTypeField(index int, fieldType dtype.Type) *TupleTypeField {
	return &TupleTypeField{index: index, fieldType: fieldType}
}

func (s *TupleTypeField) SetIndexBySorter(index int) {
	s.index = index
}

func (s *TupleTypeField) MemoryOffset() MemoryOffset {
	return s.memoryOffset
}

func (s *TupleTypeField) MemorySize() MemorySize {
	return s.memorySize
}

func (s *TupleTypeField) Index() int {
	if s.index == -1 {
		panic("you can not read index if it isn't set properly")
	}
	return s.index
}

func (s *TupleTypeField) Type() dtype.Type {
	return s.fieldType
}

func (s *TupleTypeField) String() string {
	return fmt.Sprintf("[tuple-type-field %v (%v)]", s.fieldType, s.index)
}

func (s *TupleTypeField) HumanReadable() string {
	return fmt.Sprintf("%v", s.fieldType.HumanReadable())
}
