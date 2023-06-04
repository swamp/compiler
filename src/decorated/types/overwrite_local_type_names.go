package dectype

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/decorated/debug"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
)

type LookupTypeName interface {
	LookupTypeRef(decReference *LocalTypeNameReference) (*ResolvedLocalTypeReference, decshared.DecoratedError)
	DebugString() string
}

func replaceLocalNamesInFunctionIfNeeded(atom *FunctionAtom, lookup LookupTypeName) (*FunctionAtom, bool, error) {
	newTypes, wasReplaced, err := replaceLocalNameInSliceIfNeeded(atom.FunctionParameterTypes(), lookup)
	if err != nil {
		return nil, false, err
	}
	if !wasReplaced {
		return atom, false, nil
	}

	return NewFunctionAtom(atom.astFunctionType, newTypes), true, nil
}

func replaceLocalNamesInTupleTypeIfNeeded(atom *TupleTypeAtom, lookup LookupTypeName) (*TupleTypeAtom, bool, error) {
	newTypes, wasReplaced, err := replaceLocalNameInSliceIfNeeded(atom.ParameterTypes(), lookup)
	if err != nil {
		return nil, false, err
	}
	if !wasReplaced {
		return atom, false, nil
	}

	return NewTupleTypeAtom(atom.astTupleType, newTypes), true, nil
}

func replaceLocalNamesInCustomTypeIfNeeded(customType *CustomTypeAtom, lookup LookupTypeName) (*CustomTypeAtom, bool,
	error) {
	var newVariants []*CustomTypeVariantAtom

	newCustomType := NewCustomTypePrepare(customType.astCustomType, customType.artifactTypeName)

	someVariantWasReplaced := false
	for _, variant := range customType.variants {
		newTypes, wasReplaced, err := replaceLocalNameInSliceIfNeeded(variant.ParameterTypes(), lookup)
		if err != nil {
			return nil, false, err
		}

		if !wasReplaced {
			newVariants = append(newVariants, variant)
			continue
		}

		someVariantWasReplaced = true
		newVariant := NewCustomTypeVariant(variant.index, newCustomType, variant.astCustomTypeVariant, newTypes)
		newVariants = append(newVariants, newVariant)
	}

	if !someVariantWasReplaced {
		return customType, someVariantWasReplaced, nil
	}

	newCustomType.FinalizeVariants(newVariants)

	return newCustomType, someVariantWasReplaced, nil
}

func replaceLocalNamesInRecordAtomIfNeeded(recordType *RecordAtom, lookup LookupTypeName) (*RecordAtom, bool, error) {
	var newRecordFields []*RecordField

	//newCustomType := NewCustomTypePrepare(recordType.astCustomType, recordType.artifactTypeName)

	someFieldWasReplaced := false
	for _, recordField := range recordType.sortedFields {
		newType, wasReplaced, err := internalCollapse(recordField.Type(), lookup)
		if err != nil {
			return nil, false, err
		}

		if !wasReplaced {
			newRecordFields = append(newRecordFields, recordField)
			continue
		}

		someFieldWasReplaced = true
		newVariant := NewRecordField(recordField.FieldName(), newType)
		newRecordFields = append(newRecordFields, newVariant)
	}

	if !someFieldWasReplaced {
		return recordType, false, nil
	}

	newRecord := NewRecordType(recordType.AstRecord(), newRecordFields)

	return newRecord, true, nil
}

func replaceLocalNamesInPrimitiveIfNeeded(atom *PrimitiveAtom, lookup LookupTypeName) (*PrimitiveAtom, bool, error) {
	newTypes, wasReplaced, err := replaceLocalNameInSliceIfNeeded(atom.ParameterTypes(), lookup)
	if err != nil {
		return nil, false, err
	}
	if !wasReplaced {
		return atom, false, nil
	}

	return NewPrimitiveType(atom.PrimitiveName(), newTypes), wasReplaced, nil
}

func Collapse(typeToCheck dtype.Type, lookup LookupTypeName) (dtype.Type, error) {
	newType, _, err := internalCollapse(typeToCheck, lookup)
	log.Printf("collapsing %T %s\n%s\nTo\n%s", typeToCheck, lookup.DebugString(), debug.TreeString(typeToCheck),
		debug.TreeString(newType))

	return newType, err
}

func internalCollapse(typeToCheck dtype.Type, lookup LookupTypeName) (dtype.Type, bool, error) {
	switch t := typeToCheck.(type) {
	case *LocalTypeNameReference:
		newType, err := lookup.LookupTypeRef(t)
		return newType, true, err
	case *FunctionAtom:
		return replaceLocalNamesInFunctionIfNeeded(t, lookup)
	case *TupleTypeAtom:
		return replaceLocalNamesInTupleTypeIfNeeded(t, lookup)
	case *CustomTypeAtom:
		return replaceLocalNamesInCustomTypeIfNeeded(t, lookup)
	case *RecordAtom:
		return replaceLocalNamesInRecordAtomIfNeeded(t, lookup)
	case *PrimitiveTypeReference:
		return replaceLocalNamesInPrimitiveIfNeeded(t.primitiveType, lookup)
	case *PrimitiveAtom:
		return replaceLocalNamesInPrimitiveIfNeeded(t, lookup)
	case *LocalTypeNameOnlyContextReference:
		return internalCollapse(t.Next(), lookup)
	case *AliasReference:
		return internalCollapse(t.Next(), lookup)
	case *Alias:
		return internalCollapse(t.Next(), lookup)
	case *UnmanagedType:
		return typeToCheck, false, nil
	//return replaceLocalNameInNameOnlyContext(t, lookup)
	case *LocalTypeNameOnlyContext:
		return internalCollapse(t.Next(), lookup)
	case *ResolvedLocalTypeContext:
		return internalCollapse(t.Next(), t)
	default:
		panic(fmt.Errorf("collapse not implemented for %T", typeToCheck))
		return typeToCheck, false, nil
	}
}

func replaceLocalNameInSliceIfNeeded(types []dtype.Type, lookup LookupTypeName) ([]dtype.Type, bool, error) {
	someOneWasReplaced := false
	var newTypes []dtype.Type
	for _, x := range types {
		replaced, err := Collapse(x, lookup)
		if err != nil {
			return nil, true, err
		}
		if replaced != x {
			someOneWasReplaced = true
		}
		newTypes = append(newTypes, replaced)
	}

	return newTypes, someOneWasReplaced, nil
}
