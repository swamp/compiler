package concretize

import (
	"fmt"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"log"
)

func IsTypeCompletelyConcrete(reference dtype.Type) error {
	switch t := reference.(type) {
	case *dectype.LocalTypeNameContext:
		return fmt.Errorf("found local type name context")
	case *dectype.LocalTypeNameReference:
		return fmt.Errorf("found local type name reference")
	case *dectype.LocalTypeNameContextReference:
		return fmt.Errorf("found local type name context reference")
	case *dectype.PrimitiveAtom:
		for _, x := range t.ParameterTypes() {
			if err := IsTypeCompletelyConcrete(x); err != nil {
				return err
			}
		}
		return nil
	case *dectype.FunctionAtom:
		for _, x := range t.FunctionParameterTypes() {
			if err := IsTypeCompletelyConcrete(x); err != nil {
				return err
			}
		}
		return nil

	case *dectype.TupleTypeAtom:
		for _, x := range t.ParameterTypes() {
			if err := IsTypeCompletelyConcrete(x); err != nil {
				return fmt.Errorf("had problem in function %v %w", t, err)
			}
		}
		return nil
	default:
		next := t.Next()
		if next != nil && next != t {
			return IsTypeCompletelyConcrete(next)
		}
		return nil
	}
}

func ConcreteTypeIfNeeded(reference dtype.Type, concrete dtype.Type, resolveLocalTypeNames *dectype.TypeParameterContext) (dtype.Type, decshared.DecoratedError) {
	if dectype.IsAny(concrete) {
		newReference, err := ResolveTypeFromContext(reference, resolveLocalTypeNames)
		if err != nil {
			return nil, decorated.NewInternalError(err)
		}
		return newReference, nil
	}

	switch t := reference.(type) {
	case *dectype.PrimitiveAtom:
		return Primitive(t, concrete.(*dectype.PrimitiveTypeReference), resolveLocalTypeNames)
	case *dectype.TupleTypeAtom:
		return Tuple(t, concrete.(*dectype.TupleTypeAtom), resolveLocalTypeNames)
	case *dectype.LocalTypeNameReference:
		return concrete, nil
	case *dectype.FunctionTypeReference:
		return FunctionTypeArg(t.FunctionAtom(), concrete.(*dectype.FunctionTypeReference).FunctionAtom(), resolveLocalTypeNames)
	case *dectype.LocalTypeDefinitionReference:
		return ConcreteTypeIfNeeded(t.ReferencedType(), concrete, resolveLocalTypeNames)
	case *dectype.LocalTypeDefinition:
		return ConcreteTypeIfNeeded(t.ReferencedType(), concrete, resolveLocalTypeNames)
	case *dectype.FunctionAtom:
		return FunctionTypeArg(t, concrete.(*dectype.FunctionAtom), resolveLocalTypeNames)
	case *dectype.RecordAtom:
		return RecordArg(t, concrete.(*dectype.RecordAtom), resolveLocalTypeNames)
	case *dectype.PrimitiveTypeReference:
		return ConcreteTypeIfNeeded(t.Next(), concrete, resolveLocalTypeNames)
	default:
		panic(fmt.Errorf("concrete: what is this %T", reference))
	}

	return concrete, nil
}

func ResolveTypeFromContext(parameterType dtype.Type, resolver *dectype.TypeParameterContext) (dtype.Type, error) {
	resolvedType := parameterType

	var err error
	switch t := parameterType.(type) {
	case *dectype.LocalTypeNameContext:
		return nil, fmt.Errorf("found local type name context")
	case *dectype.LocalTypeNameReference:
		resolvedType, err = resolver.LookupTypeRef(t)
		//log.Printf("resolved after settype to %T %v", resolvedType, resolvedType)
		if err != nil {
			log.Printf("ERR: %v", err)
			return nil, decorated.NewInternalError(err)
		}
		return resolvedType, nil
	case *dectype.LocalTypeNameContextReference:
		return nil, fmt.Errorf("found local type name context reference")
	case *dectype.PrimitiveAtom:
		return PrimitiveArguments(t, t.ParameterTypes(), resolver)
	case *dectype.FunctionAtom:
		return FunctionType(t, t.FunctionParameterTypes(), resolver)
	case *dectype.TupleTypeAtom:
		return TupleArgs(t, t.ParameterTypes(), resolver)
	default:
		next := t.Next()
		if next != nil && next != t {
			return ResolveTypeFromContext(next, resolver)
		}
		return nil, nil
	}

	//log.Printf("resolved to %T", resolvedType)
	if err := IsTypeCompletelyConcrete(resolvedType); err != nil {
		log.Printf("error: %v", resolvedType)
		return nil, decorated.NewInternalError(err)
	}

	return resolvedType, nil
}

func ResolveSlices(references []dtype.Type, concretes []dtype.Type, resolver *dectype.TypeParameterContext) ([]dtype.Type, decshared.DecoratedError) {
	var resolvedTypes []dtype.Type
	if len(concretes) != len(references) {
		return nil, decorated.NewInternalError(fmt.Errorf("must have equal number of arguments to concretize slices"))
	}

	for index, parameterType := range references {
		resolvedType := parameterType
		argument := concretes[index]

		var lookupErr decshared.DecoratedError
		resolvedType, lookupErr = ConcreteTypeIfNeeded(parameterType, argument, resolver)
		if lookupErr != nil {
			log.Printf("ERR: %v", lookupErr)
			return nil, lookupErr
		}

		localTypeRef, wasLocalTypeRef := parameterType.(*dectype.LocalTypeNameReference)
		if wasLocalTypeRef {
			var err error
			resolvedType, err = resolver.SetType(localTypeRef, argument)
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

func ResolveFromContext(references []dtype.Type, resolver *dectype.TypeParameterContext) ([]dtype.Type, decshared.DecoratedError) {
	var resolvedTypes []dtype.Type
	//log.Printf("checking %d types", len(references))
	for _, parameterType := range references {
		resolvedType, err := ResolveTypeFromContext(parameterType, resolver)
		if err != nil {
			return nil, decorated.NewInternalError(err)
		}
		resolvedTypes = append(resolvedTypes, resolvedType)
	}

	return resolvedTypes, nil
}

func Primitive(reference *dectype.PrimitiveAtom, concrete_ *dectype.PrimitiveTypeReference, resolver *dectype.TypeParameterContext) (*dectype.PrimitiveAtom, decshared.DecoratedError) {
	concrete, _ := concrete_.Next().(*dectype.PrimitiveAtom)
	log.Printf("checking primitive %v <- paramters: %v", reference, concrete.ParameterTypes())

	return PrimitiveArguments(reference, concrete.ParameterTypes(), resolver)
}

func PrimitiveArguments(reference *dectype.PrimitiveAtom, arguments []dtype.Type, resolver *dectype.TypeParameterContext) (*dectype.PrimitiveAtom, decshared.DecoratedError) {
	//log.Printf("checking primitiveArguments %v <- %v", reference, arguments)

	convertedTypes, err := ResolveSlices(reference.ParameterTypes(), arguments, resolver)
	if err != nil {
		return nil, err
	}
	//log.Printf("checking primitiveArguments resolved arguments: %v", convertedTypes)

	return dectype.NewPrimitiveType(reference.PrimitiveName(), convertedTypes), nil
}

func Tuple(reference *dectype.TupleTypeAtom, concrete *dectype.TupleTypeAtom, resolver *dectype.TypeParameterContext) (*dectype.TupleTypeAtom, decshared.DecoratedError) {
	log.Printf("checking %v <- %v", reference, concrete)

	return TupleArgs(reference, concrete.ParameterTypes(), resolver)
}

func TupleArgs(reference *dectype.TupleTypeAtom, args []dtype.Type, resolver *dectype.TypeParameterContext) (*dectype.TupleTypeAtom, decshared.DecoratedError) {
	convertedTypes, err := ResolveSlices(reference.ParameterTypes(), args, resolver)
	if err != nil {
		return nil, err
	}

	return dectype.NewTupleTypeAtom(reference.AstTuple(), convertedTypes), nil
}

func CustomTypeVariant(reference *dectype.CustomTypeVariantAtom, arguments []dtype.Type, resolver *dectype.TypeParameterContext) (*dectype.CustomTypeVariantAtom, decshared.DecoratedError) {
	convertedTypes, err := ResolveSlices(reference.ParameterTypes(), arguments, resolver)
	if err != nil {
		return nil, err
	}

	newVariant := dectype.NewCustomTypeVariant(reference.Index(), nil, reference.AstCustomTypeVariant(), convertedTypes)

	return newVariant, nil
}

func CustomTypeVariantFromContext(reference *dectype.CustomTypeVariantAtom, resolver *dectype.TypeParameterContext) (*dectype.CustomTypeVariantAtom, decshared.DecoratedError) {
	convertedTypes, err := ResolveFromContext(reference.ParameterTypes(), resolver)
	if err != nil {
		return nil, err
	}

	newVariant := dectype.NewCustomTypeVariant(reference.Index(), nil, reference.AstCustomTypeVariant(), convertedTypes)

	return newVariant, nil
}

func CustomType(reference *dectype.CustomTypeAtom, arguments []dtype.Type, resolver *dectype.TypeParameterContext) (*dectype.CustomTypeAtom, decshared.DecoratedError) {
	var decVariants []*dectype.CustomTypeVariantAtom

	setErr := resolver.SetTypes(arguments)
	if setErr != nil {
		panic(setErr)
	}
	if !resolver.IsDefined() {
		panic(fmt.Errorf("it is not defined"))
	}

	//resolver.Debug()

	for _, variantReference := range reference.Variants() {
		decVariant, decVariantErr := CustomTypeVariantFromContext(variantReference, resolver)
		if decVariantErr != nil {
			return nil, decVariantErr
		}
		decVariants = append(decVariants, decVariant)
	}

	newCustomType := dectype.NewCustomTypePrepare(reference.AstCustomType(), reference.ArtifactTypeName())
	newCustomType.FinalizeVariants(decVariants)

	return newCustomType, nil
}

func Record(reference *dectype.RecordAtom, arguments []dtype.Type, resolver *dectype.TypeParameterContext) (*dectype.RecordAtom, decshared.DecoratedError) {
	var fieldTypes []dtype.Type

	setErr := resolver.SetTypes(arguments)
	if setErr != nil {
		panic(setErr)
	}

	for _, field := range reference.ParseOrderedFields() {
		fieldTypes = append(fieldTypes, field.Type())
	}

	convertedTypes, err := ResolveFromContext(fieldTypes, resolver)
	if err != nil {
		return nil, err
	}

	var newFields []*dectype.RecordField
	for index, field := range reference.ParseOrderedFields() {
		newField := dectype.NewRecordField(field.FieldName(), field.AstRecordTypeField(), convertedTypes[index])
		newFields = append(newFields, newField)
	}

	newRecord := dectype.NewRecordType(reference.AstRecord(), newFields)

	return newRecord, nil
}

func RecordArg(reference *dectype.RecordAtom, concrete *dectype.RecordAtom, resolver *dectype.TypeParameterContext) (*dectype.RecordAtom, decshared.DecoratedError) {
	var fieldTypes []dtype.Type

	for _, field := range concrete.ParseOrderedFields() {
		fieldTypes = append(fieldTypes, field.Type())
	}

	return Record(reference, fieldTypes, resolver)
}

func FunctionType(reference *dectype.FunctionAtom, arguments []dtype.Type, resolver *dectype.TypeParameterContext) (*dectype.FunctionAtom, decshared.DecoratedError) {
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

		return FunctionType(reference, allConverted, resolver)
	} else {
		if len(reference.FunctionParameterTypes()) < len(arguments) {
			return nil, decorated.NewInternalError(fmt.Errorf("too few parameter types"))
		}
	}

	arguments = append(arguments, dectype.NewAnyType())

	convertedTypes, err := ResolveSlices(reference.FunctionParameterTypes(), arguments, resolver)
	if err != nil {
		return nil, err
	}

	newFunction := dectype.NewFunctionAtom(reference.AstFunction(), convertedTypes)

	return newFunction, nil
}

func FunctionTypeArg(reference *dectype.FunctionAtom, concrete *dectype.FunctionAtom, resolver *dectype.TypeParameterContext) (*dectype.FunctionAtom, decshared.DecoratedError) {
	log.Printf("%v \n<- %v", reference, concrete)
	return FunctionType(reference, concrete.FunctionParameterTypes(), resolver)
}

func Concrete(reference dtype.Type, concrete dtype.Type) (dtype.Type, decshared.DecoratedError) {
	resolveLocalTypeNames := dectype.NewTypeParameterContext()
	concreteType, resolveErr := ConcreteTypeIfNeeded(reference, concrete, resolveLocalTypeNames)
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
	case *dectype.CustomTypeAtom:
		resolvedType, err = CustomType(t, concreteArguments, resolveLocalTypeNames)
		if err != nil {
			return nil, err
		}
	case *dectype.CustomTypeVariantAtom:
		resolvedType, err = CustomTypeVariant(t, concreteArguments, resolveLocalTypeNames)
		if err != nil {
			return nil, err
		}
	case *dectype.FunctionTypeReference:
		resolvedType, err = FunctionType(t.FunctionAtom(), concreteArguments, resolveLocalTypeNames)
		if err != nil {
			return nil, err
		}
	case *dectype.FunctionAtom:
		resolvedType, err = FunctionType(t, concreteArguments, resolveLocalTypeNames)
		if err != nil {
			return nil, err
		}
	case *dectype.PrimitiveAtom:
		resolvedType, err = PrimitiveArguments(t, concreteArguments, resolveLocalTypeNames)
		if err != nil {
			return nil, err
		}
	case *dectype.RecordAtom:
		resolvedType, err = Record(t, concreteArguments, resolveLocalTypeNames)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Errorf("not handled concrete %T", localTypeNameContext.Next()))
	}

	if !resolveLocalTypeNames.IsDefined() {
		return nil, decorated.NewInternalError(fmt.Errorf("not all local type names where resolved, sorry about that %T %v", localTypeNameContext.Next(), localTypeNameContext))
	}

	/*

		if err := IsTypeCompletelyConcrete(resolvedType); err != nil {
			log.Printf("error: %v", resolvedType)
			return nil, decorated.NewInternalError(err)
		}
	*/
	return resolvedType, nil
}
