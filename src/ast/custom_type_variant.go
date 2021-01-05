/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"
)

type CustomTypeVariant struct {
	typeIdentifier *TypeIdentifier
	userTypes      []Type
	parent *CustomType
	index int
}

func NewCustomTypeVariant(index int, typeIdentifier *TypeIdentifier, userTypes []Type) *CustomTypeVariant {
	return &CustomTypeVariant{index: index, typeIdentifier: typeIdentifier, userTypes: userTypes}
}

func (i *CustomTypeVariant) TypeIdentifier() *TypeIdentifier {
	return i.typeIdentifier
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
	return fmt.Sprintf("[variant %v%v]", i.typeIdentifier, i.userTypes)
}
