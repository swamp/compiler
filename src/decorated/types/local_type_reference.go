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

type LocalTypeDefinitionReference struct {
	identifier     *LocalTypeNameReference
	typeDefinition *LocalTypeDefinition
}

func (u *LocalTypeDefinitionReference) String() string {
	return fmt.Sprintf("[ConcreteGenericRef %v => %v]", u.typeDefinition.identifier, u.typeDefinition.referencedType)
}

func (u *LocalTypeDefinitionReference) FetchPositionLength() token.SourceFileReference {
	return u.identifier.FetchPositionLength()
}

func (u *LocalTypeDefinitionReference) HumanReadable() string {
	return fmt.Sprintf("%v", u.typeDefinition.HumanReadable())
}

func (u *LocalTypeDefinitionReference) TypeDefinition() *LocalTypeDefinition {
	return u.typeDefinition
}

func (u *LocalTypeDefinitionReference) Identifier() *LocalTypeNameReference {
	return u.identifier
}

func (u *LocalTypeDefinitionReference) IsEqual(_ dtype.Atom) error {
	return nil
}

func (u *LocalTypeDefinitionReference) Resolve() (dtype.Atom, error) {
	return u.typeDefinition.Resolve()
}

func (u *LocalTypeDefinitionReference) ReferencedType() dtype.Type {
	return u.typeDefinition
}

func (u *LocalTypeDefinitionReference) ParameterCount() int {
	return 0
}

func (u *LocalTypeDefinitionReference) Next() dtype.Type {
	return u.typeDefinition
}

func (u *LocalTypeDefinitionReference) WasReferenced() bool {
	return false
}

func NewLocalTypeDefinitionReference(identifier *LocalTypeNameReference, referencedDefinition *LocalTypeDefinition) *LocalTypeDefinitionReference {
	x := &LocalTypeDefinitionReference{identifier: identifier, typeDefinition: referencedDefinition}
	referencedDefinition.AddReference(x)

	return x
}
