/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
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

func ConvertFromAstToDecorated(astType ast.Type,
	t *dectype.TypeRepo) (dtype.Type, dectype.DecoratedTypeError) {
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
		functionType := t.AddFunctionAtom(info, convertedParameters)
		return dectype.NewFunctionTypeReference(info, functionType), nil

	case *ast.Alias:
		subType, subTypeErr := ConvertFromAstToDecorated(info.ReferencedType(), t)
		if subTypeErr != nil {
			return nil, subTypeErr
		}
		return t.DeclareTypeAlias(info, subType)

	case *ast.Record:
		return DecorateRecordType(info, t)

	case *ast.TypeIdentifier:
		refName := info.Symbol().Name()
		foundType := t.CreateTypeReference(info)
		if foundType == nil {
			return nil, fmt.Errorf("couldn't find type identifier %v", refName)
		}
		return foundType, nil

	case *ast.LocalType:
		return dectype.NewLocalType(info.TypeParameter()), nil

	case *ast.TypeReference:
		refName := info.TypeResolver()
		foundType := t.CreateTypeReference(refName)
		if foundType == nil {
			return nil, fmt.Errorf("couldn't find type reference %v", refName)
		}

		unaliasedTypeToCheck := dectype.Unalias(foundType)

		if unaliasedTypeToCheck.ParameterCount() != len(info.Arguments()) {
			return nil, fmt.Errorf("problems number of arguments %v (\n\n%v\n vs\n\n%v\n) (found %T, expected %T) found %v vs expected %v (%v)", refName, foundType, info, foundType, info, foundType.ParameterCount(), len(info.Arguments()), info.Arguments())
		}

		if foundType.ParameterCount() == 0 {
			return foundType, nil
		}

		var convertedParameters []dtype.Type

		for _, a := range info.Arguments() {
			convertedParameter, convertedParameterErr := ConvertFromAstToDecorated(a, t)
			if convertedParameterErr != nil {
				return nil, convertedParameterErr
			}

			convertedParameters = append(convertedParameters, convertedParameter)
		}

		return dectype.NewInvokerType(foundType, convertedParameters)
	case *ast.TypeParameter:
		return dectype.NewLocalType(info), nil
	default:
		return nil, decorated.NewInternalError(fmt.Errorf("xcan not convert this ast type %v %T", astType, astType))
	}
}
