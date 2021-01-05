/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func finalizeCustomType(customType *dectype.CustomTypeAtom) {
	for _, variant := range customType.Variants() {
		variant.AttachToCustomType(customType)
	}
}

func decorateCustomTypeVariantConstructors(customType *dectype.CustomTypeAtom, typeRepo *dectype.TypeRepo) {
	for _, variant := range customType.Variants() {
		constructorType := dectype.NewCustomTypeVariantConstructorType(typeRepo, variant)
		typeRepo.DeclareType(constructorType)
	}
}

func DecorateCustomType(customTypeDefinition *ast.CustomType,
	typeRepo *dectype.TypeRepo) (*dectype.CustomTypeAtom, decshared.DecoratedError) {
	var variants []*dectype.CustomTypeVariant

	for astVariantIndex, astVariant := range customTypeDefinition.Variants() {
		var astVariantTypes []dtype.Type

		for _, astVariantType := range astVariant.Types() {
			newType, newTypeErr := ConvertFromAstToDecorated(astVariantType, typeRepo)
			if newTypeErr != nil {
				return nil, decorated.NewUnknownTypeInCustomTypeVariant(astVariant, newTypeErr)
			}
			astVariantTypes = append(astVariantTypes, newType.(dtype.Type))
		}

		variant := dectype.NewCustomTypeVariant(astVariantIndex, astVariant.TypeIdentifier(), astVariantTypes)

		variants = append(variants, variant)
	}

	var decoratedTypeParameters []dtype.Type

	for _, typeParameter := range customTypeDefinition.FindAllLocalTypes() {
		converted, convertErr := ConvertFromAstToDecorated(typeParameter, typeRepo)
		if convertErr != nil {
			return nil, decorated.NewInternalError(convertErr)
		}

		decoratedTypeParameters = append(decoratedTypeParameters, converted)
	}

	decoratedNames := AstParametersToArgumentNames(customTypeDefinition.FindAllLocalTypes())
	artifactTypeName := typeRepo.ModuleName().JoinTypeIdentifier(customTypeDefinition.Identifier())
	s := dectype.NewCustomType(customTypeDefinition.Identifier(), artifactTypeName, decoratedNames, variants)
	finalizeCustomType(s)
	decorateCustomTypeVariantConstructors(s, typeRepo)

	return s, nil
}
