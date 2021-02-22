/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type TypeId struct {
	typeRef     Type
	typeIdToken token.TypeId
}

func NewTypeId(typeIdToken token.TypeId, typeRef Type) *TypeId {
	return &TypeId{typeRef: typeRef, typeIdToken: typeIdToken}
}

func (i *TypeId) TypeRef() Type {
	return i.typeRef
}

func (i *TypeId) String() string {
	return fmt.Sprintf("[type-id %v]", i.typeRef)
}

func (i *TypeId) Name() string {
	return fmt.Sprintf("id<%v>", i.typeRef.Name())
}

func (i *TypeId) DebugString() string {
	return i.Name()
}

func (i *TypeId) FetchPositionLength() token.Range {
	return i.typeIdToken.FetchPositionLength()
}
