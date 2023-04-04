/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/concretize"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"log"
)

func AstLocalTypeNamesToTypeArgumentName(typeParameters []*ast.LocalTypeName) []*dtype.LocalTypeName {
	var argumentNames []*dtype.LocalTypeName
	for _, param := range typeParameters {
		argumentNames = append(argumentNames, dtype.NewLocalTypeName(param))
	}
	return argumentNames
}

func DecorateTupleType(tupleType *ast.TupleType, t decorated.TypeAddAndReferenceMaker) (dtype.Type, decshared.DecoratedError) {
	var convertedParameters []dtype.Type
	for _, a := range tupleType.Types() {
		convertedParameter, convertedParameterErr := ConvertFromAstToDecorated(a, t)
		if convertedParameterErr != nil {
			return nil, convertedParameterErr
		}
		convertedParameters = append(convertedParameters, convertedParameter)
	}

	return dectype.NewTupleTypeAtom(tupleType, convertedParameters), nil
}

func ConvertFromAstToDecoratedSlice(astTypes []ast.Type, t decorated.TypeAddAndReferenceMaker) ([]dtype.Type, decshared.DecoratedError) {
	var types []dtype.Type
	for _, astType := range astTypes {
		convertedParameter, convertedParameterErr := ConvertFromAstToDecorated(astType, t)
		if convertedParameterErr != nil {
			return nil, convertedParameterErr
		}
		types = append(types, convertedParameter)
	}

	return types, nil
}

func ConvertFromAstToDecorated(astType ast.Type,
	t decorated.TypeAddAndReferenceMaker) (dtype.Type, decshared.DecoratedError) {
	switch info := astType.(type) {
	case *ast.FunctionType:
		var convertedParameters []dtype.Type
		for _, a := range info.FunctionParameters() {
			convertedParameter, convertedParameterErr := ConvertFromAstToDecorated(a, t)
			if convertedParameterErr != nil {
				return nil, convertedParameterErr
			}
			convertedParameters = append(convertedParameters, convertedParameter)
		}
		functionType := dectype.NewFunctionAtom(info, convertedParameters)
		return functionType, nil

	case *ast.Alias:
		subType, subTypeErr := ConvertFromAstToDecorated(info.ReferencedType(), t)
		if subTypeErr != nil {
			return nil, subTypeErr
		}
		artifactTypeName := t.SourceModule().FullyQualifiedModuleName().JoinTypeIdentifier(info.Identifier())
		newType := dectype.NewAliasType(info, artifactTypeName, subType)
		return newType, t.AddTypeAlias(newType)

	case *ast.LocalTypeNameDefinitionContext:
		decContext := dectype.NewLocalTypeNameContext()
		for _, name := range info.LocalTypeNames() {
			decName := dtype.NewLocalTypeName(name)
			decContext.AddDef(decName)
		}
		subContext := t.MakeLocalNameContext(decContext)
		subType, subTypeErr := ConvertFromAstToDecorated(info.Next(), subContext)
		if subTypeErr != nil {
			return nil, subTypeErr
		}
		decContext.SetType(subType)

		return decContext, nil
	case *ast.Record:
		return DecorateRecordType(info, t)

	case *ast.TupleType:
		return DecorateTupleType(info, t)

	case *ast.TypeIdentifier:
		refName := info.Symbol().Name()
		foundType, err := t.CreateSomeTypeReference(info)
		if err != nil {
			return nil, err
		}
		if foundType == nil {
			return nil, decorated.NewInternalError(fmt.Errorf("couldn't find type identifier %v", refName))
		}
		return foundType, nil

	case *ast.LocalTypeNameReference:
		//artifactTypeName := t.SourceModule().FullyQualifiedModuleName().JoinTypeIdentifier(info.Identifier())
		return t.CreateLocalTypeNameOnlyReference(info)
		//return dectype.NewLocalTypeDefinitionReference(info, dectype.NewAnyType()), nil

	case *ast.TypeReferenceScoped:
		foundType, err := t.CreateSomeTypeReference(info.TypeResolver())
		if err != nil {
			return nil, err
		}
		if foundType == nil {
			return nil, decorated.NewInternalError(fmt.Errorf("coulfdn't find type reference %v", info))
		}
		return foundType, nil
	case *ast.TypeReference:
		refName := info.TypeIdentifier()
		foundType, err := t.CreateSomeTypeReference(refName)
		if err != nil {
			return nil, err
		}
		if foundType == nil {
			return nil, decorated.NewInternalError(fmt.Errorf("couldn't find type reference %v", refName))
		}
		types, sliceErr := ConvertFromAstToDecoratedSlice(info.Arguments(), t)
		if sliceErr != nil {
			return nil, sliceErr
		}

		nameOnlyContext := dectype.FindNameOnlyContextWithUnalias(foundType)
		if nameOnlyContext != nil {
			newType, concreteErr := concretize.ConcreteArguments(nameOnlyContext, types)
			if concreteErr != nil {
				return nil, concreteErr
			}
			switch newTypeT := newType.(type) {
			case *dectype.PrimitiveAtom:
				foundType = dectype.NewPrimitiveTypeReference(foundType.NameReference(), newTypeT)
				break
			case *dectype.CustomTypeAtom:
				foundType = dectype.NewCustomTypeReference(foundType.NameReference(), newTypeT)
			case *dectype.RecordAtom:
			default:
				log.Printf("nameOnly what is this: %T", newType)
			}
		}
		return foundType, nil
	//case *ast.LocalTypeNameDefinition:
	//	return dectype.NewLocalTypeNameDefinition(info, dectype.NewAnyType()), nil
	case *ast.AnyMatchingType:
		return dectype.NewAnyMatchingTypes(info), nil
	case *ast.UnmanagedType:
		return dectype.NewUnmanagedType(info), nil
	default:
		return nil, decorated.NewInternalError(fmt.Errorf("can't convert this ast type %v %T", astType, astType))
	}
}
