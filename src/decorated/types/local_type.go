/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type LocalTypeDefinition struct {
	identifier     *dtype.LocalTypeName
	referencedType dtype.Type
	hasBeenDefined bool
	wasReferenced  bool
	references     []*LocalTypeDefinitionReference
}

func (u *LocalTypeDefinition) String() string {
	return fmt.Sprintf("[ConcreteGeneric %v %v]", u.identifier.Name(), u.referencedType)
}

func (u *LocalTypeDefinition) AddReference(ref *LocalTypeDefinitionReference) {
	if ref == u.referencedType {
		panic("problem")
	}
	u.references = append(u.references, ref)
}

func (u *LocalTypeDefinition) FetchPositionLength() token.SourceFileReference {
	return u.identifier.LocalType().FetchPositionLength()
}

func (u *LocalTypeDefinition) HumanReadable() string {
	return fmt.Sprintf("%v", u.referencedType.HumanReadable())
}

func (u *LocalTypeDefinition) Identifier() *dtype.LocalTypeName {
	return u.identifier
}

func (u *LocalTypeDefinition) WantsToBeReplaced() bool {
	return IsAny(u.referencedType)
}

func (u *LocalTypeDefinition) IsEqual(_ dtype.Atom) error {
	return nil
}

func (u *LocalTypeDefinition) Resolve() (dtype.Atom, error) {
	return u.referencedType.Resolve()
}

func (u *LocalTypeDefinition) ReferencedType() dtype.Type {
	return u.referencedType
}

func (u *LocalTypeDefinition) ParameterCount() int {
	return 0
}

func (u *LocalTypeDefinition) Next() dtype.Type {
	return u.referencedType
}

func (u *LocalTypeDefinition) WasReferenced() bool {
	return u.wasReferenced
}

func (u *LocalTypeDefinition) MarkAsReferenced() {
	u.wasReferenced = true
}

func (u *LocalTypeDefinition) Verify(referencedType dtype.Type) {
	if referencedType == nil {
		return
	}

	ref, wasRef := referencedType.(*LocalTypeDefinitionReference)
	if wasRef {
		log.Printf("referencedType %v", ref.identifier)
		if ref.typeDefinition == u {
			panic("circular")
		}
	}

	//	TypeChain(referencedType, 0)
}

func (u *LocalTypeDefinition) SetDefinition(referencedType dtype.Type) error {
	u.Verify(referencedType)
	u.referencedType = referencedType
	u.hasBeenDefined = true
	return nil
}

func NewLocalTypeDefinition(identifier *dtype.LocalTypeName) *LocalTypeDefinition {
	//TypeChain(referencedType, 0)
	x := &LocalTypeDefinition{identifier: identifier, referencedType: NewAnyType()}
	x.hasBeenDefined = false
	return x
}
