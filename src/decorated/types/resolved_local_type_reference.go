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

type ResolvedLocalTypeReference struct {
	identifier     *LocalTypeName
	typeDefinition *ResolvedLocalType `debug:"true"`
}

func (u *ResolvedLocalTypeReference) String() string {
	return fmt.Sprintf("[ConcreteGenericRef %v => %v]", u.typeDefinition.debugLocalTypeName,
		u.typeDefinition.referencedType)
}

func (u *ResolvedLocalTypeReference) FetchPositionLength() token.SourceFileReference {
	return u.identifier.FetchPositionLength()
}

func (u *ResolvedLocalTypeReference) HumanReadable() string {
	return fmt.Sprintf("%v", u.typeDefinition.HumanReadable())
}

func (u *ResolvedLocalTypeReference) TypeDefinition() *ResolvedLocalType {
	return u.typeDefinition
}

func (u *ResolvedLocalTypeReference) Identifier() *LocalTypeName {
	return u.identifier
}

func (u *ResolvedLocalTypeReference) IsEqual(_ dtype.Atom) error {
	return nil
}

func (u *ResolvedLocalTypeReference) Resolve() (dtype.Atom, error) {
	return u.typeDefinition.Resolve()
}

func (u *ResolvedLocalTypeReference) ReferencedType() dtype.Type {
	return u.typeDefinition
}

func (u *ResolvedLocalTypeReference) Next() dtype.Type {
	return u.typeDefinition
}

func (u *ResolvedLocalTypeReference) WasReferenced() bool {
	return false
}

func NewResolvedLocalTypeReference(identifier *LocalTypeName,
	referencedDefinition *ResolvedLocalType) *ResolvedLocalTypeReference {
	x := &ResolvedLocalTypeReference{identifier: identifier, typeDefinition: referencedDefinition}
	referencedDefinition.AddReference(x)

	return x
}
