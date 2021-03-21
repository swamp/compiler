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

type Any struct {
	astTypeIdentifier *ast.TypeIdentifier
}

func (u *Any) String() string {
	return fmt.Sprintf("[any]")
}

func (u *Any) HumanReadable() string {
	return fmt.Sprintf("ANY")
}

func (u *Any) DecoratedName() string {
	return "any"
}

func (u *Any) AtomName() string {
	return u.DecoratedName()
}

func (u *Any) IsEqual(_ dtype.Atom) error {
	return nil
}

func (u *Any) ParameterCount() int {
	return 0
}

func (u *Any) Resolve() (dtype.Atom, error) {
	return u, nil
}

func (u *Any) Next() dtype.Type {
	return nil
}

func (u *Any) FetchPositionLength() token.SourceFileReference {
	return u.astTypeIdentifier.FetchPositionLength()
}

func NewAnyType(astTypeIdentifier *ast.TypeIdentifier) *Any {
	return &Any{astTypeIdentifier: astTypeIdentifier}
}
