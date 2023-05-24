/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"

	opcode_sp_type "github.com/swamp/opcodes/type"
)

type RecordAtom struct {
	nameToField       map[string]*RecordField
	parsedOrderFields []*RecordField
	sortedFields      []*RecordField `debug:"true"`
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

func (s *RecordAtom) AstRecord() *ast.Record {
	return s.record
}

func (s *RecordAtom) String() string {
	return fmt.Sprintf("[RecordType %v]", s.sortedFields)
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
	if len(s.parsedOrderFields) == 0 {
		panic(fmt.Errorf("not allowed to have zero record"))
	}
	lastType := s.parsedOrderFields[len(s.parsedOrderFields)-1].fieldType
	inclusive := token.MakeInclusiveSourceFileReference(s.parsedOrderFields[0].name.FetchPositionLength(),
		lastType.FetchPositionLength())
	return inclusive
}

func TypeChain(p dtype.Type, tabs int) {
	if p == nil {
		log.Printf("end")
		return
	}
	log.Printf("%v%v =>", strings.Repeat("  ", tabs), p)

	if tabs > 2 {
		log.Printf("break here")
	}

	if tabs > 20 {
		panic("too far")
	}

	TypeChain(p.Next(), tabs+1)
}

/*
	case *ResolvedLocalType:
		return GetMemorySizeAndAlignmentInternal(t.Next())
	case *ResolvedLocalTypeReference:
		return GetMemorySizeAndAlignmentInternal(t.Next())
	case *CustomTypeVariantReference:
		return GetMemorySizeAndAlignmentInternal(t.Next())
	case *LocalTypeNameOnlyContextReference:
		return 0, 8
	case *AliasReference:
		return GetMemorySizeAndAlignmentInternal(t.Next())

*/

func GetMemorySizeAndAlignmentInternal(atom dtype.Atom) (MemorySize, MemoryAlign) {
	if atom == nil {
		panic(fmt.Errorf("nil is not allowed"))
	}
	//log.Printf("unaliased: %T %v", unaliased, unaliased)
	switch t := atom.(type) {
	case *RecordAtom:
		return t.MemorySize(), t.MemoryAlignment()
	case *PrimitiveAtom:
		{
			name := t.PrimitiveName().Name()
			switch name {
			case "List":
				return MemorySize(opcode_sp_type.Sizeof64BitPointer), MemoryAlign(opcode_sp_type.Alignof64BitPointer)
			case "Array":
				return MemorySize(opcode_sp_type.Sizeof64BitPointer), MemoryAlign(opcode_sp_type.Alignof64BitPointer)
			case "Blob":
				return MemorySize(opcode_sp_type.Sizeof64BitPointer), MemoryAlign(opcode_sp_type.Alignof64BitPointer)
			case "Bool":
				return MemorySize(opcode_sp_type.SizeofSwampBool), MemoryAlign(opcode_sp_type.AlignOfSwampBool)
			case "Int":
				return MemorySize(opcode_sp_type.SizeofSwampInt), MemoryAlign(opcode_sp_type.AlignOfSwampInt)
			case "Fixed":
				return MemorySize(opcode_sp_type.SizeofSwampInt), MemoryAlign(opcode_sp_type.AlignOfSwampInt)
			case "ResourceName": // Resource names are translated to integers
				return MemorySize(opcode_sp_type.SizeofSwampInt), MemoryAlign(opcode_sp_type.AlignOfSwampInt)
			case "TypeRef":
				return MemorySize(opcode_sp_type.SizeofSwampInt), MemoryAlign(opcode_sp_type.AlignOfSwampInt)
			case "Char":
				return MemorySize(opcode_sp_type.SizeofSwampInt), MemoryAlign(opcode_sp_type.AlignOfSwampInt)
			case "String":
				return MemorySize(opcode_sp_type.Sizeof64BitPointer), MemoryAlign(opcode_sp_type.Alignof64BitPointer)
			case "Any":
				return MemorySize(opcode_sp_type.Sizeof64BitPointer), MemoryAlign(opcode_sp_type.Alignof64BitPointer)
			}
			panic(fmt.Errorf("do not know primitive atom of '%s' %v %T", name, atom, atom))
		}
	case *CustomTypeAtom:
		return t.MemorySize(), t.MemoryAlignment()
	case *CustomTypeVariantAtom:
		return t.MemorySize(), t.MemoryAlignment()
	case *FunctionAtom:
		return MemorySize(opcode_sp_type.Sizeof64BitPointer), MemoryAlign(opcode_sp_type.Alignof64BitPointer)
	case *UnmanagedType:
		return MemorySize(opcode_sp_type.Sizeof64BitPointer), MemoryAlign(opcode_sp_type.Alignof64BitPointer)
	case *TupleTypeAtom:
		return t.MemorySize(), t.MemoryAlignment()
	case *LocalTypeNameOnlyContext:
		log.Printf("LocalTypeNameOnlyContext: %v", t)
		return 0, 8
	default:
		panic(fmt.Errorf("calc: do not know memory size of %v %T", atom, atom))
	}
}

func GetMemorySizeAndAlignment(p dtype.Type) (MemorySize, MemoryAlign) {
	atom, resolveErr := p.Resolve()
	if resolveErr != nil {
		panic(resolveErr)
	}
	memorySize, memoryAlign := GetMemorySizeAndAlignmentInternal(atom)
	if memoryAlign == 0 {
		panic(fmt.Errorf("unsupported Type %T %v", p, p))
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

func NewRecordType(info *ast.Record, fields []*RecordField) *RecordAtom {
	sortedFields := make([]*RecordField, len(fields))
	copy(sortedFields, fields)
	sort.Sort(ByFieldName(sortedFields))

	nameToField := make(map[string]*RecordField)
	for index, field := range sortedFields {
		name := field.Name()
		if nameToField[name] != nil {
			panic("we already have that struct name")
		}
		/*
			_, wasPrimitive := field.Type().(*PrimitiveAtom)
			if wasPrimitive {
				panic("not allowed with primitive")
			}

		*/
		if !field.Type().FetchPositionLength().Verify() {
			log.Print("problem:")
			for index2, field2 := range sortedFields {
				if index == index2 {
					log.Printf("* here-> %v %v", field2.name, field2.fieldType)
				} else {
					log.Printf("%v %v", field2.name, field2.fieldType)
				}
			}
			panic(fmt.Errorf("stopping %T %v", field.Type(), field.Type()))
		}
		field.SetIndexBySorter(index)
		nameToField[name] = field
	}

	memorySize, memoryAlign := calculateFieldOffsetsAndRecordMemorySizeAndAlign(sortedFields)

	return &RecordAtom{
		sortedFields: sortedFields, record: info, parsedOrderFields: fields,
		nameToField: nameToField,
		memorySize:  memorySize, memoryAlign: memoryAlign,
	}
}

func (s *RecordAtom) SortedFields() []*RecordField {
	return s.sortedFields
}

func (s *RecordAtom) NameFromSortedFields() string {
	out := ""
	for index, field := range s.sortedFields {
		unaliasedType := UnaliasWithResolveInvoker(field.Type())
		if index > 0 {
			out += ":"
		}
		out += field.FieldName().Name().Name() + "_" + unaliasedType.AtomName()
	}

	return out
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
	return 0
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
		return fmt.Errorf("wasn't a record even %T %v", other_, other_)
	}
	otherFields := other.sortedFields
	if len(u.sortedFields) != len(otherFields) {
		return fmt.Errorf("wrong number of fields\n  %v\nvs\n   %v", u.HumanReadable(), other.HumanReadable())
	}
	for index, field := range u.sortedFields {
		otherField := otherFields[index]

		if otherField.Name() != field.Name() {
			return fmt.Errorf("field names differ '%v' <-> '%v'\n %v\n %v", otherField.name.Name(), field.name.Name(),
				otherField, field)
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

func (u *RecordAtom) WasReferenced() bool {
	return false // Record atom types are not reused.
}
