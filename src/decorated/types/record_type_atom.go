/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"sort"

	"github.com/swamp/compiler/src/decorated/dtype"
)

type RecordAtom struct {
	nameToField  map[string]*RecordField
	parsedOrderFields []*RecordField
	sortedFields []*RecordField
	genericLocalTypeNames []*dtype.TypeArgumentName
}

func (s *RecordAtom) String() string {
	return fmt.Sprintf("record-type %v%v]", s.sortedFields, TypeArgumentNamesSuffix(s.genericLocalTypeNames))
}

func (s *RecordAtom) HumanReadable() string {
	str := "{"
	for index, field := range s.sortedFields {
		if index > 0 {
			str += ", "
		}
		str += field.HumanReadable()
	}
	str += "}"
	return str
}


func (s *RecordAtom) ShortString() string {
	str := "[record-type "
	for _, field := range s.sortedFields {
		str += " " + field.ShortString()
	}
	str += "]"

	return str
}

func (s *RecordAtom) DecoratedName() string {
	str := "{"
	for index, field := range s.sortedFields {
		if index > 0 {
			str += ";"
		}
		str += field.DecoratedName()
	}
	str += "}"

	return str
}

func (s *RecordAtom) ShortName() string {
	return s.DecoratedName()
}

func (s *RecordAtom) ConcretizedName() string {
	return s.DecoratedName()
}

type ByFieldName []*RecordField

func (a ByFieldName) Len() int           { return len(a) }
func (a ByFieldName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFieldName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func NewRecordType(fields []*RecordField, genericLocalTypeNames []*dtype.TypeArgumentName) *RecordAtom {
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

	return &RecordAtom{sortedFields: sortedFields, parsedOrderFields: fields, nameToField: nameToField, genericLocalTypeNames:genericLocalTypeNames}
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
	return len(s.genericLocalTypeNames)
}

func (u *RecordAtom) Apply(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("record type does not have apply")
}

func (u *RecordAtom) Generate(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("record type does not have apply")
}

func (u *RecordAtom) Resolve() (dtype.Atom, error) {
	return u, nil
}

func (u *RecordAtom) Next() dtype.Type {
	return nil
}

func (u *RecordAtom) IsEqual(other_ dtype.Atom) error {
	_, isAny := other_.(*Any)
	if isAny {
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

		if otherField.name.Name() != field.name.Name() {
			return fmt.Errorf("field names differ\n %v\n %v", otherField, field)
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
