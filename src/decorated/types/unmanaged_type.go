/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type UnmanagedType struct {
	identifier    *ast.UnmanagedType
	wasReferenced bool
}

func (u *UnmanagedType) String() string {
	return fmt.Sprintf("[unmanaged %v]", u.identifier.Name())
}

func (u *UnmanagedType) FetchPositionLength() token.SourceFileReference {
	return u.identifier.FetchPositionLength()
}

func (u *UnmanagedType) HumanReadable() string {
	return fmt.Sprintf("%v", u.identifier.Name())
}

func (u *UnmanagedType) Identifier() *ast.UnmanagedType {
	return u.identifier
}

func (u *UnmanagedType) AtomName() string {
	return u.identifier.Name()
}

func (u *UnmanagedType) IsEqual(other_ dtype.Atom) error {
	if IsAtomAny(other_) {
		return nil
	}

	other, wasUnmanaged := other_.(*UnmanagedType)
	if !wasUnmanaged {
		return fmt.Errorf("wasn't unmanaged even %v", other)
	}

	if other.identifier.Name() == u.identifier.Name() {
		return nil
	}

	return fmt.Errorf("not the same unmanaged %v vs %v", u, other)
}

func (u *UnmanagedType) ParameterCount() int {
	return 0
}

func (u *UnmanagedType) Resolve() (dtype.Atom, error) {
	return u, nil
}

func (u *UnmanagedType) Next() dtype.Type {
	return nil
}

func (u *UnmanagedType) WasReferenced() bool {
	return u.wasReferenced
}

func (u *UnmanagedType) MarkAsReferenced() {
	u.wasReferenced = true
}

func NewUnmanagedType(identifier *ast.UnmanagedType) *UnmanagedType {
	return &UnmanagedType{identifier: identifier}
}
