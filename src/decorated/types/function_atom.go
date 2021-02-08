/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
)

type FunctionAtom struct {
	parameterTypes []dtype.Type
}

func NewFunctionAtom(parameterTypes []dtype.Type) *FunctionAtom {
	if len(parameterTypes) < 1 {
		//panic("function atoms must have at least a return type (1).")
	}

	for _, param := range parameterTypes {
		if param == nil {
			panic("stop here")
		}
	}
	return &FunctionAtom{parameterTypes: parameterTypes}
}

func (u *FunctionAtom) FunctionParameterTypes() []dtype.Type {
	return u.parameterTypes
}

func (u *FunctionAtom) ParameterAndReturn() ([]dtype.Type, dtype.Type) {
	count := len(u.parameterTypes)
	ret := u.parameterTypes[count-1]
	params := u.parameterTypes[:count-1]
	return params, ret
}

func (u *FunctionAtom) ReturnType() dtype.Type {
	_, returnType := u.ParameterAndReturn()
	return returnType
}

func (u *FunctionAtom) Resolve() (dtype.Atom, error) {
	return u, nil
}

func (u *FunctionAtom) Next() dtype.Type {
	return nil
}

func (u *FunctionAtom) ParameterCount() int {
	return len(u.parameterTypes)
}

func (u *FunctionAtom) Apply(params []dtype.Type) (dtype.Type, error) {
	return u, nil
}

func (u *FunctionAtom) Generate(params []dtype.Type) (dtype.Type, error) {
	return u, nil
}

func (u *FunctionAtom) String() string {
	return fmt.Sprintf("[functype %v]", u.parameterTypes)
}

func (u *FunctionAtom) ShortString() string {
	s := "[func "
	for _, param := range u.parameterTypes {
		s += " " + param.ShortString()
	}
	s += "]"
	return s
}

func (u *FunctionAtom) HumanReadable() string {
	str := "("
	for index, parameterType := range u.parameterTypes {
		if index > 0 {
			str += " -> "
		}
		str += parameterType.HumanReadable()
	}
	str += ")"
	str += " "
	return str
}

func (u *FunctionAtom) DecoratedName() string {
	s := "func("
	for index, param := range u.parameterTypes {
		if index > 0 {
			s += " -> "
		}
		s += param.DecoratedName()
	}
	s += ")"
	return s
}

func (u *FunctionAtom) AtomName() string {
	s := "func("
	for index, param := range u.parameterTypes {
		if index > 0 {
			s += " -> "
		}
		s += param.ShortName()
	}
	s += ")"
	return s
}

func (u *FunctionAtom) ShortName() string {
	s := "func("
	for index, param := range u.parameterTypes {
		if index > 0 {
			s += " -> "
		}
		s += param.ShortName()
	}
	s += ")"
	return s
}

type FunctionAtomMismatch struct {
	Expected dtype.Atom
	Encountered dtype.Atom
}

func (e FunctionAtomMismatch) Error() string {
	return fmt.Sprintf("expected %v, but encountered %v", e.Expected, e.Encountered)
}

func (u *FunctionAtom) IsEqual(other_ dtype.Atom) error {
	other, wasFunctionAtom := other_.(*FunctionAtom)
	if !wasFunctionAtom {
		return fmt.Errorf("wasnt a function %v", other)
	}

	otherParams := other.parameterTypes
	if len(u.parameterTypes) != len(otherParams) {
		return fmt.Errorf("different argument count ")
	}

	for index, parameter := range u.parameterTypes {
		otherParam, otherParamErr := otherParams[index].Resolve()
		if otherParamErr != nil {
			return fmt.Errorf("parameter couldn't resolve %v %w", otherParam, otherParamErr)
		}
		param, paramErr := parameter.Resolve()
		if paramErr != nil {
			return fmt.Errorf("couldn't resolve it %w", paramErr)
		}

		equalErr := param.IsEqual(otherParam)
		if equalErr != nil {
			return FunctionAtomMismatch{param, otherParam}
		}
	}

	return nil
}
