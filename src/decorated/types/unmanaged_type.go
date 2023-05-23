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
	return fmt.Sprintf("[unmanaged %v]", u.identifier)
}

func (u *UnmanagedType) FetchPositionLength() token.SourceFileReference {
	return u.identifier.FetchPositionLength()
}

func (u *UnmanagedType) HumanReadable() string {
	return fmt.Sprintf("%v", u.identifier.NativeLanguageTypeName().Name())
}

func (u *UnmanagedType) Identifier() *ast.UnmanagedType {
	return u.identifier
}

func (u *UnmanagedType) AtomName() string {
	return u.identifier.Name()
}

func (u *UnmanagedType) IsEqualUnmanaged(other *UnmanagedType) error {
	if other.identifier.NativeLanguageTypeName().Name() == u.identifier.NativeLanguageTypeName().Name() {
		return nil
	}

	return fmt.Errorf("not equal unmanaged %v vs %v", other.identifier.NativeLanguageTypeName(),
		u.identifier.NativeLanguageTypeName())
}

func (u *UnmanagedType) IsEqual(other_ dtype.Atom) error {
	if IsAtomAny(other_) {
		return nil
	}

	other, wasUnmanaged := other_.(*UnmanagedType)
	if !wasUnmanaged {
		return fmt.Errorf("wasn't unmanaged even %T %v", other_, other_)
	}

	if err := u.IsEqualUnmanaged(other); err != nil {
		return err
	}

	return nil
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
