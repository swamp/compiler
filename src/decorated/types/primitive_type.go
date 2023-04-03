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
	if !wasPrimitive || len(primitive.ParameterTypes()) != 1 {
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

func FindNameOnlyContextWithUnalias(checkType dtype.Type) *LocalTypeNameContext {
	unliased := Unalias(checkType)
	localTypeNameContext, wasLocalTypeNameContext := unliased.(*LocalTypeNameContext)
	if !wasLocalTypeNameContext {
		localTypeNameContextRef, wasLocalTypeNameContextRef := unliased.(*LocalTypeNameContextReference)
		if wasLocalTypeNameContextRef {
			return localTypeNameContextRef.nameContext
		}
		return nil
	}

	return localTypeNameContext
}

func IsListLike(typeToCheck dtype.Type) bool {
	unaliasType := UnaliasWithResolveInvoker(typeToCheck)

	primitiveAtom, _ := unaliasType.(*PrimitiveAtom)
	if primitiveAtom == nil {
		return false
	}

	name := primitiveAtom.PrimitiveName().Name()

	return name == "List"
}

func IsIntLike(typeToCheck dtype.Type) bool {
	unaliasType := UnaliasWithResolveInvoker(typeToCheck)

	primitiveAtom, _ := unaliasType.(*PrimitiveAtom)
	if primitiveAtom == nil {
		return false
	}

	name := primitiveAtom.AtomName()

	return name == "Int" || name == "Fixed" || name == "Char"
}

func IsListAny(checkType dtype.Type) bool {
	unliased := UnaliasWithResolveInvoker(checkType)
	listAtom, err := GetListType(unliased)
	if err != nil {
		return false
	}
	return IsAny(listAtom.ParameterTypes()[0])
}

func IsLocalType(checkType dtype.Type) bool {
	unliased := UnaliasWithResolveInvoker(checkType)
	_, wasLocalType := unliased.(*LocalTypeDefinitionReference)
	return wasLocalType
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

func TryTypeIdRef(checkType dtype.Type) (*PrimitiveAtom, bool) {
	unliased := UnaliasWithResolveInvoker(checkType)
	primitive, wasPrimitive := unliased.(*PrimitiveAtom)
	if !wasPrimitive {
		return nil, false
	}

	wasTypeRef := primitive.AtomName() == "TypeRef"
	if !wasTypeRef {
		return nil, false
	}

	return primitive, true
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
	name           *ast.TypeIdentifier
	parameterTypes []dtype.Type
	references     []*PrimitiveTypeReference
	inclusive      token.SourceFileReference
}

func NewPrimitiveType(name *ast.TypeIdentifier, parameterTypes []dtype.Type) *PrimitiveAtom {
	for _, parameterType := range parameterTypes {
		if parameterType == nil {
			panic("not allowed to be nil parameterType")
		}
	}
	inclusive := name.FetchPositionLength()
	if len(parameterTypes) > 0 {
		inclusive = token.MakeInclusiveSourceFileReference(name.FetchPositionLength(), parameterTypes[len(parameterTypes)-1].FetchPositionLength())
	}
	return &PrimitiveAtom{name: name, parameterTypes: parameterTypes, inclusive: inclusive}
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

	for index, genericType := range u.parameterTypes {
		otherGenericType := other.parameterTypes[index]
		if err := CompatibleTypes(genericType, otherGenericType); err != nil {
			return fmt.Errorf("not same generic type %v and %v %v", genericType.HumanReadable(), otherGenericType.HumanReadable(), err)
		}
	}

	return nil
}

func (u *PrimitiveAtom) FetchPositionLength() token.SourceFileReference {
	return u.inclusive
}

func (u *PrimitiveAtom) PrimitiveName() *ast.TypeIdentifier {
	return u.name
}

func (u *PrimitiveAtom) String() string {
	return fmt.Sprintf("[Primitive %v%v]", u.name.Name(), TypeParametersSuffix(u.parameterTypes))
}

func (u *PrimitiveAtom) HumanReadable() string {
	return fmt.Sprintf("%v%v", u.name.Name(), TypesToHumanReadableWithinBrackets(u.parameterTypes))
}

func (u *PrimitiveAtom) AtomName() string {
	return u.name.Name()
}

func (u *PrimitiveAtom) ParameterTypes() []dtype.Type {
	return u.parameterTypes
}

func (u *PrimitiveAtom) Next() dtype.Type {
	return nil
}

func (u *PrimitiveAtom) ParameterCount() int {
	return len(u.parameterTypes)
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

func (u *PrimitiveAtom) WasReferenced() bool {
	return len(u.references) > 0
}
