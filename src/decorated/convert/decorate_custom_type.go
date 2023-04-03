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

func DecorateCustomType(astCustomType *ast.CustomType,
	typeRepo decorated.TypeAddAndReferenceMaker) (*dectype.CustomTypeAtom, decshared.DecoratedError) {
	var variants []*dectype.CustomTypeVariantAtom

	artifactTypeName := typeRepo.SourceModule().FullyQualifiedModuleName().JoinTypeIdentifier(astCustomType.Identifier())

	var astArgumentTypes []dtype.Type

	for _, astArgumentType := range astCustomType.Arguments() {
		newType, newTypeErr := ConvertFromAstToDecorated(astArgumentType, typeRepo)
		if newTypeErr != nil {
			return nil, decorated.NewInternalError(newTypeErr)
		}
		astArgumentTypes = append(astArgumentTypes, newType)
	}

	s := dectype.NewCustomTypePrepare(astCustomType, astArgumentTypes, artifactTypeName)

	for astVariantIndex, astVariant := range astCustomType.Variants() {
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
