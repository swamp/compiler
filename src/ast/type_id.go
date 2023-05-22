/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"
	"reflect"

	"github.com/swamp/compiler/src/token"
)

type TypeId struct {
	typeRef     TypeIdentifierNormalOrScoped
	typeIdToken token.TypeId `debug:"true"`
	inclusive   token.SourceFileReference
}

func NewTypeId(typeIdToken token.TypeId, typeIdentifier TypeIdentifierNormalOrScoped) *TypeId {
	if typeIdentifier == nil || reflect.ValueOf(typeIdentifier).IsNil() {
		panic("type identifier can not be nil")
	}
	inclusive := token.MakeInclusiveSourceFileReference(typeIdToken.SourceFileReference, typeIdentifier.FetchPositionLength())
	return &TypeId{typeRef: typeIdentifier, typeIdToken: typeIdToken, inclusive: inclusive}
}

func (i *TypeId) TypeIdentifier() TypeIdentifierNormalOrScoped {
	return i.typeRef
}

func (i *TypeId) String() string {
	return fmt.Sprintf("[TypeId %v]", i.typeRef)
}

func (i *TypeId) Name() string {
	return fmt.Sprintf("id<%v>", i.typeRef.Name())
}

func (i *TypeId) DebugString() string {
	return i.Name()
}

func (i *TypeId) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}
