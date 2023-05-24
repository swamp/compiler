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
	parameterType dtype.Type `debug:"true"`
}

func NewCustomTypeVariantField(index uint, fieldType dtype.Type) *CustomTypeVariantField {
	return &CustomTypeVariantField{index: index, parameterType: fieldType}
}

func (c *CustomTypeVariantField) MemoryOffset() MemoryOffset {
	return c.memoryOffset
}

func (c *CustomTypeVariantField) MemorySize() MemorySize {
	return c.memorySize
}

func (c *CustomTypeVariantField) Type() dtype.Type {
	return c.parameterType
}

func (c *CustomTypeVariantField) String() string {
	return c.parameterType.String()
}

func CustomTypeGVariantFieldsToStringSuffix(types []*CustomTypeVariantField) string {
	if len(types) == 0 {
		return ""
	}
	s := " ["
	for index, t := range types {
		if index > 0 {
			s += ","
		}
		s += t.String()
	}
	s += "]"

	return s
}

type CustomTypeVariantAtom struct {
	index                int
	astCustomTypeVariant *ast.CustomTypeVariant    `debug:"true"`
	parameterFields      []*CustomTypeVariantField `debug:"true"`
	parent               dtype.Type
	inCustomType         *CustomTypeAtom
	references           []*CustomTypeVariantReference
	debugMemorySize      MemorySize
	debugMemoryAlign     MemoryAlign
}

func NewCustomTypeVariant(index int, inCustomType *CustomTypeAtom, astCustomTypeVariant *ast.CustomTypeVariant,
	parameterTypes []dtype.Type) *CustomTypeVariantAtom {
	if inCustomType == nil {
		//		panic("must have valid in custom type")
	}
	for _, paramType := range parameterTypes {
		if paramType == nil {
			panic("paramtype is nil")
		}
	}

	var fields []*CustomTypeVariantField

	pos := MemoryOffset(1) // Leave room for the custom type localTypeNameReference
	var biggestMemoryAlign MemoryAlign
	biggestMemoryAlign = 1
	for index, paramType := range parameterTypes {
		if paramType == nil {
			panic("paramtype is nil")
		}

		_, wasLocalType := paramType.(*LocalTypeNameReference)
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

		if memoryAlign > biggestMemoryAlign {
			biggestMemoryAlign = memoryAlign
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
	if biggestMemoryAlign > 0 {
		rest := pos % MemoryOffset(biggestMemoryAlign)
		if rest != 0 {
			pos += MemoryOffset(uint(biggestMemoryAlign) - uint(rest))
		}
	}

	return &CustomTypeVariantAtom{
		index: index, astCustomTypeVariant: astCustomTypeVariant, inCustomType: inCustomType,
		parameterFields: fields, debugMemorySize: MemorySize(pos), debugMemoryAlign: biggestMemoryAlign,
	}
}

func (s *CustomTypeVariantAtom) AstCustomTypeVariant() *ast.CustomTypeVariant {
	return s.astCustomTypeVariant
}

func (s *CustomTypeVariantAtom) FetchPositionLength() token.SourceFileReference {
	return s.astCustomTypeVariant.FetchPositionLength()
}

func (s *CustomTypeVariantAtom) Fields() []*CustomTypeVariantField {
	return s.parameterFields
}

func (s *CustomTypeVariantAtom) InCustomType() *CustomTypeAtom {
	return s.inCustomType
}

func (s *CustomTypeVariantAtom) ParameterTypes() []dtype.Type {
	var types []dtype.Type

	for _, x := range s.parameterFields {
		types = append(types, x.parameterType)
	}

	return types
}

func (s *CustomTypeVariantAtom) ParentType() dtype.Type {
	if s.parent == nil {
		panic("can not fetch nil parent type")
	}
	return s.parent
}

func (s *CustomTypeVariantAtom) Index() int {
	return s.index
}

func (s *CustomTypeVariantAtom) Resolve() (dtype.Atom, error) {
	return s, nil
}

func (s *CustomTypeVariantAtom) AtomName() string {
	return "CustomTypeVariant"
}

func (s *CustomTypeVariantAtom) IsEqual(other_ dtype.Atom) error {
	other, wasCustomTypeVariant := other_.(*CustomTypeVariantAtom)
	if !wasCustomTypeVariant {
		return fmt.Errorf("was not even a custom type variant %T %v", other_, other_)
	}

	// If our union is the same, then it is ok

	if s.index != other.index {
		return fmt.Errorf("different index %d vs %d", s.index, other.index)
	}

	if s.astCustomTypeVariant.Name() != other.astCustomTypeVariant.Name() {
		return fmt.Errorf("custom type must have same name %s vs %s", s.astCustomTypeVariant.Name(),
			other.astCustomTypeVariant.Name())
	}

	if s.ParameterCount() != other.ParameterCount() {
		return fmt.Errorf("wrong parameter count %d vs %d", s.ParameterCount(), other.ParameterCount())
	}

	otherParams := other.parameterFields
	for index, parameter := range s.parameterFields {
		equalErr := CompatibleTypes(parameter.parameterType, otherParams[index].parameterType)
		if equalErr != nil {
			return FunctionAtomMismatch{parameter.parameterType, otherParams[index].parameterType}
		}
	}

	return nil
}

func (s *CustomTypeVariantAtom) Next() dtype.Type {
	return s.parent
}

func (s *CustomTypeVariantAtom) Name() *ast.TypeIdentifier {
	return s.astCustomTypeVariant.TypeIdentifier()
}

func (s *CustomTypeVariantAtom) ParameterCount() int {
	return len(s.parameterFields)
}

func (s *CustomTypeVariantAtom) String() string {
	return fmt.Sprintf("[Variant %v%v]", s.astCustomTypeVariant.TypeIdentifier(),
		CustomTypeGVariantFieldsToStringSuffix(s.parameterFields))
}

func (s *CustomTypeVariantAtom) MemorySize() MemorySize {
	return s.inCustomType.MemorySize()
}

func (s *CustomTypeVariantAtom) MemoryAlignment() MemoryAlign {
	return s.inCustomType.MemoryAlignment()
}

func (s *CustomTypeVariantAtom) HumanReadable() string {
	str := fmt.Sprintf("%v", s.astCustomTypeVariant.TypeIdentifier().Name())
	for _, parameterType := range s.parameterFields {
		str += " "
		str += parameterType.parameterType.HumanReadable()
	}

	return str
}

func (s *CustomTypeVariantAtom) DecoratedName() string {
	return s.Name().Name()
}

func (s *CustomTypeVariantAtom) AddReferee(reference *CustomTypeVariantReference) {
	s.references = append(s.references, reference)
}

func (s *CustomTypeVariantAtom) References() []*CustomTypeVariantReference {
	return s.references
}

func (s *CustomTypeVariantAtom) WasReferenced() bool {
	return len(s.references) > 0
}
