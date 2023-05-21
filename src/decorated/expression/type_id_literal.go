/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	dectype "github.com/swamp/compiler/src/decorated/types"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type TypeIdLiteral struct {
	typeId                *ast.TypeId
	constructedTypeIdType dtype.Type
	containedType         dectype.TypeReferenceScopedOrNormal
}

func NewTypeIdLiteral(typeId *ast.TypeId, constructedTypeIdType dtype.Type, containedType dectype.TypeReferenceScopedOrNormal) *TypeIdLiteral {
	return &TypeIdLiteral{typeId: typeId, constructedTypeIdType: constructedTypeIdType, containedType: containedType}
}

func (i *TypeIdLiteral) Type() dtype.Type {
	return i.constructedTypeIdType
}

func (i *TypeIdLiteral) ContainedType() dectype.TypeReferenceScopedOrNormal {
	return i.containedType
}

func (i *TypeIdLiteral) String() string {
	return fmt.Sprintf("[TypeIdLiteral %v]", i.constructedTypeIdType)
}

func (i *TypeIdLiteral) HumanReadable() string {
	return "Type ID Literal"
}

func (i *TypeIdLiteral) FetchPositionLength() token.SourceFileReference {
	return i.typeId.FetchPositionLength()
}
