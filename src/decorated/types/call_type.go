/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"github.com/swamp/compiler/src/decorated/dtype"
)

type Lookup interface {
	LookupType(name string) (dtype.Type, error)
}

/*
func replaceCustomTypeFromContext(customType *CustomTypeAtom, lookup ParseReferenceFromName) (*CustomTypeAtom, error) {
	var replacedVariants []*CustomTypeVariantAtom

	var replacedGenerics []dtype.Type
	for _, genericType := range customType.Parameters() {
		lookedUpType, err := lookup.LookupType(genericType.identifier.Name())
		if err != nil {
			return nil, err
		}
		replacedGenerics = append(replacedGenerics, genericType)
	}

	newCustomType := NewCustomTypePrepare(customType.astCustomType, customType.artifactTypeName, replacedGenerics)
	for _, field := range customType.Variants() {
		var variantParameters []dtype.Type
		for _, param := range field.parameterFields {
			converted, convertedErr := ReplaceTypeFromContext(param.parameterType, lookup)
			if convertedErr != nil {
				return nil, convertedErr
			}
			variantParameters = append(variantParameters, converted)
		}

		newField, newErr := replaceCustomTypeVariantFromContext(newCustomType, field, lookup)
		if newErr != nil {
			return nil, newErr
		}

		replacedVariants = append(replacedVariants, newField)
	}

	newCustomType.FinalizeVariants(replacedVariants)

	return newCustomType, nil
}
*/
