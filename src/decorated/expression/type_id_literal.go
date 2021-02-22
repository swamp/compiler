/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type TypeIdLiteral struct {
	typeId                *ast.TypeId
	constructedTypeIdType dtype.Type
	containedType         dtype.Type
}

func NewTypeIdLiteral(typeId *ast.TypeId, constructedTypeIdType dtype.Type, containedType dtype.Type) *TypeIdLiteral {
	return &TypeIdLiteral{typeId: typeId, constructedTypeIdType: constructedTypeIdType, containedType: containedType}
}

func (i *TypeIdLiteral) Type() dtype.Type {
	return i.constructedTypeIdType
}

func (i *TypeIdLiteral) ContainedType() dtype.Type {
	return i.containedType
}

func (i *TypeIdLiteral) String() string {
	return fmt.Sprintf("[typeid %v]", i.typeId)
}

func (i *TypeIdLiteral) FetchPositionLength() token.Range {
	return i.typeId.FetchPositionLength()
}
