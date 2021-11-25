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

func GetListType(p dtype.Type) (*PrimitiveAtom, error) {
	unresolvedType := UnaliasWithResolveInvoker(p)
	primitive, wasPrimitive := unresolvedType.(*PrimitiveAtom)
	if !wasPrimitive || len(primitive.GenericTypes()) != 1 {
		return nil, fmt.Errorf("wasnt a list type")
	}

	return primitive, nil
}

func IsAny(checkType dtype.Type) bool {
	unliased := UnaliasWithResolveInvoker(checkType)
	primitive, wasPrimitive := unliased.(*PrimitiveAtom)
	if !wasPrimitive {
		return false
	}

	return primitive.PrimitiveName().Name() == "Any"
}

func IsListAny(checkType dtype.Type) bool {
	unliased := UnaliasWithResolveInvoker(checkType)
	listAtom, err := GetListType(unliased)
	if err != nil {
		return false
	}
	return IsAny(listAtom.GenericTypes()[0])
}

func IsTypeIdRef(checkType dtype.Type) bool {
	unliased := UnaliasWithResolveInvoker(checkType)
	primitive, wasPrimitive := unliased.(*PrimitiveAtom)
	if !wasPrimitive {
		return false
	}

	wasTypeRef := primitive.AtomName() == "TypeRef"

	return wasTypeRef
}

func ArgumentNeedsTypeIdInsertedBefore(p dtype.Type) bool {
	unaliased := UnaliasWithResolveInvoker(p)
	return IsAny(unaliased)
}

func IsAnyOrFunctionWithAnyMatching(p dtype.Type) bool {
	if IsAny(p) {
		return true
	}

	unalias := UnaliasWithResolveInvoker(p)
	functionAtom, wasFunctionAtom := unalias.(*FunctionAtom)
	if wasFunctionAtom {
		for _, param := range functionAtom.FunctionParameterTypes() {
			_, isAnyMatching := param.(*AnyMatchingTypes)
			if isAnyMatching {
				return true
			}
		}
	}
	return false
}

func IsAtomAny(checkType dtype.Atom) bool {
	primitive, wasPrimitive := checkType.(*PrimitiveAtom)
	if !wasPrimitive {
		return false
	}

	return primitive.AtomName() == "Any"
}

type PrimitiveAtom struct {
	name         *ast.TypeIdentifier
	genericTypes []dtype.Type
	references   []*PrimitiveTypeReference
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
	other, wasPrimitive := other_.(*PrimitiveAtom)
	if !wasPrimitive {
		return fmt.Errorf("wasn't same primitive %v", other)
	}

	isAny := other.name.Name() == "Any"
	if isAny {
		return nil
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

func (u *PrimitiveAtom) FetchPositionLength() token.SourceFileReference {
	return u.name.FetchPositionLength()
}

func (u *PrimitiveAtom) PrimitiveName() *ast.TypeIdentifier {
	return u.name
}

func (u *PrimitiveAtom) String() string {
	return fmt.Sprintf("[primitive %v%v]", u.name.Name(), TypeParametersSuffix(u.genericTypes))
}

func (u *PrimitiveAtom) HumanReadable() string {
	return fmt.Sprintf("%v", u.name.Name())
}

func (u *PrimitiveAtom) AtomName() string {
	return u.name.Name()
}

func (u *PrimitiveAtom) GenericTypes() []dtype.Type {
	return u.genericTypes
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

func (u *PrimitiveAtom) AddReferee(reference *PrimitiveTypeReference) {
	u.references = append(u.references, reference)
}

func (u *PrimitiveAtom) References() []*PrimitiveTypeReference {
	return u.references
}
