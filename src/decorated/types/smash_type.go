/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"reflect"

	"github.com/swamp/compiler/src/decorated/dtype"
)

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
	case *FunctionTypeReference:
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

func UnaliasWithResolveInvoker(t dtype.Type) dtype.Type {
	/*
		alias, wasAlias := t.(*Alias)
		if wasAlias {
			return Unalias(alias.referencedType)
		}

		typeRef, wasTypeRef := t.(*TypeReference)
		if wasTypeRef {
			return Unalias(typeRef.Next())
		}

		scopedTypeRef, wasScopedTypeRef := t.(*TypeReferenceScoped)
		if wasScopedTypeRef {
			return Unalias(scopedTypeRef.Next())
		}

	*/
	unaliased := Unalias(t)

	call, wasCall := unaliased.(*InvokerType)
	if wasCall {
		resolved, resolveErr := CallType(call.typeToInvoke, call.params)
		if resolveErr != nil {
			panic(resolveErr)
		}
		return Unalias(resolved)
	}

	return unaliased
}

func fillContextFromPrimitive(context *TypeParameterContextOther, original *PrimitiveAtom, other *PrimitiveAtom) (*PrimitiveAtom, error) {
	var converted []dtype.Type

	wasConverted := false
	for index, funcParam := range original.GenericTypes() {
		otherType := other.GenericTypes()[index]
		convertedType, convertErr := smashTypes(context, funcParam, otherType)
		if convertErr != nil {
			return nil, convertErr
		}
		if convertedType != funcParam {
			wasConverted = true
		}

		if convertedType == nil {
			panic("strange")
		}
		converted = append(converted, convertedType)
	}

	if !wasConverted {
		return original, nil
	}

	return NewPrimitiveType(other.name, converted), nil
}

func fillContextFromRecordType(context *TypeParameterContextOther, original *RecordAtom, other *RecordAtom) (*RecordAtom, error) {
	var converted []dtype.Type

	wasConverted := false
	for index, funcParam := range original.GenericTypes() {
		otherType := other.GenericTypes()[index]
		convertedType, convertErr := smashTypes(context, funcParam, otherType)
		if convertErr != nil {
			return nil, convertErr
		}
		if convertedType != funcParam {
			wasConverted = true
		}

		if convertedType == nil {
			panic("strange")
		}
		converted = append(converted, convertedType)
	}

	if !wasConverted {
		return original, nil
	}

	return NewRecordType(original.AstRecord(), other.SortedFields(), converted), nil
}

func fillContextFromCustomType(context *TypeParameterContextOther, original *CustomTypeAtom, other *CustomTypeAtom) (*CustomTypeAtom, error) {
	if len(original.Variants()) != len(other.Variants()) {
		return nil, fmt.Errorf("not the same number of variants")
	}

	wasConverted := false
	var convertedVariants []*CustomTypeVariant
	for index, originalVariant := range original.Variants() {
		otherVariant := other.Variants()[index]
		if otherVariant.index != originalVariant.index {
			return nil, fmt.Errorf("index error")
		}
		if len(originalVariant.ParameterTypes()) != len(otherVariant.ParameterTypes()) {
			return nil, fmt.Errorf("wrong number of parameter types in variant %v", otherVariant)
		}
		var convertedParams []dtype.Type
		for paramIndex, originalParam := range originalVariant.ParameterTypes() {
			otherParam := otherVariant.ParameterTypes()[paramIndex]
			convertedParam, smashErr := smashTypes(context, originalParam, otherParam)
			if smashErr != nil {
				return nil, smashErr
			}
			if convertedParam != originalParam {
				wasConverted = true
			}
			convertedParams = append(convertedParams, convertedParam)
		}

		convertedVariant := NewCustomTypeVariant(originalVariant.index, originalVariant.astCustomTypeVariant, convertedParams)
		convertedVariants = append(convertedVariants, convertedVariant)
	}

	if !wasConverted {
		return original, nil
	}

	return NewCustomType(original.astCustomType, ArtifactFullyQualifiedTypeName{ModuleName{path: nil}}, nil, convertedVariants), nil
}

func fillContextFromFunctions(context *TypeParameterContextOther, original *FunctionAtom, other *FunctionAtom) (*FunctionAtom, error) {
	var converted []dtype.Type

	if hasAnyMatching, startIndex := HasAnyMatchingTypes(original.parameterTypes); hasAnyMatching {
		originalInitialCount := startIndex
		originalEndCount := len(original.parameterTypes) - startIndex - 2

		originalFirst := append([]dtype.Type{}, original.parameterTypes[0:startIndex]...)

		if len(other.parameterTypes) >= len(original.parameterTypes) {
			originalEndCount++
		}

		otherMiddle := other.parameterTypes[originalInitialCount : len(other.parameterTypes)-originalEndCount]
		if len(otherMiddle) < 1 {
			return nil, fmt.Errorf("currently, you must have at least one wildcard parameter")
		}

		originalEnd := original.parameterTypes[startIndex+1 : len(original.parameterTypes)]

		allConverted := append(originalFirst, otherMiddle...)
		allConverted = append(allConverted, originalEnd...)

		created := NewFunctionAtom(original.astFunctionType, allConverted)

		return fillContextFromFunctions(context, created, other)
	} else {
		if len(original.parameterTypes) < len(other.parameterTypes) {
			return nil, fmt.Errorf("too few parameter types")
		}
	}

	for index, otherParam := range other.parameterTypes {
		originalParam := original.parameterTypes[index]
		convertedType, convertErr := smashTypes(context, originalParam, otherParam)
		if convertErr != nil {
			return nil, convertErr
		}
		if convertedType == nil {
			panic("converted was nil")
		}

		converted = append(converted, convertedType)
	}

	for index := len(other.parameterTypes); index < len(original.parameterTypes); index++ {
		originalParam := original.parameterTypes[index]
		convertedType, replaceErr := ReplaceTypeFromContext(originalParam, context)
		if replaceErr != nil {
			return nil, replaceErr
		}
		if convertedType == nil {
			panic(fmt.Sprintf("conversion is not working %v %T", originalParam, originalParam))
		}
		converted = append(converted, convertedType)
	}

	return NewFunctionAtom(original.astFunctionType, converted), nil
}

func fillContextFromTuples(context *TypeParameterContextOther, original *TupleTypeAtom, other *TupleTypeAtom) (*TupleTypeAtom, error) {
	var converted []*TupleTypeField

	if len(original.parameterTypes) < len(other.parameterTypes) {
		return nil, fmt.Errorf("too few parameter types")
	}

	for index, otherParam := range other.parameterTypes {
		originalParam := original.parameterTypes[index]
		convertedType, convertErr := smashTypes(context, originalParam, otherParam)
		if convertErr != nil {
			return nil, convertErr
		}
		if convertedType == nil {
			panic("converted was nil")
		}

		field := NewTupleTypeField(index, convertedType)
		converted = append(converted, field)
	}

	for index := len(other.parameterTypes); index < len(original.parameterTypes); index++ {
		originalParam := original.parameterTypes[index]
		convertedType, replaceErr := ReplaceTypeFromContext(originalParam, context)
		if replaceErr != nil {
			return nil, replaceErr
		}
		if convertedType == nil {
			panic(fmt.Sprintf("conversion is not working %v %T", originalParam, originalParam))
		}
		field := NewTupleTypeField(index, convertedType)
		converted = append(converted, field)
	}

	return NewTupleTypeAtom(original.astTupleType, converted), nil
}

func smashTypes(context *TypeParameterContextOther, originalUnchanged dtype.Type, otherUnchanged dtype.Type) (dtype.Type, error) {
	if reflect.ValueOf(originalUnchanged).IsNil() {
		panic("original was nil")
	}
	if reflect.ValueOf(otherUnchanged).IsNil() {
		panic("otherUnchanged was nil")
	}

	original := UnaliasWithResolveInvoker(originalUnchanged)
	other := UnaliasWithResolveInvoker(otherUnchanged)

	if reflect.ValueOf(other).IsNil() {
		panic("other was nil")
	}

	localType, wasLocalType := original.(*LocalType)
	if wasLocalType {
		_, wasLocalType := other.(*LocalType)
		if wasLocalType {
			return nil, fmt.Errorf("not great")
		}
		return context.SpecialSet(localType.identifier.Name(), otherUnchanged)
	}

	otherIsAny := IsAny(originalUnchanged)
	originalIsAny := IsAny(original)

	if originalIsAny {
		original = otherUnchanged
	}

	if !otherIsAny && !originalIsAny {
		sameType := reflect.TypeOf(original) == reflect.TypeOf(other)
		if !sameType {
			return nil, fmt.Errorf("not even same reflect type %T vs %T\n%v\n vs\n%v", original, other, original.HumanReadable(), other.HumanReadable())
		}
	}

	switch t := original.(type) {
	case *AliasReference:
		{
			return originalUnchanged, nil
		}
	case *CustomTypeReference:
		{
			return originalUnchanged, nil
		}
	case *PrimitiveTypeReference:
		{
			return originalUnchanged, nil
		}
	case *FunctionTypeReference:
		{
			return originalUnchanged, nil
		}

	case *FunctionAtom:
		{
			otherFunc := other.(*FunctionAtom)
			return fillContextFromFunctions(context, t, otherFunc)
		}
	case *TupleTypeAtom:
		{
			otherTuple := other.(*TupleTypeAtom)
			return fillContextFromTuples(context, t, otherTuple)
		}
	case *PrimitiveAtom:
		{
			if otherIsAny {
				return originalUnchanged, nil
			}
			otherPrimitive := other.(*PrimitiveAtom)
			if otherPrimitive.PrimitiveName().Name() != t.PrimitiveName().Name() {
				return nil, fmt.Errorf("not same primitive type. %v vs %v", t.PrimitiveName(), otherPrimitive.PrimitiveName())
			}
			return fillContextFromPrimitive(context, t, otherPrimitive)
		}
	case *CustomTypeAtom:
		{
			otherCustomType := other.(*CustomTypeAtom)
			return fillContextFromCustomType(context, t, otherCustomType)
		}
	case *CustomTypeVariant:
		{
			fmt.Printf("\n\nFOUND CUSTOM TYPE ATOM variant:%v\n\n", t)
		}
	case *RecordAtom:
		{
			otherRecordType, _ := other.(*RecordAtom)
			if otherRecordType == nil {
				return nil, fmt.Errorf("how can this happen %T and %T", t, other)
			}
			return fillContextFromRecordType(context, t, otherRecordType)
		}
	case *InvokerType:
		var converted []dtype.Type
		for _, param := range t.params {
			localType, wasLocal := param.(*LocalType)
			if wasLocal {
				foundType := context.LookupTypeFromName(localType.identifier.Name())
				if foundType == nil {
					return nil, fmt.Errorf("couldn't find lookup from name %v", localType.identifier)
				}
				converted = append(converted, foundType)
			} else {
				converted = append(converted, param)
			}
		}
		resolved, resolveErr := CallType(t.typeToInvoke, converted)
		if resolveErr != nil {
			return nil, resolveErr
		}

		return smashTypes(context, resolved, other)
	case *UnmanagedType:
		{
			if otherIsAny {
				return originalUnchanged, nil
			}
			otherUnmanaged := other.(*UnmanagedType)
			if err := t.IsEqualUnmanaged(otherUnmanaged); err != nil {
				return nil, fmt.Errorf("not same unmanaged type. %v", err)
			}
			return originalUnchanged, nil
		}
	default:
		return nil, fmt.Errorf("smash type: Not handled:%T %v\n", t, t)
	}

	return nil, fmt.Errorf("unhandled type: %T %v\n", original, original)
}

func SmashFunctions(original *FunctionAtom, otherFunc *FunctionAtom) (*FunctionAtom, error) {
	context := NewTypeParameterContextOther()

	result, resultErr := fillContextFromFunctions(context, original, otherFunc)
	if resultErr != nil {
		return nil, resultErr
	}

	return result, nil
}
