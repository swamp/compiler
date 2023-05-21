/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
	"log"
)

type ResolvedLocalType struct {
	identifier     *LocalTypeName
	referencedType dtype.Type
	wasReferenced  bool
	references     []*ResolvedLocalTypeReference
}

func (u *ResolvedLocalType) String() string {
	return fmt.Sprintf("%v:%v", u.identifier.Name(), u.referencedType)
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
	return u.identifier
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

func (u *ResolvedLocalType) ParameterCount() int {
	return 0
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

	//	TypeChain(referencedType, 0)
}

/*

func (u *ResolvedLocalType) SetDefinition(referencedType dtype.Type) error {
	u.Verify(referencedType)
	u.referencedType = referencedType
	u.hasBeenDefined = true
	return nil
}
*/

func NewResolvedLocalType(identifier *LocalTypeName, referencedType dtype.Type) *ResolvedLocalType {
	if !identifier.identifier.LocalType().FetchPositionLength().Verify() {
		panic(fmt.Errorf("wrong identifier"))
	}
	if identifier.identifier.LocalType().FetchPositionLength().Range.Position().Line() == 1 && identifier.identifier.LocalType().FetchPositionLength().Range.Position().Column() == 0 {
		log.Printf("found")
	}
	log.Printf("NewResolvedLocalType %T %v %v", referencedType, referencedType.FetchPositionLength().ToCompleteReferenceString(), referencedType)
	x := &ResolvedLocalType{identifier: identifier, referencedType: referencedType}
	return x
}
