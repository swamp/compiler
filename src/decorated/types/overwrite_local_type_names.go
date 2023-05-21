package dectype

import (
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
)

type LookupTypeName interface {
	LookupTypeRef(decReference *LocalTypeNameReference) (*ResolvedLocalTypeReference, decshared.DecoratedError)
}

func replaceLocalNamesInFunctionIfNeeded(atom *FunctionAtom, lookup LookupTypeName) (*FunctionAtom, error) {
	newTypes, wasReplaced, err := replaceLocalNameInSliceIfNeeded(atom.FunctionParameterTypes(), lookup)
	if err != nil {
		return nil, err
	}
	if !wasReplaced {
		return atom, nil
	}

	return NewFunctionAtom(atom.astFunctionType, newTypes), nil
}

func replaceLocalNamesInPrimitiveIfNeeded(atom *PrimitiveAtom, lookup LookupTypeName) (*PrimitiveAtom, error) {
	newTypes, wasReplaced, err := replaceLocalNameInSliceIfNeeded(atom.ParameterTypes(), lookup)
	if err != nil {
		return nil, err
	}
	if !wasReplaced {
		return atom, nil
	}

	return NewPrimitiveType(atom.PrimitiveName(), newTypes), nil
}

func ReplaceLocalNameIfNeeded(typeToCheck dtype.Type, lookup LookupTypeName) (dtype.Type, error) {
	switch t := typeToCheck.(type) {
	case *LocalTypeNameReference:
		return lookup.LookupTypeRef(t)
	case *FunctionAtom:
		return replaceLocalNamesInFunctionIfNeeded(t, lookup)
	case *PrimitiveAtom:
		return replaceLocalNamesInPrimitiveIfNeeded(t, lookup)
	default:
		return typeToCheck, nil
	}
}

func replaceLocalNameInSliceIfNeeded(types []dtype.Type, lookup LookupTypeName) ([]dtype.Type, bool, error) {
	someOneWasReplaced := false
	var newTypes []dtype.Type
	for _, x := range types {
		replaced, err := ReplaceLocalNameIfNeeded(x, lookup)
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
