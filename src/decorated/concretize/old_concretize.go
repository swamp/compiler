package concretize

/*
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
	case *dectype.LocalTypeNameOnlyContext:
		return fmt.Errorf("found local type name context")
	case *dectype.LocalTypeNameReference:
		return fmt.Errorf("found local type name reference")
	case *dectype.LocalTypeNameOnlyContextReference:
		return fmt.Errorf("found local type name context reference")
	case *dectype.PrimitiveAtom:
		for _, x := range t.ParameterTypes() {
			if err := IsTypeCompletelyConcrete(x); err != nil {
				return err
			}
		}
		return nil
	case *dectype.CustomTypeAtom:
		for _, x := range t.Arguments() {
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

func ConcreteTypeIfNeeded(reference dtype.Type, concrete dtype.Type, resolveLocalTypeNames *dectype.ResolvedLocalTypeContext) (dtype.Type, decshared.DecoratedError) {
	if dectype.IsAny(concrete) {
		newReference, err := ResolveTypeFromContext(reference, resolveLocalTypeNames)
		if err != nil {
			return nil, decorated.NewInternalError(err)
		}

		if newReference == nil {
			panic(fmt.Errorf("newReference is nil"))
		}
		return newReference, nil
	}

	//log.Printf("concrete: \n %v\n %v\n    %v\n<-  %v\n (%v %v)", reference.HumanReadable(), concrete.HumanReadable(), reference, concrete, reference.FetchPositionLength().ToStandardReferenceString(), concrete.FetchPositionLength().ToCompleteReferenceString())

	switch t := reference.(type) {
	case *dectype.Alias:
		return ConcreteTypeIfNeeded(t.Next(), concrete, resolveLocalTypeNames)
	case *dectype.PrimitiveAtom:
		if typeIdRef, wasTypeIdRef := dectype.TryTypeIdRef(reference); wasTypeIdRef {
			//log.Printf("detected type id ref! %v", typeIdRef)
			pointingToRef := typeIdRef.ParameterTypes()[0]
			localTypeDefRef, wasRef := pointingToRef.(*dectype.ResolvedLocalTypeReference)
			if wasRef {
				concreteTypeRef, wasConcreteTypeRef := dectype.TryTypeIdRef(concrete)
				if wasConcreteTypeRef {
					concrete = concreteTypeRef.ParameterTypes()[0]
				}
				ref, err := resolveLocalTypeNames.SetType(localTypeDefRef.Identifier(), concrete)
				if err != nil {
					return nil, decorated.NewInternalError(err)
				}
				//log.Printf("ref:%v <- %v (%T)", localTypeDefRef.Identifier().Identifier().Name(), concrete.HumanReadable(), ref)
				return ref, nil
			} else {
				log.Printf("what is this %T", pointingToRef)
			}
		}
		//log.Printf("checking primitive %v <- %v", reference.HumanReadable(), concrete.HumanReadable())

		ref, wasRef := concrete.(*dectype.PrimitiveTypeReference)
		if !wasRef {
			return concrete, nil
		}
		return Primitive(t, ref, resolveLocalTypeNames)
	case *dectype.TupleTypeAtom:
		return Tuple(t, concrete.(*dectype.TupleTypeAtom), resolveLocalTypeNames)
	case *dectype.LocalTypeNameReference:
		return concrete, nil
	case *dectype.FunctionTypeReference:
		return FunctionTypeArg(t.FunctionAtom(), concrete.(*dectype.FunctionTypeReference).FunctionAtom(), resolveLocalTypeNames)
	case *dectype.ResolvedLocalTypeReference:
		return ConcreteTypeIfNeeded(t.ReferencedType(), concrete, resolveLocalTypeNames)
	case *dectype.ResolvedLocalType:
		return ConcreteTypeIfNeeded(t.ReferencedType(), concrete, resolveLocalTypeNames)
	case *dectype.FunctionAtom:
		return FunctionTypeArg(t, concrete.(*dectype.FunctionAtom), resolveLocalTypeNames)
	case *dectype.RecordAtom:
		return RecordArg(t, concrete.(*dectype.RecordAtom), resolveLocalTypeNames)
	case *dectype.PrimitiveTypeReference:
		return ConcreteTypeIfNeeded(t.Next(), concrete, resolveLocalTypeNames)
	case *dectype.AliasReference:
		return ConcreteTypeIfNeeded(t.Next(), concrete, resolveLocalTypeNames)
	case *dectype.UnmanagedType:
		return concrete, nil
	case *dectype.LocalTypeNameOnlyContext:
		return concrete, nil
	default:
		panic(fmt.Errorf("concrete: what is this %T", reference))
	}

	return concrete, nil
}

func ResolveTypeFromContext(parameterType dtype.Type, resolver *dectype.ResolvedLocalTypeContext) (dtype.Type, error) {
	resolvedType := parameterType

	//log.Printf("ResolveTypeFromContext:\n  %v\n  %v", parameterType.HumanReadable(), resolver)

	var err error
	switch t := parameterType.(type) {
	case *dectype.LocalTypeNameOnlyContext:
		return nil, fmt.Errorf("found local type name context")
	case *dectype.LocalTypeNameReference:
		resolvedType, err = resolver.LookupTypeRef(t)
		//log.Printf("resolved after settype to %T %v", resolvedType, resolvedType)
		if err != nil {
			log.Printf("ERR: %v", err)
			return nil, decorated.NewInternalError(err)
		}
		return resolvedType, nil
	case *dectype.LocalTypeNameOnlyContextReference:
		return nil, fmt.Errorf("found local type name context reference")
	case *dectype.PrimitiveAtom:
		return PrimitiveArguments(t, t.ParameterTypes(), resolver)
	case *dectype.FunctionAtom:
		return FunctionType(t, t.FunctionParameterTypes(), resolver)
	case *dectype.TupleTypeAtom:
		return TupleArgs(t, t.ParameterTypes(), resolver)
	case *dectype.UnmanagedType:
		return t, nil
	case *dectype.CustomTypeAtom:
		return CustomType(t, t.Arguments(), resolver)
	case *dectype.PrimitiveTypeReference:
		return t, nil
	default:
		next := t.Next()
		if next != nil && next != t {
			return ResolveTypeFromContext(next, resolver)
		}
		return nil, decorated.NewInternalError(fmt.Errorf("do not know what this is %T %v", parameterType, parameterType))
	}

	//log.Printf("resolved to %T", resolvedType)
	if err := IsTypeCompletelyConcrete(resolvedType); err != nil {
		log.Printf("error: %v", resolvedType)
		return nil, decorated.NewInternalError(err)
	}

	return resolvedType, nil
}

func FillLocalTypesFromSlice(references []dtype.Type, concretes []dtype.Type, resolver *dectype.ResolvedLocalTypeContext) ([]dtype.Type, decshared.DecoratedError) {
	var resolvedTypes []dtype.Type
	if len(concretes) != len(references) {
		return nil, decorated.NewInternalError(fmt.Errorf("must have equal number of arguments to concretize slices %v vs %v", concretes, references))
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

		if resolvedType == nil {
			panic(fmt.Errorf("how can resolvedType be nil %T %T %v", parameterType, argument, argument))
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

		if resolvedType == nil {
			panic(fmt.Errorf("how can the resolvedType be nil %T %v", argument, argument))
		}
		//log.Printf("resolved to %T", resolvedType)

		resolvedTypes = append(resolvedTypes, resolvedType)
	}

	return resolvedTypes, nil
}

func ResolveFromContext(references []dtype.Type, resolver *dectype.ResolvedLocalTypeContext) ([]dtype.Type, decshared.DecoratedError) {
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

func Primitive(reference *dectype.PrimitiveAtom, concrete_ *dectype.PrimitiveTypeReference, resolver *dectype.ResolvedLocalTypeContext) (*dectype.PrimitiveAtom, decshared.DecoratedError) {
	concrete, _ := concrete_.Next().(*dectype.PrimitiveAtom)

	return PrimitiveArguments(reference, concrete.ParameterTypes(), resolver)
}

func PrimitiveArguments(reference *dectype.PrimitiveAtom, arguments []dtype.Type, resolver *dectype.ResolvedLocalTypeContext) (*dectype.PrimitiveAtom, decshared.DecoratedError) {
	convertedTypes, err := FillLocalTypesFromSlice(reference.ParameterTypes(), arguments, resolver)
	if err != nil {
		return nil, err
	}
	//log.Printf("checking primitiveArguments resolved arguments: %v", convertedTypes)

	return dectype.NewPrimitiveType(reference.PrimitiveName(), convertedTypes), nil
}

func Tuple(reference *dectype.TupleTypeAtom, concrete *dectype.TupleTypeAtom, resolver *dectype.ResolvedLocalTypeContext) (*dectype.TupleTypeAtom, decshared.DecoratedError) {
	log.Printf("checking %v <- %v", reference, concrete)

	return TupleArgs(reference, concrete.ParameterTypes(), resolver)
}

func TupleArgs(reference *dectype.TupleTypeAtom, args []dtype.Type, resolver *dectype.ResolvedLocalTypeContext) (*dectype.TupleTypeAtom, decshared.DecoratedError) {
	convertedTypes, err := FillLocalTypesFromSlice(reference.ParameterTypes(), args, resolver)
	if err != nil {
		return nil, err
	}

	return dectype.NewTupleTypeAtom(reference.AstTuple(), convertedTypes), nil
}

func CustomTypeVariant(reference *dectype.CustomTypeVariantAtom, arguments []dtype.Type, resolver *dectype.ResolvedLocalTypeContext) (*dectype.CustomTypeVariantAtom, decshared.DecoratedError) {
	convertedTypes, err := FillLocalTypesFromSlice(reference.ParameterTypes(), arguments, resolver)
	if err != nil {
		return nil, err
	}

	newVariant := dectype.NewCustomTypeVariant(reference.Index(), nil, reference.AstCustomTypeVariant(), convertedTypes)

	return newVariant, nil
}

func CustomTypeVariantFromContext(reference *dectype.CustomTypeVariantAtom, resolver *dectype.ResolvedLocalTypeContext) (*dectype.CustomTypeVariantAtom, decshared.DecoratedError) {
	convertedTypes, err := ResolveFromContext(reference.ParameterTypes(), resolver)
	if err != nil {
		return nil, err
	}

	newVariant := dectype.NewCustomTypeVariant(reference.Index(), nil, reference.AstCustomTypeVariant(), convertedTypes)

	return newVariant, nil
}

func CustomType(reference *dectype.CustomTypeAtom, arguments []dtype.Type, resolver *dectype.ResolvedLocalTypeContext) (*dectype.CustomTypeAtom, decshared.DecoratedError) {
	//convertedTypes, err := FillLocalTypesFromSlice(reference.Arguments(), arguments, resolver)
	//if err != nil {
	//	return nil, err
	//}
	setErr := resolver.SetTypes(arguments)
	if setErr != nil {
		panic(setErr)
	}
	/*
		setErr := resolver.SetTypes(arguments)
		if setErr != nil {
			panic(setErr)
		}
		if !resolver.IsDefined() {
			panic(fmt.Errorf("it is not defined"))
		}
		&

	if !resolver.IsDefined() {
		//panic(fmt.Errorf("it is not defined"))
	}

	convertedTypes, err := FillLocalTypesFromSlice(reference.Arguments(), arguments, resolver)
	if err != nil {
		return nil, err
	}

	newCustomType := dectype.NewCustomTypePrepare(reference.AstCustomType(), convertedTypes, reference.ArtifactTypeName())

	var newVariants []*dectype.CustomTypeVariantAtom
	for _, variant := range reference.Variants() {
		variantParameterTypes, variantErr := ResolveFromContext(variant.ParameterTypes(), resolver)
		if variantErr != nil {
			return nil, variantErr
		}
		newVariant := dectype.NewCustomTypeVariant(variant.Index(), newCustomType, variant.AstCustomTypeVariant(), variantParameterTypes)
		newVariants = append(newVariants, newVariant)
	}

	//resolver.Debug()

	newCustomType.FinalizeVariants(newVariants)

	return newCustomType, nil
}

func CustomTypeHelper(reference *dectype.CustomTypeAtom, resolver *dectype.ResolvedLocalTypeContext) (*dectype.CustomTypeAtom, decshared.DecoratedError) {
	var decArguments []dtype.Type
	for _, argument := range reference.Arguments() {
		if localDef, wasLocalDef := dectype.TryLocalTypeDef(argument); wasLocalDef {
			decVariant, decVariantErr := resolver.LookupTypeName(localDef.TypeDefinition().Identifier().LocalTypeName().LocalType())
			if decVariantErr != nil {
				return nil, decorated.NewInternalError(decVariantErr)
			}
			decArguments = append(decArguments, decVariant)
		} else {
			decArguments = append(decArguments, argument)
		}
	}

	var decVariants []*dectype.CustomTypeVariantAtom

	for _, variantReference := range reference.Variants() {
		decVariant, decVariantErr := CustomTypeVariantFromContext(variantReference, resolver)
		if decVariantErr != nil {
			return nil, decVariantErr
		}
		decVariants = append(decVariants, decVariant)
	}

	newCustomType := dectype.NewCustomTypePrepare(reference.AstCustomType(), decArguments, reference.ArtifactTypeName())
	newCustomType.FinalizeVariants(decVariants)

	return newCustomType, nil
}

func Record(reference *dectype.RecordAtom, arguments []dtype.Type, resolver *dectype.ResolvedLocalTypeContext) (*dectype.RecordAtom, decshared.DecoratedError) {
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
		newField := dectype.NewRecordField(field.FieldName(), convertedTypes[index])
		newFields = append(newFields, newField)
	}

	newRecord := dectype.NewRecordType(reference.AstRecord(), newFields)

	log.Printf("newRecord:%v", newRecord)

	return newRecord, nil
}

func RecordArg(reference *dectype.RecordAtom, concrete *dectype.RecordAtom, resolver *dectype.ResolvedLocalTypeContext) (*dectype.RecordAtom, decshared.DecoratedError) {
	var fieldTypes []dtype.Type

	for _, field := range concrete.ParseOrderedFields() {
		fieldTypes = append(fieldTypes, field.Type())
	}

	return Record(reference, fieldTypes, resolver)
}

func FunctionType(reference *dectype.FunctionAtom, arguments []dtype.Type, resolver *dectype.ResolvedLocalTypeContext) (*dectype.FunctionAtom, decshared.DecoratedError) {
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

	convertedTypes, err := FillLocalTypesFromSlice(reference.FunctionParameterTypes(), arguments, resolver)
	if err != nil {
		return nil, err
	}

	newFunction := dectype.NewFunctionAtom(reference.AstFunction(), convertedTypes)

	return newFunction, nil
}

func FunctionTypeArg(reference *dectype.FunctionAtom, concrete *dectype.FunctionAtom, resolver *dectype.ResolvedLocalTypeContext) (*dectype.FunctionAtom, decshared.DecoratedError) {
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

func ConcretizeLocalTypeContextUsingArguments(localTypeNameContext *dectype.LocalTypeNameOnlyContext, concreteArguments []dtype.Type) (dtype.Type, decshared.DecoratedError) {
	if localTypeNameContext == nil {
		panic(fmt.Errorf("localTypeNameContext is nil"))
	}
	//log.Printf("ConcretizeLocalTypeContextUsingArguments %v '%v'", localTypeNameContext, dectype.TypesToHumanReadable(concreteArguments))
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
		concreteArguments = append(concreteArguments, dectype.NewAnyType())
		resolvedType, err = FunctionType(t.FunctionAtom(), concreteArguments, resolveLocalTypeNames)
		if err != nil {
			return nil, err
		}
	case *dectype.FunctionAtom:
		concreteArguments = append(concreteArguments, dectype.NewAnyType())
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

	/*
		if !resolveLocalTypeNames.IsDefined() {
			return nil, decorated.NewInternalError(fmt.Errorf("not all local type names where resolved, sorry about that \nNOT DEFINED: %v\n %T %v", resolveLocalTypeNames.DebugAllNotDefined(), localTypeNameContext.Next(), localTypeNameContext))
		}



	/*

		if err := IsTypeCompletelyConcrete(resolvedType); err != nil {
			log.Printf("error: %v", resolvedType)
			return nil, decorated.NewInternalError(err)
		}
	return resolvedType, nil
}


*/
