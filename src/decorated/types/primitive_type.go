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

type PrimitiveAtom struct {
	name         *ast.TypeIdentifier
	genericTypes []dtype.Type
}

func NewPrimitiveType(name *ast.TypeIdentifier, genericTypes []dtype.Type) *PrimitiveAtom {
	for _, generic := range genericTypes {
		if generic == nil {
			panic("not allowed to be nil generic")
		}
	}
	return &PrimitiveAtom{name: name, genericTypes: genericTypes}
}

func (u *PrimitiveAtom) IsEqual(other_ dtype.Atom) error {
	_, isAny := other_.(*Any)
	if isAny {
		return nil
	}
	other, wasPrimitive := other_.(*PrimitiveAtom)
	if !wasPrimitive {
		return fmt.Errorf("wasn't same primitive %v", other)
	}

	if other.name.Name() != u.name.Name() {
		return fmt.Errorf("not same primitive '%v' vs '%v'", u.name, other.name)
	}

	if other.ParameterCount() != u.ParameterCount() {
		return fmt.Errorf("different number of parameters")
	}

	for index, genericType := range u.genericTypes {
		otherGenericType := other.genericTypes[index]
		if err := CompatibleTypes(genericType, otherGenericType); err != nil {
			return fmt.Errorf("not same generic type %v and %v %v", genericType.HumanReadable(), otherGenericType.HumanReadable(), err)
		}
	}

	return nil
}

func (u *PrimitiveAtom) PrimitiveName() *ast.TypeIdentifier {
	return u.name
}

func (u *PrimitiveAtom) String() string {
	return fmt.Sprintf("[primitive %v%v]", u.name.Name(), TypeParametersSuffix(u.genericTypes))
}

func (u *PrimitiveAtom) ShortString() string {
	return fmt.Sprintf("[primitive %v%v]", u.name.Name(), TypeParametersShortSuffix(u.genericTypes))
}

func (u *PrimitiveAtom) HumanReadable() string {
	return fmt.Sprintf("%v%v", u.name.Name(), TypeParametersHumanReadableSuffix(u.genericTypes))
}

func (u *PrimitiveAtom) DecoratedName() string {
	return fmt.Sprintf("%v%v", u.name.Name(), TypeParametersShortSuffix(u.genericTypes))
}

func (u *PrimitiveAtom) AtomName() string {
	return u.DecoratedName()
}

func (u *PrimitiveAtom) ShortName() string {
	return u.DecoratedName()
}

func (u *PrimitiveAtom) GenericTypes() []dtype.Type {
	return u.genericTypes
}

func (u *PrimitiveAtom) Apply(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("can not apply primitive %v", u.name)
}

func (u *PrimitiveAtom) Generate(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("can not generate primitive %v", u.name)
}

func (u *PrimitiveAtom) Next() dtype.Type {
	return nil
}

func (u *PrimitiveAtom) ParameterCount() int {
	return len(u.genericTypes)
}

func (u *PrimitiveAtom) Resolve() (dtype.Atom, error) {
	return u, nil
}
