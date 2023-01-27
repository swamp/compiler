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

type LocalTypeNameDefinition struct {
	identifier    *dtype.LocalTypeName
	wasReferenced bool
}

func (u *LocalTypeNameDefinition) String() string {
	return fmt.Sprintf("[GenericParam %v]", u.identifier.Name())
}

func (u *LocalTypeNameDefinition) FetchPositionLength() token.SourceFileReference {
	return u.identifier.LocalType().FetchPositionLength()
}

func (u *LocalTypeNameDefinition) HumanReadable() string {
	return fmt.Sprintf("%v", u.identifier.Name())
}

func (u *LocalTypeNameDefinition) Identifier() *dtype.LocalTypeName {
	return u.identifier
}

func (u *LocalTypeNameDefinition) IsEqual(_ dtype.Atom) error {
	return nil
}

func (u *LocalTypeNameDefinition) WasReferenced() bool {
	return u.wasReferenced
}

func (u *LocalTypeNameDefinition) MarkAsReferenced() {
	u.wasReferenced = true
}

func NewLocalTypeNameDefinition(identifier *dtype.LocalTypeName) *LocalTypeNameDefinition {
	return &LocalTypeNameDefinition{identifier: identifier}
}
