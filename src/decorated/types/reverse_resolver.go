/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"github.com/swamp/compiler/src/decorated/dtype"
)

type DynamicReverseResolver struct {
	reverseLookup map[string]*LocalTypeNameReference
	nextResolver  DynamicResolver
}

func NewDynamicReverseResolver(reverseLookup map[string]*LocalTypeNameReference,
	nextResolver DynamicResolver) *DynamicReverseResolver {
	return &DynamicReverseResolver{reverseLookup: reverseLookup, nextResolver: nextResolver}
}

func (t *DynamicReverseResolver) SetType(defRef *LocalTypeName, definedType dtype.Type) error {
	reverse, foundReverse := t.reverseLookup[defRef.Name()]
	nameToUse := defRef
	if foundReverse {
		nameToUse = reverse.localTypeName
	}
	return t.nextResolver.SetType(nameToUse, definedType)
}
