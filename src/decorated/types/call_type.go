/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/decorated/dtype"
)

type Lookup interface {
	LookupType(name string) (dtype.Type, error)
}

func ReplaceTypeFromContext(originalTarget dtype.Type, lookup Lookup) (dtype.Type, error) {
	target := Unalias(originalTarget)
	switch info := originalTarget.(type) {
	case *LocalType:
		newType, newTypeErr := lookup.LookupType(info.identifier.Name())
		if newTypeErr != nil {
			return nil, newTypeErr
		}
		if newType == nil {
			panic(fmt.Sprintf("couldn't lookup %v", info.identifier.Name()))
		}
		return newType, nil
	case *PrimitiveAtom:
		var convertedTypes []dtype.Type
		for _, someType := range info.genericTypes {
			convertedType, convertedErr := ReplaceTypeFromContext(someType, lookup)
			if convertedErr != nil {
				return nil, convertedErr
			}
			convertedTypes = append(convertedTypes, convertedType)
		}

		return NewPrimitiveType(info.name, convertedTypes), nil
	case *PrimitiveTypeReference:
		return originalTarget, nil
	case *AliasReference:
		return originalTarget, nil
	case *RecordAtom:
		return replaceRecordFromContext(info, lookup)
	case *TupleTypeAtom:
		return replaceTupleTypeFromContext(info, lookup)
	case *CustomTypeVariantAtom:
		return replaceCustomTypeVariantFromContext(info.inCustomType, info, lookup)
	case *InvokerType:
		return replaceInvokerTypeFromContext(info, lookup)
	case *CustomTypeAtom:
		// TODO: fix this
		return originalTarget, nil
	case *CustomTypeReference:
		// TODO: fix this
		return originalTarget, nil
	case *FunctionTypeReference:
		return originalTarget, nil
	default:
		log.Printf("warning: not sure what to do with %T %v. Returning same type for now", originalTarget, originalTarget)
		return nil, fmt.Errorf("not sure what to do with %T %v. Returning same type for now", target, target)
	}

	return originalTarget, nil
}

func replaceRecordFromContext(record *RecordAtom, lookup Lookup) (*RecordAtom, error) {
	var replacedFields []*RecordField

	hasLocalTypes := false
	for _, field := range record.SortedFields() {
		converted, convertedErr := ReplaceTypeFromContext(field.Type(), lookup)
		if convertedErr != nil {
			return nil, convertedErr
		}
		if converted == nil {
			panic("converted is nil")
		}

		if IsLocalType(converted) {
			hasLocalTypes = true
		}

		newField := NewRecordField(field.name, field.AstRecordTypeField(), converted)

		replacedFields = append(replacedFields, newField)
	}

	var genericTypes []dtype.Type
	if hasLocalTypes {
		genericTypes = record.genericTypes
	}

	return NewRecordType(record.AstRecord(), replacedFields, genericTypes), nil
}

func replaceInvokerTypeFromContext(invoker *InvokerType, lookup Lookup) (*InvokerType, error) {
	var convertedTypes []dtype.Type

	for _, param := range invoker.params {
		converted, convertedErr := ReplaceTypeFromContext(param, lookup)
		if convertedErr != nil {
			return nil, convertedErr
		}
		if converted == nil {
			panic(fmt.Sprintf("couldn't replace from context %v %T", param, param))
		}

		convertedTypes = append(convertedTypes, converted)
	}

	return NewInvokerType(invoker.typeToInvoke, convertedTypes)
}

func replaceCustomTypeFromContext(customType *CustomTypeAtom, lookup Lookup) (*CustomTypeAtom, error) {
	var replacedVariants []*CustomTypeVariantAtom

	var replacedGenerics []dtype.Type
	for _, genericType := range customType.Parameters() {
		localType, wasLocalType := genericType.(*LocalType)
		foundGeneric := genericType
		if wasLocalType {
			lookedUpType, err := lookup.LookupType(localType.identifier.Name())
			if err != nil {
				return nil, err
			}
			foundGeneric = lookedUpType
		}
		replacedGenerics = append(replacedGenerics, foundGeneric)
	}

	newCustomType := NewCustomTypePrepare(customType.astCustomType, customType.artifactTypeName, replacedGenerics)
	for _, field := range customType.Variants() {
		var variantParameters []dtype.Type
		for _, param := range field.parameterFields {
			converted, convertedErr := ReplaceTypeFromContext(param.parameterType, lookup)
			if convertedErr != nil {
				return nil, convertedErr
			}
			variantParameters = append(variantParameters, converted)
		}

		newField, newErr := replaceCustomTypeVariantFromContext(newCustomType, field, lookup)
		if newErr != nil {
			return nil, newErr
		}

		replacedVariants = append(replacedVariants, newField)
	}

	newCustomType.FinalizeVariants(replacedVariants)

	return newCustomType, nil
}

func replaceTupleTypeFromContext(tupleType *TupleTypeAtom, lookup Lookup) (dtype.Type, error) {
	var convertedTypes []*TupleTypeField
	for index, someType := range tupleType.parameterTypes {
		convertedType, convertedErr := ReplaceTypeFromContext(someType, lookup)
		if convertedErr != nil {
			return nil, convertedErr
		}
		field := NewTupleTypeField(index, convertedType)
		convertedTypes = append(convertedTypes, field)
	}

	return NewTupleTypeAtom(tupleType.astTupleType, convertedTypes), nil
}

func replaceCustomTypeVariantFromContext(inCustomType *CustomTypeAtom, customTypeVariant *CustomTypeVariantAtom, lookup Lookup) (*CustomTypeVariantAtom, error) {
	var convertedTypes []dtype.Type
	for _, someType := range customTypeVariant.parameterFields {
		convertedType, convertedErr := ReplaceTypeFromContext(someType.parameterType, lookup)
		if convertedErr != nil {
			return nil, convertedErr
		}
		convertedTypes = append(convertedTypes, convertedType)
	}

	return NewCustomTypeVariant(customTypeVariant.index, inCustomType, customTypeVariant.astCustomTypeVariant, convertedTypes), nil
}

func callRecordType(record *RecordAtom, arguments []dtype.Type) (dtype.Type, error) {
	genericTypes := record.GenericTypes()
	if len(arguments) == 0 {
		return nil, fmt.Errorf("no arguments for type %v", record)
	}

	if len(genericTypes) != len(arguments) {
		return nil, fmt.Errorf("call record type illegal count %v %v %v", record, genericTypes, arguments)
	}

	var convertedTypes []dtype.Type

	foundLocal := false
	for index, foundType := range genericTypes {
		_, wasLocal := foundType.(*LocalType)
		argument := arguments[index]
		convertedType := foundType
		if wasLocal {
			convertedType = argument
			foundLocal = true
		} else {
			if compatibleErr := CompatibleTypes(foundType, argument); compatibleErr != nil {
				return nil, compatibleErr
			}
		}
		convertedTypes = append(convertedTypes, convertedType)
	}

	if !foundLocal {
		return nil, fmt.Errorf("no local types, why did you call it? %v", record)
	}

	if len(arguments) > 0 {
		context := NewTypeParameterContextOther()
		for index, arg := range arguments {
			genericType := genericTypes[index]
			localType, wasLocal := genericType.(*LocalType)
			if wasLocal {
				//argument := arguments[index]
				context.SpecialSet(localType.Identifier().Name(), arg)
			}
		}
		return replaceRecordFromContext(record, context)
	}

	return record, nil
}

func callCustomType(customType *CustomTypeAtom, calledGenericTypes []dtype.Type) (dtype.Type, error) {
	generics := customType.parameters
	if len(generics) == 0 {
		return nil, fmt.Errorf("no calledGenericTypes for type %v", customType)
	}

	if len(generics) != len(calledGenericTypes) {
		return nil, fmt.Errorf("call custom type: illegal count %v %v %v", customType,
			calledGenericTypes, generics)
	}

	context := NewTypeParameterContextDynamic(customType.GenericNames())
	for index, genericType := range generics {
		localType, wasLocalType := genericType.(*LocalType)
		if wasLocalType {
			context.SpecialSet(localType.identifier.Name(), calledGenericTypes[index])
		}
	}

	return replaceCustomTypeFromContext(customType, context)
}

func callCustomTypeVariant(variant *CustomTypeVariantAtom, arguments []dtype.Type) (*CustomTypeVariantAtom, error) {
	params := variant.parameterFields

	if len(arguments) == 0 {
		panic("no arguments")
	}

	if len(params) != len(arguments) {
		return nil, fmt.Errorf("call custom type variant: illegal count %v %v %v", variant, arguments, params)
	}

	var convertedTypes []dtype.Type

	foundLocal := false
	for index, foundType := range variant.ParameterTypes() {
		_, wasLocal := foundType.(*LocalType)
		argument := arguments[index]
		convertedType := foundType
		if wasLocal {
			convertedType = argument
			foundLocal = true
		} else {
			if compatibleErr := CompatibleTypes(foundType, argument); compatibleErr != nil {
				return nil, compatibleErr
			}
		}
		convertedTypes = append(convertedTypes, convertedType)
	}

	if !foundLocal {
		//return nil, fmt.Errorf("no local types, why did you call it? %v", variant)
	}

	return NewCustomTypeVariant(variant.index, variant.inCustomType, variant.astCustomTypeVariant, convertedTypes), nil
}

func callPrimitiveType(primitive *PrimitiveAtom, arguments []dtype.Type) (*PrimitiveAtom, error) {
	genericTypes := primitive.GenericTypes()
	if len(arguments) == 0 {
		return nil, fmt.Errorf("no arguments for type %v", primitive)
	}

	if len(genericTypes) != len(arguments) {
		return nil, fmt.Errorf("call primitive type illegal count %v %v %v", primitive,
			genericTypes, arguments)
	}

	var convertedTypes []dtype.Type

	foundLocal := false
	for index, foundType := range genericTypes {
		_, wasLocal := foundType.(*LocalType)
		argument := arguments[index]
		convertedType := foundType
		if wasLocal {
			convertedType = argument
			foundLocal = true
		} else {
			if compatibleErr := CompatibleTypes(foundType, argument); compatibleErr != nil {
				return nil, compatibleErr
			}
		}
		convertedTypes = append(convertedTypes, convertedType)
	}

	if !foundLocal {
		return nil, fmt.Errorf("no local types, why did you call it? %v", primitive)
	}

	return NewPrimitiveType(primitive.name, convertedTypes), nil
}

func callTypeHelper(unaliasTarget dtype.Type, arguments []dtype.Type) (dtype.Type, error) {
	switch info := unaliasTarget.(type) {
	case *RecordAtom:
		{
			return callRecordType(info, arguments)
		}
	case *PrimitiveAtom:
		{
			return callPrimitiveType(info, arguments)
		}
	case *CustomTypeAtom:
		{
			return callCustomType(info, arguments)
		}
	case *CustomTypeVariantAtom:
		{
			return callCustomTypeVariant(info, arguments)
		}

	}
	return nil, fmt.Errorf("not found %T", unaliasTarget)
}

func CallType(target dtype.Type, arguments []dtype.Type) (dtype.Type, error) {
	unaliasTarget := Unalias(target)
	calledType, err := callTypeHelper(unaliasTarget, arguments)
	return calledType, err
}
