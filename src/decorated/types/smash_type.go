/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"github.com/swamp/compiler/src/decorated/dtype"
)

func TypesIsTemplateHasLocalTypes(p []dtype.Type) bool {
	for _, x := range p {
		if TypeIsTemplateHasLocalTypes(x) {
			return true
		}
	}

	return false
}

func TypeIsTemplateHasLocalTypes(p dtype.Type) bool {
	atom := UnaliasWithResolveInvoker(p)
	switch t := atom.(type) {
	case *CustomTypeAtom:
		for _, variant := range t.Variants() {
			if TypesIsTemplateHasLocalTypes(variant.ParameterTypes()) {
				return true
			}
		}
	case *CustomTypeVariantAtom:
		if TypesIsTemplateHasLocalTypes(t.ParameterTypes()) {
			return true
		}
	case *FunctionAtom:
		if TypesIsTemplateHasLocalTypes(t.FunctionParameterTypes()) && !IsAnyOrFunctionWithAnyMatching(t) {
			return true
		}
	}

	return false
}

func UnReference(t dtype.Type) dtype.Type {
	fnTypeRef, wasFnTypeRef := t.(*FunctionTypeReference)
	if wasFnTypeRef {
		return Unalias(fnTypeRef.Next())
	}

	switch info := t.(type) {
	case *AliasReference:
		return UnReference(info.reference)
	case *PrimitiveTypeReference:
		return UnReference(info.primitiveType)
	case *CustomTypeReference:
		return UnReference(info.customType)
	case *CustomTypeVariantReference:
		return UnReference(info.customTypeVariant)
	case *FunctionTypeReference:
		return UnReference(info.referencedType)
	case *ResolvedLocalType:
		return UnReference(info.referencedType)
	}

	return t
}

func Unalias(t dtype.Type) dtype.Type {
	unref := UnReference(t)
	alias, wasAlias := unref.(*Alias)
	if wasAlias {
		return Unalias(alias.referencedType)
	}

	return unref
}

func UnaliasWithResolveInvoker(t dtype.Type) dtype.Atom {
	unaliased := Unalias(t)

	resolved, err := unaliased.Resolve()
	if err != nil {
		panic(err)
	}

	return resolved
}
