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

type LocalType struct {
	identifier *ast.TypeParameter
}

func (u *LocalType) String() string {
	return fmt.Sprintf("[localtype %v]", u.identifier.Name())
}

func (u *LocalType) FetchPositionLength() token.SourceFileReference {
	return u.identifier.Identifier().FetchPositionLength()
}

func (u *LocalType) HumanReadable() string {
	return fmt.Sprintf("%v", u.identifier.Name())
}

func (u *LocalType) Identifier() *ast.TypeParameter {
	return u.identifier
}

func (u *LocalType) ShortString() string {
	return fmt.Sprintf("[localtype %v]", u.identifier.Name())
}

func (u *LocalType) DecoratedName() string {
	return u.identifier.Name()
}

func (u *LocalType) ShortName() string {
	return u.DecoratedName()
}

func (u *LocalType) AtomName() string {
	return u.DecoratedName()
}

func (u *LocalType) IsEqual(_ dtype.Atom) error {
	return nil
}

func (u *LocalType) ParameterCount() int {
	return 0
}

func (u *LocalType) Apply(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("LocalType can not be applied")
}

func (u *LocalType) Generate(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("LocalType can not be applied")
}

func (u *LocalType) Resolve() (dtype.Atom, error) {
	return u, nil
}

func (u *LocalType) Next() dtype.Type {
	return nil
}

func NewLocalType(identifier *ast.TypeParameter) *LocalType {
	return &LocalType{identifier: identifier}
}
