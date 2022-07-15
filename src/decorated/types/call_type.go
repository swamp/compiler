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
	case *InvokerType:
		return replaceInvokerTypeFromContext(info, lookup)
	case *CustomTypeAtom:
		// TODO: fix this
		return originalTarget, nil
	case *CustomTypeReference:
		// TODO: fix this
		return originalTarget, nil
	default:
		log.Printf("warning: not sure what to do with %T %v. Returning same type for now", originalTarget, originalTarget)
		return nil, fmt.Errorf("not sure what to do with %T %v. Returning same type for now", target, target)
	}

	// fmt.Printf("warning: not sure what to do with %T %v. Returning same type for now\n", target, target)

	return originalTarget, nil
}

func replaceRecordFromContext(record *RecordAtom, lookup Lookup) (*RecordAtom, error) {
	var replacedFields []*RecordField

	for _, field := range record.SortedFields() {
		converted, convertedErr := ReplaceTypeFromContext(field.Type(), lookup)
		if convertedErr != nil {
			return nil, convertedErr
		}
		if converted == nil {
			panic("converted is nil")
		}

		newField := NewRecordField(field.name, field.AstRecordTypeField(), converted)

		replacedFields = append(replacedFields, newField)
	}

	return NewRecordType(record.AstRecord(), replacedFields, nil), nil
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
	var replacedVariants []*CustomTypeVariant

	for index, field := range customType.Variants() {
		var variantParameters []dtype.Type
		for _, param := range field.parameterTypes {
			converted, convertedErr := ReplaceTypeFromContext(param, lookup)
			if convertedErr != nil {
				return nil, convertedErr
			}
			variantParameters = append(variantParameters, converted)
		}

		NewCustomTypeVariantReference(nil, field)
		newField := NewCustomTypeVariant(index, field.astCustomTypeVariant, variantParameters)

		replacedVariants = append(replacedVariants, newField)
	}

	return NewCustomType(customType.astCustomType, customType.artifactTypeName, nil, replacedVariants), nil
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

	return NewRecordType(record.AstRecord(), record.SortedFields(), convertedTypes), nil
}

func callCustomType(customType *CustomTypeAtom, arguments []dtype.Type) (dtype.Type, error) {
	argumentNames := customType.genericLocalTypeNames
	if len(argumentNames) == 0 {
		return nil, fmt.Errorf("no arguments for type %v", customType)
	}

	if len(argumentNames) != len(arguments) {
		return nil, fmt.Errorf("call custom type: illegal count %v %v %v", customType,
			arguments, argumentNames)
	}

	context := NewTypeParameterContextDynamic(argumentNames)
	for index, name := range argumentNames {
		context.SpecialSet(name.Name(), arguments[index])
	}

	return replaceCustomTypeFromContext(customType, context)
}

func callCustomTypeVariant(variant *CustomTypeVariant, arguments []dtype.Type) (dtype.Type, error) {
	params := variant.parameterTypes

	if len(params) != len(arguments) {
		return nil, fmt.Errorf("call custom type variant: illegal count %v %v %v", variant, arguments, params)
	}

	argumentNames := variant.InCustomType().genericLocalTypeNames
	context := NewTypeParameterContextDynamic(argumentNames)

	for index, param := range params {
		localType, wasLocalType := param.(*LocalType)
		if wasLocalType {
			context.SpecialSet(localType.identifier.Name(), arguments[index])
		}
	}

	context.FillOutTheRestWithAny()
	return replaceCustomTypeFromContext(variant.InCustomType(), context)
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

func CallType(target dtype.Type, arguments []dtype.Type) (dtype.Type, error) {
	unaliasTarget := Unalias(target)
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
	case *CustomTypeVariant:
		{
			return callCustomTypeVariant(info, arguments)
		}
	}

	return nil, fmt.Errorf(" noot found %T", unaliasTarget)
}
