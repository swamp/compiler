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

func DecorateCustomType(customTypeDefinition *ast.CustomType,
	typeRepo decorated.TypeAddAndReferenceMaker) (*dectype.CustomTypeAtom, decshared.DecoratedError) {
	var variants []*dectype.CustomTypeVariantAtom

	genericLocalTypes := AstParametersToLocalTypes(customTypeDefinition.FindAllLocalTypes())
	artifactTypeName := typeRepo.SourceModule().FullyQualifiedModuleName().JoinTypeIdentifier(customTypeDefinition.Identifier())

	s := dectype.NewCustomTypePrepare(customTypeDefinition, artifactTypeName, genericLocalTypes)

	for astVariantIndex, astVariant := range customTypeDefinition.Variants() {
		var astVariantTypes []dtype.Type

		for _, astVariantType := range astVariant.Types() {
			newType, newTypeErr := ConvertFromAstToDecorated(astVariantType, typeRepo)
			if newTypeErr != nil {
				return nil, decorated.NewUnknownTypeInCustomTypeVariant(astVariant, newTypeErr)
			}
			astVariantTypes = append(astVariantTypes, newType)
		}

		variant := dectype.NewCustomTypeVariant(astVariantIndex, s, astVariant, astVariantTypes)

		variants = append(variants, variant)
	}

	s.FinalizeVariants(variants)

	return s, nil
}
