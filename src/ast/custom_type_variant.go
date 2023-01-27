/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type CustomTypeVariant struct {
	typeIdentifier *TypeIdentifier
	typeParameters *LocalTypeNameDefinitionContext
	userTypes      []Type
	parent         *CustomType
	index          int
	comment        token.Comment
}

func NewCustomTypeVariant(index int, typeIdentifier *TypeIdentifier, userTypes []Type, comment token.Comment) *CustomTypeVariant {
	return &CustomTypeVariant{index: index, typeIdentifier: typeIdentifier, userTypes: userTypes, comment: comment}
}

func (i *CustomTypeVariant) TypeIdentifier() *TypeIdentifier {
	return i.typeIdentifier
}

func (i *CustomTypeVariant) Comment() token.Comment {
	return i.comment
}

func (i *CustomTypeVariant) TypeParameterContext() *LocalTypeNameDefinitionContext {
	return i.typeParameters
}

func (i *CustomTypeVariant) Name() string {
	return i.typeIdentifier.Name()
}

func (i *CustomTypeVariant) Types() []Type {
	return i.userTypes
}

func (i *CustomTypeVariant) Parent() *CustomType {
	return i.parent
}

func (i *CustomTypeVariant) SetParent(parent *CustomType) {
	i.parent = parent
}

func (i *CustomTypeVariant) Index() int {
	return i.index
}

func (i *CustomTypeVariant) String() string {
	return fmt.Sprintf("[Variant %v %v]", i.typeIdentifier, i.userTypes)
}

func (i *CustomTypeVariant) FetchPositionLength() token.SourceFileReference {
	return i.typeIdentifier.FetchPositionLength()
}
