/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"github.com/swamp/compiler/src/decorated/dtype"
)

func TypeIsTemplateHasLocalTypes(p dtype.Type) bool {
	switch p.(type) {
	case *LocalTypeNameOnlyContext:
		return true
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

func ResolveToAtom(t dtype.Type) dtype.Atom {
	resolved, err := t.Resolve()
	if err != nil {
		panic(err)
	}

	return resolved
}

func ResolveToFunctionAtom(t dtype.Type) *FunctionAtom {
	atom := ResolveToAtom(t)
	functionAtom, _ := atom.(*FunctionAtom)
	return functionAtom
}

func NonResolvedFunctionParameters(t dtype.Type) []dtype.Type {
	switch info := t.(type) {
	case *FunctionAtom:
		return info.parameterTypes
	case *FunctionTypeReference:
		return info.referencedType.parameterTypes
	case *LocalTypeNameOnlyContext:
		return NonResolvedFunctionParameters(info.Next())
	}

	panic("can not convert to function parameter types")
}
