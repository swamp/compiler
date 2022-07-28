/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func AstParametersToArgumentNames(typeParameters []*ast.TypeParameter) []*dtype.TypeArgumentName {
	var argumentNames []*dtype.TypeArgumentName
	for _, param := range typeParameters {
		argumentNames = append(argumentNames, dtype.NewTypeArgumentName(param.Identifier()))
	}
	return argumentNames
}

func AstParametersToLocalTypes(typeParameters []*ast.TypeParameter) []dtype.Type {
	var argumentNames []dtype.Type
	for _, param := range typeParameters {
		astLocalType := ast.NewLocalType(param)
		argumentNames = append(argumentNames, dectype.NewLocalType(astLocalType.TypeParameter()))
	}
	return argumentNames
}

func newInvokerType(astTypeReference ast.TypeReferenceScopedOrNormal, foundType dectype.TypeReferenceScopedOrNormal, t decorated.TypeAddAndReferenceMaker) (dtype.Type, decshared.DecoratedError) {
	unaliasedTypeToCheck := dectype.Unalias(foundType)

	if unaliasedTypeToCheck.ParameterCount() != len(astTypeReference.Arguments()) {
		return nil, decorated.NewInternalError(fmt.Errorf("problems number of arguments %v (\n\n%v\n vs\n\n%v\n) (found %T, expected %T) found %v vs expected %v (%v)", astTypeReference, foundType, astTypeReference, foundType, astTypeReference, foundType.ParameterCount(), len(astTypeReference.Arguments()), astTypeReference.Arguments()))
	}

	if foundType.ParameterCount() == 0 {
		return foundType, nil
	}

	var convertedParameters []dtype.Type

	for _, a := range astTypeReference.Arguments() {
		convertedParameter, convertedParameterErr := ConvertFromAstToDecorated(a, t)
		if convertedParameterErr != nil {
			return nil, convertedParameterErr
		}

		convertedParameters = append(convertedParameters, convertedParameter)
	}

	return dectype.NewInvokerType(foundType, convertedParameters)
}

func DecorateTupleType(tupleType *ast.TupleType, t decorated.TypeAddAndReferenceMaker) (dtype.Type, decshared.DecoratedError) {
	var convertedParameters []*dectype.TupleTypeField
	for index, a := range tupleType.Types() {
		convertedParameter, convertedParameterErr := ConvertFromAstToDecorated(a, t)
		if convertedParameterErr != nil {
			return nil, convertedParameterErr
		}
		field := dectype.NewTupleTypeField(index, convertedParameter)
		convertedParameters = append(convertedParameters, field)
	}

	return dectype.NewTupleTypeAtom(tupleType, convertedParameters), nil
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
		return dectype.NewFunctionTypeReference(info, functionType, nil), nil

	case *ast.Alias:
		subType, subTypeErr := ConvertFromAstToDecorated(info.ReferencedType(), t)
		if subTypeErr != nil {
			return nil, subTypeErr
		}
		artifactTypeName := t.SourceModule().FullyQualifiedModuleName().JoinTypeIdentifier(info.Identifier())
		newType := dectype.NewAliasType(info, artifactTypeName, subType)
		return newType, t.AddTypeAlias(newType)

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

	case *ast.LocalType:
		return dectype.NewLocalType(info.TypeParameter()), nil
	case *ast.TypeReferenceScoped:
		foundType, err := t.CreateSomeTypeReference(info.TypeResolver())
		if err != nil {
			return nil, err
		}
		if foundType == nil {
			return nil, decorated.NewInternalError(fmt.Errorf("coulfdn't find type reference %v", info))
		}
		return newInvokerType(info, foundType, t)
	case *ast.TypeReference:
		refName := info.TypeIdentifier()
		foundType, err := t.CreateSomeTypeReference(refName)
		if err != nil {
			return nil, err
		}
		if foundType == nil {
			return nil, decorated.NewInternalError(fmt.Errorf("coulfdn't find type reference %v", refName))
		}

		return newInvokerType(info, foundType, t)
	case *ast.TypeParameter:
		return dectype.NewLocalType(info), nil
	case *ast.AnyMatchingType:
		return dectype.NewAnyMatchingTypes(info), nil
	case *ast.UnmanagedType:
		return dectype.NewUnmanagedType(info), nil
	default:
		return nil, decorated.NewInternalError(fmt.Errorf("can't convert this ast type %v %T", astType, astType))
	}
}
