/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"github.com/swamp/compiler/src/token"
)

type Type interface {
	Name() string
	String() string
	FetchPositionLength() token.SourceFileReference
}

func NewCustomType(keywordType token.Keyword, customTypeName *TypeIdentifier, variants []*CustomTypeVariant, typeParameterIdentifiers []*TypeParameter) *CustomType {
	c := &CustomType{keywordType: keywordType, name: customTypeName, variants: variants, typeParameters: typeParameterIdentifiers}
	for _, variant := range variants {
		variant.SetParent(c)
	}
	return c
}

func NewLocalType(typeArgument *TypeParameter) *LocalType {
	return &LocalType{typeParameterReference: typeArgument}
}

func NewTypeReference(ident *TypeIdentifier, arguments []Type) *TypeReference {
	return &TypeReference{ident: ident, arguments: arguments}
}

func NewScopedTypeReference(ident *TypeIdentifierScoped, arguments []Type) *TypeReferenceScoped {
	return &TypeReferenceScoped{ident: ident, arguments: arguments}
}
