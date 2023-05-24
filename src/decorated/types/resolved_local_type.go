/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type ResolvedLocalType struct {
	debugLocalTypeName *LocalTypeName `debug:"true"`
	referencedType     dtype.Type     `debug:"true"`
	wasReferenced      bool
	references         []*ResolvedLocalTypeReference
}

func (u *ResolvedLocalType) String() string {
	return fmt.Sprintf("%v:%v", u.debugLocalTypeName.Identifier().Name(), u.referencedType)
}

func (u *ResolvedLocalType) AddReference(ref *ResolvedLocalTypeReference) {
	if ref == u.referencedType {
		panic("problem")
	}
	u.references = append(u.references, ref)
}

func (u *ResolvedLocalType) FetchPositionLength() token.SourceFileReference {
	return u.referencedType.FetchPositionLength()
}

func (u *ResolvedLocalType) HumanReadable() string {
	return fmt.Sprintf("%v", u.referencedType.HumanReadable())
}

func (u *ResolvedLocalType) Identifier() *LocalTypeName {
	return u.debugLocalTypeName
}

func (u *ResolvedLocalType) WantsToBeReplaced() bool {
	return IsAny(u.referencedType)
}

func (u *ResolvedLocalType) IsEqual(_ dtype.Atom) error {
	return nil
}

func (u *ResolvedLocalType) Resolve() (dtype.Atom, error) {
	return u.referencedType.Resolve()
}

func (u *ResolvedLocalType) ReferencedType() dtype.Type {
	return u.referencedType
}

func (u *ResolvedLocalType) Next() dtype.Type {
	return u.referencedType
}

func (u *ResolvedLocalType) WasReferenced() bool {
	return u.wasReferenced
}

func (u *ResolvedLocalType) MarkAsReferenced() {
	u.wasReferenced = true
}

func (u *ResolvedLocalType) Verify(referencedType dtype.Type) {
	if referencedType == nil {
		return
	}

	ref, wasRef := referencedType.(*ResolvedLocalTypeReference)
	if wasRef {
		if ref.typeDefinition == u {
			panic("circular")
		}
	}

	//	TypeChain(localTypeName, 0)
}

/*

func (u *ResolvedLocalType) SetDefinition(localTypeName dtype.Type) error {
	u.Verify(localTypeName)
	u.localTypeName = localTypeName
	u.hasBeenDefined = true
	return nil
}
*/

func NewResolvedLocalType(localTypeName *LocalTypeName,
	referencedType dtype.Type) *ResolvedLocalType {
	localIdent := localTypeName.identifier
	if !localIdent.LocalType().FetchPositionLength().Verify() {
		panic(fmt.Errorf("wrong localTypeName"))
	}
	x := &ResolvedLocalType{debugLocalTypeName: localTypeName, referencedType: referencedType}
	return x
}
