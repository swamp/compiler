/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type Record struct {
	typeParameters []*TypeParameter
	fields         []*RecordTypeField
	startParen     token.ParenToken
	endParen       token.ParenToken
	comment        *MultilineComment

	inclusive token.SourceFileReference
}

func NewRecordType(startParen token.ParenToken, endParen token.ParenToken, fields []*RecordTypeField, typeParameters []*TypeParameter, comment *MultilineComment) *Record {
	inclusive := token.MakeInclusiveSourceFileReference(startParen.SourceFileReference, endParen.SourceFileReference)
	return &Record{fields: fields, typeParameters: typeParameters, inclusive: inclusive, comment: comment}
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

func (i *Record) Fields() []*RecordTypeField {
	return i.fields
}

func (i *Record) FindField(name string) *RecordTypeField {
	for _, f := range i.fields {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

func (i *Record) Comment() *MultilineComment {
	return i.comment
}

func (i *Record) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}
