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
	variants          []*CustomTypeVariant
	keywordType       token.Keyword
	precedingComments *MultilineComment
	parameterTypes    []Type
}

func (i *CustomType) Arguments() []Type {
	return i.parameterTypes
}

func (i *CustomType) String() string {
	return fmt.Sprintf("[CustomType %v<%v> %v]", i.name, i.parameterTypes, i.variants)
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

func (i *CustomType) FetchPositionLength() token.SourceFileReference {
	return i.name.FetchPositionLength()
}

func (i *CustomType) KeywordType() token.Keyword {
	return i.keywordType
}

func (i *CustomType) Comment() *MultilineComment {
	return i.precedingComments
}

func (i *CustomType) DebugString() string {
	return fmt.Sprintf("[CustomType %v]", i.name)
}

func NewCustomType(keywordType token.Keyword, customTypeName *TypeIdentifier, parameterTypes []Type, variants []*CustomTypeVariant, comment *MultilineComment) *CustomType {
	c := &CustomType{keywordType: keywordType, name: customTypeName, parameterTypes: parameterTypes, variants: variants, precedingComments: comment}
	for _, variant := range variants {
		variant.SetParent(c)
	}
	return c
}
