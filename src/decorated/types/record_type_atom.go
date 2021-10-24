/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"sort"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type RecordAtom struct {
	nameToField       map[string]*RecordField
	parsedOrderFields []*RecordField
	sortedFields      []*RecordField
	genericTypes      []dtype.Type
	record            *ast.Record
	memorySize        MemorySize
	memoryAlign       MemoryAlign
}

func (s *RecordAtom) MemorySize() MemorySize {
	return s.memorySize
}

func (s *RecordAtom) MemoryAlignment() MemoryAlign {
	return s.memoryAlign
}

func (s *RecordAtom) GenericTypes() []dtype.Type {
	return s.genericTypes
}

func (s *RecordAtom) AstRecord() *ast.Record {
	return s.record
}

func (s *RecordAtom) String() string {
	return fmt.Sprintf("record-type %v%v]", s.sortedFields, s.genericTypes)
}

func (s *RecordAtom) HumanReadable() string {
	str := "{ "
	for index, field := range s.sortedFields {
		if index > 0 {
			str += ", "
		}
		str += field.HumanReadable()
	}
	str += " }"
	return str
}

func (s *RecordAtom) FetchPositionLength() token.SourceFileReference {
	return s.record.FetchPositionLength()
}

const (
	Sizeof64BitPointer  MemorySize  = 8
	Alignof64BitPointer MemoryAlign = 8
	SizeofSwampInt      MemorySize  = 4
	SizeofSwampRune     MemorySize  = 4
	SizeofSwampBool     MemorySize  = 1

	AlignOfSwampBool = MemoryAlign(SizeofSwampBool)
	AlignOfSwampRune = MemoryAlign(SizeofSwampRune)
	AlignOfSwampInt  = MemoryAlign(SizeofSwampInt)
)

func GetMemorySizeAndAlignmentInternal(p dtype.Type) (MemorySize, MemoryAlign) {
	if p == nil {
		panic(fmt.Errorf("nil is not allowed"))
	}
	unaliased := UnaliasWithResolveInvoker(p)
	switch t := unaliased.(type) {
	case *RecordAtom:
		return t.MemorySize(), t.MemoryAlignment()
	case *PrimitiveAtom:
		{
			name := t.PrimitiveName().Name()
			switch name {
			case "List":
				return Sizeof64BitPointer, Alignof64BitPointer
			case "Array":
				return Sizeof64BitPointer, Alignof64BitPointer
			case "Blob":
				return Sizeof64BitPointer, Alignof64BitPointer
			case "Bool":
				return SizeofSwampBool, AlignOfSwampBool
			case "Int":
				return SizeofSwampInt, AlignOfSwampInt
			case "Fixed":
				return SizeofSwampInt, AlignOfSwampInt
			case "ResourceName": // Resource names are translated to integers
				return SizeofSwampInt, AlignOfSwampInt
			case "TypeRef":
				return SizeofSwampInt, AlignOfSwampInt
			case "Char":
				return SizeofSwampInt, AlignOfSwampInt
			case "String":
				return Sizeof64BitPointer, Alignof64BitPointer
			case "Any":
				return Sizeof64BitPointer, Alignof64BitPointer
			}
			panic(fmt.Errorf("do not know primitive atom of '%s' %v %T", name, p, unaliased))
		}
	case *CustomTypeAtom:
		return t.MemorySize(), t.MemoryAlignment()
	case *CustomTypeVariant:
		return t.debugMemorySize, t.debugMemoryAlign
	case *FunctionAtom:
		return Sizeof64BitPointer, Alignof64BitPointer
	case *UnmanagedType:
		return Sizeof64BitPointer, Alignof64BitPointer
	case *TupleTypeAtom:
		return t.MemorySize(), t.MemoryAlignment()
	case *LocalType:
		return 0, 0
	default:
		panic(fmt.Errorf("calc: do not know memory size of %v %T", p, unaliased))
	}
}

func GetMemorySizeAndAlignment(p dtype.Type) (MemorySize, MemoryAlign) {
	memorySize, memoryAlign := GetMemorySizeAndAlignmentInternal(p)
	if memorySize == 0 || memoryAlign == 0 {
		panic(fmt.Errorf("can not be correct size and align %T %v", p, p))
	}

	return memorySize, memoryAlign
}

func calculateFieldOffsetsAndRecordMemorySizeAndAlign(fields []*RecordField) (MemorySize, MemoryAlign) {
	offset := MemoryOffset(0)
	maxMemoryAlign := MemoryAlign(0)

	for _, field := range fields {
		memorySize, memoryAlign := GetMemorySizeAndAlignment(field.fieldType)
		rest := MemoryAlign(uint32(offset) % uint32(memoryAlign))
		if rest != 0 {
			offset += MemoryOffset(memoryAlign - rest)
		}
		if memoryAlign > maxMemoryAlign {
			maxMemoryAlign = memoryAlign
		}

		field.memoryOffset = offset
		field.memorySize = memorySize

		offset += MemoryOffset(memorySize)
	}

	rest := MemoryAlign(uint32(offset) % uint32(maxMemoryAlign))
	if rest != 0 {
		offset += MemoryOffset(maxMemoryAlign - rest)
	}

	return MemorySize(offset), maxMemoryAlign
}

type ByFieldName []*RecordField

func (a ByFieldName) Len() int           { return len(a) }
func (a ByFieldName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFieldName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func NewRecordType(info *ast.Record, fields []*RecordField, genericTypes []dtype.Type) *RecordAtom {
	sortedFields := make([]*RecordField, len(fields))
	copy(sortedFields, fields)
	sort.Sort(ByFieldName(sortedFields))

	nameToField := make(map[string]*RecordField)
	for index, field := range sortedFields {
		name := field.Name()
		if nameToField[name] != nil {
			panic("we already have that struct name")
		}
		field.SetIndexBySorter(index)
		nameToField[name] = field
	}

	memorySize, memoryAlign := calculateFieldOffsetsAndRecordMemorySizeAndAlign(sortedFields)

	return &RecordAtom{
		sortedFields: sortedFields, record: info, parsedOrderFields: fields,
		nameToField: nameToField, genericTypes: genericTypes,
		memorySize: memorySize, memoryAlign: memoryAlign,
	}
}

func (s *RecordAtom) SortedFields() []*RecordField {
	return s.sortedFields
}

func (s *RecordAtom) ParseOrderedFields() []*RecordField {
	return s.parsedOrderFields
}

func (s *RecordAtom) FieldCount() int {
	return len(s.sortedFields)
}

func (s *RecordAtom) AtomName() string {
	return "recordatom"
}

func (s *RecordAtom) FindField(name string) *RecordField {
	return s.nameToField[name]
}

func (s *RecordAtom) ParameterCount() int {
	return len(s.genericTypes)
}

func (u *RecordAtom) Resolve() (dtype.Atom, error) {
	return u, nil
}

func (u *RecordAtom) Next() dtype.Type {
	return nil
}

func (u *RecordAtom) IsEqual(other_ dtype.Atom) error {
	if IsAtomAny(other_) {
		return nil
	}

	other, wasFunctionAtom := other_.(*RecordAtom)
	if !wasFunctionAtom {
		return fmt.Errorf("wasn't a record even %v", other)
	}
	otherFields := other.sortedFields
	if len(u.sortedFields) != len(otherFields) {
		return fmt.Errorf("wrong number of fields\n  %v\nvs\n   %v", u.HumanReadable(), other.HumanReadable())
	}
	for index, field := range u.sortedFields {
		otherField := otherFields[index]

		if otherField.Name() != field.Name() {
			return fmt.Errorf("field names differ '%v' <-> '%v'\n %v\n %v", otherField.name.Name(), field.name.Name(), otherField, field)
		}
		otherFieldType, otherFieldTypeErr := otherField.Type().Resolve()
		if otherFieldTypeErr != nil {
			return fmt.Errorf("couldn't resolve %w", otherFieldTypeErr)
		}

		fieldType, fieldTypeErr := field.fieldType.Resolve()
		if fieldTypeErr != nil {
			return fmt.Errorf("couldn't resolve %w", fieldTypeErr)
		}
		equalErr := fieldType.IsEqual(otherFieldType)
		if equalErr != nil {
			return fmt.Errorf("field type differs %w", equalErr)
		}
	}

	return nil
}
