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

type LocalTypeName struct {
	identifier    *dtype.LocalTypeName
	wasReferenced bool
}

func (u *LocalTypeName) String() string {
	return fmt.Sprintf("[GenericParam %v]", u.identifier.Name())
}

func (u *LocalTypeName) FetchPositionLength() token.SourceFileReference {
	return u.identifier.LocalType().FetchPositionLength()
}

func (u *LocalTypeName) HumanReadable() string {
	return fmt.Sprintf("%v", u.identifier.Name())
}

func (u *LocalTypeName) Identifier() *dtype.LocalTypeName {
	return u.identifier
}

func (u *LocalTypeName) Name() string {
	return u.identifier.Name()
}

func (u *LocalTypeName) LocalTypeName() *dtype.LocalTypeName {
	return u.identifier
}

func (u *LocalTypeName) IsEqual(_ dtype.Atom) error {
	return nil
}

func (u *LocalTypeName) WasReferenced() bool {
	return u.wasReferenced
}

func (u *LocalTypeName) MarkAsReferenced() {
	u.wasReferenced = true
}

func NewLocalTypeName(identifier *dtype.LocalTypeName) *LocalTypeName {
	return &LocalTypeName{identifier: identifier}
}
