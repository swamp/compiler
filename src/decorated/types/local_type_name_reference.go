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

type LocalTypeNameReference struct {
	identifier     *ast.LocalTypeNameReference
	referencedType *LocalTypeName `debug:"true"`
}

func (u *LocalTypeNameReference) String() string {
	return fmt.Sprintf("[LocalTypeNameRef %v]", u.identifier.Name())
}

func (u *LocalTypeNameReference) FetchPositionLength() token.SourceFileReference {
	return u.identifier.FetchPositionLength()
}

func (u *LocalTypeNameReference) HumanReadable() string {
	return fmt.Sprintf("%v", u.identifier.Name())
}

func (u *LocalTypeNameReference) Identifier() *ast.LocalTypeNameReference {
	return u.identifier
}

func (u *LocalTypeNameReference) IsEqual(_ dtype.Atom) error {
	return nil
}

func (u *LocalTypeNameReference) Resolve() (dtype.Atom, error) {
	return NewAnyType(), nil
}

func (u *LocalTypeNameReference) ReferencedType() dtype.Type {
	return NewAnyType()
}

func (u *LocalTypeNameReference) ParameterCount() int {
	return 0
}

func (u *LocalTypeNameReference) Next() dtype.Type {
	return NewAnyType()
}

func (u *LocalTypeNameReference) WasReferenced() bool {
	return false
}

func NewLocalTypeNameReference(identifier *ast.LocalTypeNameReference,
	referencedType *LocalTypeName) *LocalTypeNameReference {
	return &LocalTypeNameReference{identifier: identifier, referencedType: referencedType}
}
