package concretize

import (
	"fmt"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"log"
)

func IfNeeded(reference dtype.Type, concrete dtype.Type, resolveLocalTypeNames *dectype.TypeParameterContext) (dtype.Type, decshared.DecoratedError) {
	switch t := reference.(type) {
	case *dectype.PrimitiveAtom:
		return Primitive(t, concrete.(*dectype.PrimitiveAtom), resolveLocalTypeNames)
	case *dectype.TupleTypeAtom:
		return Tuple(t, concrete.(*dectype.TupleTypeAtom), resolveLocalTypeNames)
	case *dectype.LocalTypeNameReference:
		return concrete, nil
	default:
		log.Printf("what is this %T", reference)
	}

	return concrete, nil
}

func ResolveSlices(references []dtype.Type, concretes []dtype.Type, resolver *dectype.TypeParameterContext) ([]dtype.Type, decshared.DecoratedError) {
	var resolvedTypes []dtype.Type
	if len(concretes) != len(references) {
		return nil, decorated.NewInternalError(fmt.Errorf("must have equal number of arguments to concretize custom type variant"))
	}
	log.Printf("checking %d types", len(references))
	for index, parameterType := range references {
		resolvedType := parameterType

		argument := concretes[index]
		log.Printf("checking index %d, %v <- %v", index, parameterType, argument)
		var lookupErr decshared.DecoratedError
		resolvedType, lookupErr = IfNeeded(parameterType, argument, resolver)
		if lookupErr != nil {
			log.Printf("ERR: %v", lookupErr)
			return nil, lookupErr
		}

		localTypeRef, wasLocalTypeRef := parameterType.(*dectype.LocalTypeNameReference)
		if wasLocalTypeRef {
			var err error
			resolvedType, err = resolver.SetType(localTypeRef, resolvedType)
			//log.Printf("resolved after settype to %T %v", resolvedType, resolvedType)
			if err != nil {
				log.Printf("ERR: %v", err)
				return nil, decorated.NewInternalError(err)
			}
		}
		//log.Printf("resolved to %T", resolvedType)

		resolvedTypes = append(resolvedTypes, resolvedType)
	}

	return resolvedTypes, nil
}

func Primitive(reference *dectype.PrimitiveAtom, concrete *dectype.PrimitiveAtom, resolver *dectype.TypeParameterContext) (*dectype.PrimitiveAtom, decshared.DecoratedError) {
	log.Printf("checking %v <- %v", reference, concrete)

	convertedTypes, err := ResolveSlices(reference.ParameterTypes(), concrete.ParameterTypes(), resolver)
	if err != nil {
		return nil, err
	}

	return dectype.NewPrimitiveType(concrete.PrimitiveName(), convertedTypes), nil
}

func Tuple(reference *dectype.TupleTypeAtom, concrete *dectype.TupleTypeAtom, resolver *dectype.TypeParameterContext) (*dectype.TupleTypeAtom, decshared.DecoratedError) {
	log.Printf("checking %v <- %v", reference, concrete)

	convertedTypes, err := ResolveSlices(reference.ParameterTypes(), concrete.ParameterTypes(), resolver)
	if err != nil {
		return nil, err
	}

	return dectype.NewTupleTypeAtom(concrete.AstTuple(), convertedTypes), nil
}

func CustomTypeVariant(reference *dectype.CustomTypeVariantAtom, arguments []dtype.Type, resolver *dectype.TypeParameterContext) (*dectype.CustomTypeVariantAtom, decshared.DecoratedError) {
	convertedTypes, err := ResolveSlices(reference.ParameterTypes(), arguments, resolver)
	if err != nil {
		return nil, err
	}

	newVariant := dectype.NewCustomTypeVariant(reference.Index(), nil, reference.AstCustomTypeVariant(), convertedTypes)

	return newVariant, nil
}

func Function(reference *dectype.FunctionAtom, arguments []dtype.Type, resolver *dectype.TypeParameterContext) (*dectype.FunctionAtom, decshared.DecoratedError) {
	if hasAnyMatching, startIndex := dectype.HasAnyMatchingTypes(reference.FunctionParameterTypes()); hasAnyMatching {
		originalInitialCount := startIndex
		originalEndCount := len(reference.FunctionParameterTypes()) - startIndex - 2

		originalFirst := append([]dtype.Type{}, reference.FunctionParameterTypes()[0:startIndex]...)

		if len(arguments) >= len(reference.FunctionParameterTypes()) {
			originalEndCount++
		}

		otherMiddle := arguments[originalInitialCount : len(arguments)-originalEndCount]
		if len(otherMiddle) < 1 {
			return nil, decorated.NewInternalError(fmt.Errorf("currently, you must have at least one wildcard parameter"))
		}

		originalEnd := reference.FunctionParameterTypes()[startIndex+1 : len(reference.FunctionParameterTypes())]

		allConverted := append(originalFirst, otherMiddle...)
		allConverted = append(allConverted, originalEnd...)

		//created := dectype.NewFunctionAtom(reference.AstFunction(), allConverted)

		return Function(reference, allConverted, resolver)
	} else {
		if len(reference.FunctionParameterTypes()) < len(arguments) {
			return nil, decorated.NewInternalError(fmt.Errorf("too few parameter types"))
		}
	}

	convertedTypes, err := ResolveSlices(reference.FunctionParameterTypes(), arguments, resolver)
	if err != nil {
		return nil, err
	}

	newFunction := dectype.NewFunctionAtom(reference.AstFunction(), convertedTypes)

	return newFunction, nil
}

func Concrete(reference dtype.Type, concrete dtype.Type) (dtype.Type, decshared.DecoratedError) {
	resolveLocalTypeNames := dectype.NewTypeParameterContext()
	concreteType, resolveErr := IfNeeded(reference, concrete, resolveLocalTypeNames)
	if resolveErr != nil {
		return nil, resolveErr
	}

	if !resolveLocalTypeNames.IsDefined() {
		return nil, decorated.NewInternalError(fmt.Errorf("not all local type names where resolved, sorry about that"))
	}

	return concreteType, nil
}

func ConcreteArguments(localTypeNameContext *dectype.LocalTypeNameContext, concreteArguments []dtype.Type) (dtype.Type, decshared.DecoratedError) {
	resolveLocalTypeNames := dectype.NewTypeParameterContext()
	resolveLocalTypeNames.AddExpectedDefs(localTypeNameContext.Names())

	var err decshared.DecoratedError
	var resolvedType dtype.Type

	switch t := localTypeNameContext.Next().(type) {
	case *dectype.CustomTypeVariantAtom:
		resolvedType, err = CustomTypeVariant(t, concreteArguments, resolveLocalTypeNames)
		if err != nil {
			return nil, err
		}
	case *dectype.FunctionAtom:
		resolvedType, err = Function(t, concreteArguments, resolveLocalTypeNames)
		if err != nil {
			return nil, err
		}
	}

	if !resolveLocalTypeNames.IsDefined() {
		return nil, decorated.NewInternalError(fmt.Errorf("not all local type names where resolved, sorry about that"))
	}

	return resolvedType, nil
}
