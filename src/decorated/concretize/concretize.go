package concretize

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/decorated/debug"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

/*
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

case *dectype.ResolvedLocalTypeReference:

	return ConcreteTypeIfNeeded(t.ReferencedType(), concrete, resolveLocalTypeNames)

case *dectype.ResolvedLocalType:

	return ConcreteTypeIfNeeded(t.ReferencedType(), concrete, resolveLocalTypeNames)

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
*/

func FillContextFromPrimitive(primitiveAtom *dectype.PrimitiveAtom, concretes []dtype.Type,
	resolveLocalTypeNames dectype.DynamicResolver) decshared.DecoratedError {
	return FillLocalTypesFromSlice(primitiveAtom.ParameterTypes(), concretes, resolveLocalTypeNames)
}

func fillContextFromLocalContext(localContext *dectype.LocalTypeNameOnlyContext,
	concrete *dectype.ResolvedLocalTypeContext,
	resolver dectype.DynamicResolver) decshared.DecoratedError {
	if len(localContext.Names()) != len(concrete.Definitions()) {
		return decorated.NewInternalError(fmt.Errorf("not same definitions %d vs %d", len(localContext.Names()),
			len(concrete.Definitions())))
	}

	for _, resolvedDef := range concrete.Definitions() {
		resolver.SetType(resolvedDef.Identifier().LocalTypeName(), resolvedDef.ReferencedType())
	}

	return nil
}

func ConcreteTypeIfNeeded(reference dtype.Type, concrete dtype.Type,
	resolveLocalTypeNames dectype.DynamicResolver) decshared.DecoratedError {

	//log.Printf("concrete: \n %v\n %v\n    %v\n<-  %v\n (%v %v)", reference.HumanReadable(), concrete.HumanReadable(), reference, concrete, reference.FetchPositionLength().ToStandardReferenceString(), concrete.FetchPositionLength().ToCompleteReferenceString())

	if dectype.IsAny(concrete) {
		return nil
	}

	switch t := reference.(type) {
	case *dectype.FunctionTypeReference:
		return FillContextFromFunction(
			t.FunctionAtom(), concrete.(*dectype.FunctionTypeReference).FunctionAtom().FunctionParameterTypes(),
			resolveLocalTypeNames,
		)
	case *dectype.FunctionAtom:
		return FillContextFromFunction(
			t, concrete.(*dectype.FunctionAtom).FunctionParameterTypes(), resolveLocalTypeNames,
		)
	case *dectype.TupleTypeAtom:
		return FillContextFromTuple(
			t, concrete.(*dectype.TupleTypeAtom).ParameterTypes(), resolveLocalTypeNames,
		)
	case *dectype.RecordAtom:
		unaliasRecordAtom := dectype.Unalias(concrete)
		return FillContextFromRecordAtom(t, unaliasRecordAtom.(*dectype.RecordAtom), resolveLocalTypeNames)
	case *dectype.ResolvedLocalTypeContext:
		reverseLookup := make(map[string]*dectype.LocalTypeNameReference)
		for _, mappedDefinition := range t.Definitions() {
			localTypeNameRef, wasLocalTypeNameRef := mappedDefinition.ReferencedType().(*dectype.LocalTypeNameReference)
			if wasLocalTypeNameRef {
				reverseLookup[mappedDefinition.Identifier().LocalTypeName().Name()] = localTypeNameRef
			}
		}

		subResolver := dectype.NewDynamicReverseResolver(reverseLookup, resolveLocalTypeNames)
		return ConcreteTypeIfNeeded(t.Next(), concrete, subResolver)
	case *dectype.PrimitiveTypeReference:
		concreteUnalias := dectype.Unalias(concrete)
		concreteAtom, _ := concreteUnalias.(*dectype.PrimitiveAtom)
		if concreteAtom == nil {
			ref, _ := concreteUnalias.(*dectype.PrimitiveTypeReference)
			if ref != nil {
				concreteAtom = ref.PrimitiveAtom()
			} else {
				panic(fmt.Errorf("can not get to the concrete primitive atom %T", concrete))
			}
		}
		return FillContextFromPrimitive(t.PrimitiveAtom(), concreteAtom.ParameterTypes(), resolveLocalTypeNames)

	case *dectype.LocalTypeNameReference:
		if !dectype.IsAny(concrete) {
			resolveLocalTypeNames.SetType(t.LocalTypeName(), concrete)
		}
	case *dectype.LocalTypeNameOnlyContextReference:
		resolvedContext, _ := concrete.(*dectype.ResolvedLocalTypeContext)
		if resolvedContext != nil {
			fillContextFromLocalContext(t.LocalTypeNameContext(), resolvedContext,
				resolveLocalTypeNames)
		} else {
			if !dectype.IsAny(concrete) {
				//				panic(fmt.Errorf("what is this:%T %v", concrete, concrete))
			}
		}
	case *dectype.AliasReference:
		return ConcreteTypeIfNeeded(t.Next(), concrete, resolveLocalTypeNames)
	case *dectype.Alias:
		return ConcreteTypeIfNeeded(t.Next(), concrete, resolveLocalTypeNames)
	case *dectype.UnmanagedType:
		// UnmanagedTypes can not be concrete
		return nil
	default:
		panic(fmt.Errorf("concrete: what is this %T %s", reference,
			reference.FetchPositionLength().ToCompleteReferenceString()))
	}

	return nil
}

func FillLocalTypesFromSlice(references []dtype.Type, concretes []dtype.Type,
	resolver dectype.DynamicResolver) decshared.DecoratedError {
	if len(concretes) != len(references) {
		return decorated.NewInternalError(
			fmt.Errorf(
				"must have equal number of arguments to concretize slices %v vs %v", concretes, references,
			),
		)
	}

	for index, parameterType := range references {
		resolvedType := parameterType
		argument := concretes[index]
		if dectype.IsLocalType(argument) {
			continue
		}

		var lookupErr decshared.DecoratedError
		lookupErr = ConcreteTypeIfNeeded(parameterType, argument, resolver)
		if lookupErr != nil {
			log.Printf("ERR: %v", lookupErr)
			return lookupErr
		}

		if resolvedType == nil {
			panic(fmt.Errorf("how can resolvedType be nil %T %T %v", parameterType, argument, argument))
		}

		localTypeRef, wasLocalTypeRef := parameterType.(*dectype.LocalTypeNameReference)
		if wasLocalTypeRef {
			var err error

			err = resolver.SetType(localTypeRef.LocalTypeName(), argument)
			//log.Printf("resolved after settype to %T %v", resolvedType, resolvedType)
			if err != nil {
				log.Printf("ERR: %v", err)
				return decorated.NewInternalError(err)
			}
		}

		if resolvedType == nil {
			panic(fmt.Errorf("how can the resolvedType be nil %T %v", argument, argument))
		}
		//log.Printf("resolved to %T", resolvedType)

	}

	return nil
}

func FillContextFromTuple(reference *dectype.TupleTypeAtom, encounteredArguments []dtype.Type,
	resolver dectype.DynamicResolver) decshared.DecoratedError {
	if len(encounteredArguments) != reference.ParameterCount() {
		err := fmt.Errorf("not good, wrong count %v %v %s %s",
			encounteredArguments[0].FetchPositionLength().ToCompleteReferenceString(), reference,
			debug.TreeString(reference), debug.TreeString(encounteredArguments))
		log.Panicf("%v", err)
		return decorated.NewInternalError(err)
	}

	err := FillLocalTypesFromSlice(reference.ParameterTypes(), encounteredArguments, resolver)

	return err
}

func FillContextFromRecordAtom(reference *dectype.RecordAtom, encounteredRecord *dectype.RecordAtom,
	resolver dectype.DynamicResolver) decshared.DecoratedError {
	var encounteredFieldTypes []dtype.Type
	for _, encounteredRecordField := range encounteredRecord.ParseOrderedFields() {
		encounteredFieldTypes = append(encounteredFieldTypes, encounteredRecordField.Type())
	}

	if len(encounteredFieldTypes) != reference.FieldCount() {
		err := fmt.Errorf("not good, wrong field count %v %v %s %s",
			encounteredFieldTypes[0].FetchPositionLength().ToCompleteReferenceString(), reference,
			debug.TreeString(reference), debug.TreeString(encounteredFieldTypes))
		log.Panicf("%v", err)
		return decorated.NewInternalError(err)
	}

	var fieldTypes []dtype.Type
	for _, parseField := range reference.ParseOrderedFields() {
		fieldTypes = append(fieldTypes, parseField.Type())
	}

	err := FillLocalTypesFromSlice(fieldTypes, encounteredFieldTypes, resolver)

	return err
}

func FillContextFromFunction(reference *dectype.FunctionAtom, encounteredArguments []dtype.Type,
	resolver dectype.DynamicResolver) decshared.DecoratedError {

	if hasAnyMatching, startIndex := dectype.HasAnyMatchingTypes(reference.FunctionParameterTypes()); hasAnyMatching {
		originalInitialCount := startIndex
		originalEndCount := len(reference.FunctionParameterTypes()) - startIndex - 2

		originalFirst := append([]dtype.Type{}, reference.FunctionParameterTypes()[0:startIndex]...)

		if len(encounteredArguments) >= len(reference.FunctionParameterTypes()) {
			originalEndCount++
		}

		otherMiddle := encounteredArguments[originalInitialCount : len(encounteredArguments)-originalEndCount]
		if len(otherMiddle) < 1 {
			return decorated.NewInternalError(fmt.Errorf("currently, you must have at least one wildcard parameter"))
		}

		originalEnd := reference.FunctionParameterTypes()[startIndex+1 : len(reference.FunctionParameterTypes())]

		allConverted := append(originalFirst, otherMiddle...)
		allConverted = append(allConverted, originalEnd...)

		//created := dectype.NewFunctionAtom(reference.AstFunction(), allConverted)

		return FillLocalTypesFromSlice(allConverted, encounteredArguments, resolver)
	} else {
		if len(encounteredArguments) != reference.ParameterCount() {
			err := fmt.Errorf("not good, wrong count. expected %d but got %d.\n %s %v %v %s %s",
				reference.ParameterCount(),
				len(encounteredArguments),
				reference.AstFunction().Name(),
				encounteredArguments[0].FetchPositionLength().ToCompleteReferenceString(), reference,
				debug.TreeString(reference), debug.TreeString(encounteredArguments))
			log.Panicf("%v", err)
			return decorated.NewInternalError(err)
		}
		if len(reference.FunctionParameterTypes()) < len(encounteredArguments) {
			return decorated.NewInternalError(fmt.Errorf("too few parameter types"))
		}
	}

	err := FillLocalTypesFromSlice(reference.FunctionParameterTypes(), encounteredArguments, resolver)
	return err
}

func createResolvedFromDynamic(localTypeNameContextRef *dectype.LocalTypeNameOnlyContextReference,
	dynamic *dectype.DynamicLocalTypeResolver) (
	*dectype.ResolvedLocalTypeContext, decshared.DecoratedError,
) {
	if dynamic.IsDefined() {
		resolved, resolvedErr := dectype.NewResolvedLocalTypeContext(localTypeNameContextRef, dynamic.ArgumentTypes())
		if resolvedErr != nil {
			return nil, decorated.NewInternalError(resolvedErr)
		}
		return resolved, nil
	}

	err := fmt.Errorf(
		"dynamic was not filled in %s\n%s",
		debug.TreeString(localTypeNameContextRef), dynamic.DebugAllNotDefined())
	log.Println(err)

	return nil, decorated.NewInternalError(err)
}

func handleFunction(functionAtom *dectype.FunctionAtom,
	localTypeNameContextRef *dectype.LocalTypeNameOnlyContextReference, concreteArguments []dtype.Type) (
	*dectype.ResolvedLocalTypeContext, decshared.DecoratedError,
) {
	resolver := dectype.NewDynamicLocalTypeResolver(localTypeNameContextRef.LocalTypeNameContext().Names())
	err := FillContextFromFunction(functionAtom, concreteArguments, resolver)
	if err != nil {
		return nil, err
	}

	context, resolveErr := createResolvedFromDynamic(localTypeNameContextRef, resolver)
	if resolveErr != nil {
		log.Printf("problem resolving %s", concreteArguments[0].FetchPositionLength().ToCompleteReferenceString())
	}
	return context, resolveErr
}

func ConcretizeLocalTypeContextUsingArguments(localTypeNameContext *dectype.LocalTypeNameOnlyContextReference,
	concreteArguments []dtype.Type) (
	*dectype.ResolvedLocalTypeContext, decshared.DecoratedError,
) {
	if localTypeNameContext == nil {
		panic(fmt.Errorf("localTypeNameContext is nil"))
	}
	switch t := localTypeNameContext.Next().Next().(type) {
	case *dectype.FunctionAtom:
		return handleFunction(t, localTypeNameContext, concreteArguments)
	case *dectype.FunctionTypeReference:
		return handleFunction(t.FunctionAtom(), localTypeNameContext, concreteArguments)
	case *dectype.CustomTypeAtom:
	case *dectype.RecordAtom:
	case *dectype.PrimitiveAtom:

	default:
		return nil, decorated.NewInternalError(fmt.Errorf("not sure what this is %T", t))
	}
	resolvedContext, err := dectype.NewResolvedLocalTypeContext(localTypeNameContext, concreteArguments)
	if err != nil {
		return nil, decorated.NewInternalError(err)
	}

	return resolvedContext, nil
}
