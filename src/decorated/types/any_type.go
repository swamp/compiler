/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"github.com/swamp/compiler/src/decorated/dtype"
)

import "fmt"

type Any struct {
}

func (u *Any) String() string {
	return fmt.Sprintf("[any]")
}

func (u *Any) HumanReadable() string {
	return fmt.Sprintf("ANY")
}

func (u *Any) ShortString() string {
	return fmt.Sprintf("[any]")
}

func (u *Any) DecoratedName() string {
	return "any"
}

func (u *Any) ShortName() string {
	return u.DecoratedName()
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

func (u *Any) Apply(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("Any can not be applied")
}

func (u *Any) Generate(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("Any can not be applied")
}


func (u *Any) Resolve() (dtype.Atom, error) {
	return u, nil
}

func (u *Any) Next() dtype.Type {
	return nil
}

func NewAnyType() *Any {
	return &Any{}
}
