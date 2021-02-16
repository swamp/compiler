/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import "fmt"

type Record struct {
	typeParameters []*TypeParameter
	fields         []*RecordField
}

func NewRecordType(fields []*RecordField, typeParameters []*TypeParameter) *Record {
	return &Record{fields: fields, typeParameters: typeParameters}
}

func (i *Record) TypeParameters() []*TypeParameter {
	return i.typeParameters
}

func (i *Record) Name() string {
	return "RecordType"
}

func (i *Record) String() string {
	return fmt.Sprintf("[record-type %v %v]", i.fields, i.typeParameters)
}

func (i *Record) Fields() []*RecordField {
	return i.fields
}

func (i *Record) FindField(name string) *RecordField {
	for _, f := range i.fields {
		if f.Name() == name {
			return f
		}
	}
	return nil
}
