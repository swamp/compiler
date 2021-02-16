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

type CustomTypeVariant struct {
	index          int
	name           *ast.TypeIdentifier
	parameterTypes []dtype.Type
	parent         dtype.Type
	inCustomType   *CustomTypeAtom
}

func NewCustomTypeVariant(index int, name *ast.TypeIdentifier, parameterTypes []dtype.Type) *CustomTypeVariant {
	for _, paramType := range parameterTypes {
		if paramType == nil {
			panic("paramtype is nil")
		}
	}
	return &CustomTypeVariant{index: index, name: name, parameterTypes: parameterTypes}
}

func (s *CustomTypeVariant) AttachToCustomType(c *CustomTypeAtom) {
	if s.parent != nil {
		panic("already attached")
	}

	s.parent = c
	s.inCustomType = c
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

func (s *CustomTypeVariant) Generate(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("can not generate")
}

func (s *CustomTypeVariant) Resolve() (dtype.Atom, error) {
	return nil, nil
}

func (s *CustomTypeVariant) Next() dtype.Type {
	return s.parent
}

func (s *CustomTypeVariant) Name() *ast.TypeIdentifier {
	return s.name
}

func (s *CustomTypeVariant) ParameterTypes() []dtype.Type {
	return s.parameterTypes
}

func (s *CustomTypeVariant) ParameterCount() int {
	return len(s.parameterTypes)
}

func (s *CustomTypeVariant) String() string {
	return fmt.Sprintf("[variant %v%v]", s.name, TypesToStringSuffix(s.parameterTypes))
}

func (s *CustomTypeVariant) ShortString() string {
	return fmt.Sprintf("[variant %v%v]", s.name, TypesToShortStringSuffix(s.parameterTypes))
}

func (s *CustomTypeVariant) HumanReadable() string {
	str := fmt.Sprintf("%v", s.name)
	for _, parameterType := range s.parameterTypes {
		str += " "
		str += parameterType.HumanReadable()
	}

	return str
}

func (s *CustomTypeVariant) DecoratedName() string {
	return s.Name().Name()
}

func (s *CustomTypeVariant) ShortName() string {
	return s.DecoratedName()
}
