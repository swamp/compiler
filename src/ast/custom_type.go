/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type CustomType struct {
	name              *TypeIdentifier
	typeParameters    []*TypeParameter
	variants          []*CustomTypeVariant
	keywordType       token.Keyword
	precedingComments *MultilineComment
}

func (i *CustomType) String() string {
	return fmt.Sprintf("[custom-type %v %v]", i.name, i.variants)
}

func (i *CustomType) Identifier() *TypeIdentifier {
	return i.name
}

func (i *CustomType) Name() string {
	return i.name.Name()
}

func (i *CustomType) Variants() []*CustomTypeVariant {
	return i.variants
}

func (i *CustomType) FindAllLocalTypes() []*TypeParameter {
	return i.typeParameters
}

func (i *CustomType) FetchPositionLength() token.SourceFileReference {
	return i.name.FetchPositionLength()
}

func (i *CustomType) KeywordType() token.Keyword {
	return i.keywordType
}

func (i *CustomType) DebugString() string {
	return fmt.Sprintf("[custom-type-statement %v]", i.name)
}

func NewCustomType(keywordType token.Keyword, customTypeName *TypeIdentifier, variants []*CustomTypeVariant, typeParameterIdentifiers []*TypeParameter) *CustomType {
	c := &CustomType{keywordType: keywordType, name: customTypeName, variants: variants, typeParameters: typeParameterIdentifiers}
	for _, variant := range variants {
		variant.SetParent(c)
	}
	return c
}
