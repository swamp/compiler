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

type LocalTypeDefinition struct {
	identifier     *dtype.LocalTypeName
	referencedType dtype.Type
	hasBeenDefined bool
	wasReferenced  bool
}

func (u *LocalTypeDefinition) String() string {
	return fmt.Sprintf("[ConcreteGeneric %v %v]", u.identifier.Name(), u.referencedType)
}

func (u *LocalTypeDefinition) FetchPositionLength() token.SourceFileReference {
	return u.identifier.LocalType().FetchPositionLength()
}

func (u *LocalTypeDefinition) HumanReadable() string {
	return fmt.Sprintf("%v", u.identifier.Name())
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

func (u *LocalTypeDefinition) SetDefinition(referencedType dtype.Type) error {
	u.referencedType = referencedType
	return nil
}

func NewLocalTypeDefinition(identifier *dtype.LocalTypeName, referencedType dtype.Type) *LocalTypeDefinition {
	return &LocalTypeDefinition{identifier: identifier, referencedType: referencedType}
}
