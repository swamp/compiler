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

func Unalias(t dtype.Type) dtype.Type {
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

	return t
}

func UnaliasWithResolveInvoker(t dtype.Type) dtype.Type {
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

	fnTypeRef, wasFnTypeRef := t.(*FunctionTypeReference)
	if wasFnTypeRef {
		return Unalias(fnTypeRef.Next())
	}

	call, wasCall := t.(*InvokerType)
	if wasCall {
		resolved, resolveErr := CallType(call.typeToInvoke, call.params)
		if resolveErr != nil {
			panic(resolveErr)
		}
		return Unalias(resolved)
	}

	return t
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

		convertedVariant := NewCustomTypeVariant(originalVariant.index, originalVariant.name, convertedParams)
		convertedVariants = append(convertedVariants, convertedVariant)
	}

	if !wasConverted {
		return original, nil
	}

	return NewCustomType(original.astCustomType, ArtifactFullyQualifiedTypeName{ModuleName{path: nil}}, nil, convertedVariants), nil
}

func fillContextFromFunctions(context *TypeParameterContextOther, original *FunctionAtom, other *FunctionAtom) (*FunctionAtom, error) {
	var converted []dtype.Type

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

func smashTypes(context *TypeParameterContextOther, original dtype.Type, otherUnchanged dtype.Type) (dtype.Type, error) {
	if reflect.ValueOf(original).IsNil() {
		panic("original was nil")
	}
	if reflect.ValueOf(otherUnchanged).IsNil() {
		panic("otherUnchanged was nil")
	}

	original = UnaliasWithResolveInvoker(original)
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

	_, otherIsAny := other.(*Any)

	if !otherIsAny {
		sameType := reflect.TypeOf(original) == reflect.TypeOf(other)
		if !sameType {
			fmt.Printf("\n\nNOTE SAME TYPE:%T %T\n\n", original, otherUnchanged)
			return nil, fmt.Errorf("not even same reflect type %T vs %T\n%v\n vs\n%v", original, other, original, other)
		}
	}

	switch t := original.(type) {
	case *FunctionAtom:
		{
			otherFunc := other.(*FunctionAtom)
			return fillContextFromFunctions(context, t, otherFunc)
		}
	case *PrimitiveAtom:
		{
			if otherIsAny {
				return original, nil
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
			otherRecordType := other.(*RecordAtom)
			return fillContextFromRecordType(context, t, otherRecordType)
		}
	case *InvokerType:
		resolved, resolveErr := CallType(t.typeToInvoke, t.params)
		if resolveErr != nil {
			return nil, resolveErr
		}

		fmt.Printf("\n\nafter invoker %v %v\n %T and %T\n", resolved, other, resolved, other)
		return smashTypes(context, resolved, other)
	default:
		return nil, fmt.Errorf("Not handled:%T %v\n", t, t)
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

func SmashType(original dtype.Type, otherUnchanged dtype.Type) (dtype.Type, error) {
	context := NewTypeParameterContextOther()

	return smashTypes(context, original, otherUnchanged)
}
