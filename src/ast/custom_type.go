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
	name           *TypeIdentifier
	typeParameters []*TypeParameter
	variants       []*CustomTypeVariant
	keywordType    token.Keyword
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
