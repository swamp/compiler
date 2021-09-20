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

type CustomTypeVariantField struct {
	index         uint
	memoryOffset  MemoryOffset
	memorySize    MemorySize
	parameterType dtype.Type
}

func (c *CustomTypeVariantField) MemoryOffset() MemoryOffset {
	return c.memoryOffset
}

func (c *CustomTypeVariantField) MemorySize() MemorySize {
	return c.memorySize
}

type CustomTypeVariant struct {
	index                int
	astCustomTypeVariant *ast.CustomTypeVariant
	parameterTypes       []dtype.Type
	parameterFields      []*CustomTypeVariantField
	parent               dtype.Type
	inCustomType         *CustomTypeAtom
	references           []*CustomTypeVariantReference
	debugMemorySize      MemorySize
	debugMemoryAlign     MemoryAlign
}

func NewCustomTypeVariant(index int, astCustomTypeVariant *ast.CustomTypeVariant, parameterTypes []dtype.Type) *CustomTypeVariant {
	for _, paramType := range parameterTypes {
		if paramType == nil {
			panic("paramtype is nil")
		}
	}

	var fields []*CustomTypeVariantField

	pos := MemoryOffset(1) // Leave room for the custom type identifier

	for index, paramType := range parameterTypes {
		if paramType == nil {
			panic("paramtype is nil")
		}

		_, wasLocalType := paramType.(*LocalType)
		var memorySize MemorySize
		var memoryAlign MemoryAlign

		if wasLocalType {
			memorySize = 0
			memoryAlign = 0
		} else {
			memorySize, memoryAlign = GetMemorySizeAndAlignment(paramType)
			rest := pos % MemoryOffset(memoryAlign)
			if rest != 0 {
				pos += MemoryOffset(uint(memoryAlign) - uint(rest))
			}
		}

		field := &CustomTypeVariantField{
			index:         uint(index),
			memoryOffset:  pos,
			memorySize:    memorySize,
			parameterType: paramType,
		}

		pos += MemoryOffset(memorySize)

		fields = append(fields, field)
	}
	return &CustomTypeVariant{index: index, astCustomTypeVariant: astCustomTypeVariant, parameterTypes: parameterTypes, parameterFields: fields}
}

func (s *CustomTypeVariant) AttachToCustomType(c *CustomTypeAtom) {
	if s.parent != nil {
		panic("already attached")
	}

	s.parent = c
	s.inCustomType = c
}

func (s *CustomTypeVariant) AstCustomTypeVariant() *ast.CustomTypeVariant {
	return s.astCustomTypeVariant
}

func (s *CustomTypeVariant) FetchPositionLength() token.SourceFileReference {
	return s.astCustomTypeVariant.FetchPositionLength()
}

func (s *CustomTypeVariant) Fields() []*CustomTypeVariantField {
	return s.parameterFields
}

func (s *CustomTypeVariant) InCustomType() *CustomTypeAtom {
	return s.inCustomType
}

func (s *CustomTypeVariant) ParentType() dtype.Type {
	if s.parent == nil {
		panic("can not fetch nil parent type")
	}
	return s.parent
}

func (s *CustomTypeVariant) Index() int {
	return s.index
}

func (s *CustomTypeVariant) Resolve() (dtype.Atom, error) {
	return nil, nil
}

func (s *CustomTypeVariant) Next() dtype.Type {
	return s.parent
}

func (s *CustomTypeVariant) Name() *ast.TypeIdentifier {
	return s.astCustomTypeVariant.TypeIdentifier()
}

func (s *CustomTypeVariant) ParameterTypes() []dtype.Type {
	return s.parameterTypes
}

func (s *CustomTypeVariant) ParameterCount() int {
	return len(s.parameterTypes)
}

func (s *CustomTypeVariant) String() string {
	return fmt.Sprintf("[variant %v%v]", s.astCustomTypeVariant.TypeIdentifier(), TypesToStringSuffix(s.parameterTypes))
}

func (s *CustomTypeVariant) HumanReadable() string {
	str := fmt.Sprintf("%v", s.astCustomTypeVariant.TypeIdentifier())
	for _, parameterType := range s.parameterTypes {
		str += " "
		str += parameterType.HumanReadable()
	}

	return str
}

func (s *CustomTypeVariant) DecoratedName() string {
	return s.Name().Name()
}

func (s *CustomTypeVariant) AddReferee(reference *CustomTypeVariantReference) {
	s.references = append(s.references, reference)
}

func (s *CustomTypeVariant) References() []*CustomTypeVariantReference {
	return s.references
}
