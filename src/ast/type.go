/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

type Type interface {
	Name() string
	String() string
}



func NewCustomType(customTypeName *TypeIdentifier, variants []*CustomTypeVariant, typeParameterIdentifiers []*TypeParameter) *CustomType {
	c := &CustomType{name: customTypeName, variants: variants, typeParameters: typeParameterIdentifiers}
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

