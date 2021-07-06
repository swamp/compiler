/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"reflect"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type FunctionTypeLike interface {
	dtype.Type
	ReturnType() dtype.Type
	ParameterAndReturn() ([]dtype.Type, dtype.Type)
}

type FunctionAtom struct {
	parameterTypes  []dtype.Type
	astFunctionType *ast.FunctionType
}

func NewFunctionAtom(astFunctionType *ast.FunctionType, parameterTypes []dtype.Type) *FunctionAtom {
	for _, param := range parameterTypes {
		if reflect.TypeOf(param) == nil {
			panic("function atom: nil parameter type")
		}
	}
	return &FunctionAtom{parameterTypes: parameterTypes, astFunctionType: astFunctionType}
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

func (u *FunctionAtom) String() string {
	return fmt.Sprintf("[functype %v]", u.parameterTypes)
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

func (u *FunctionAtom) AtomName() string {
	s := "func("
	for index, param := range u.parameterTypes {
		if index > 0 {
			s += " -> "
		}
		s += param.HumanReadable()
	}
	s += ")"
	return s
}

func (u *FunctionAtom) FetchPositionLength() token.SourceFileReference {
	return u.astFunctionType.FetchPositionLength()
}

type FunctionAtomMismatch struct {
	Expected    dtype.Type
	Encountered dtype.Type
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
		equalErr := CompatibleTypes(parameter, otherParams[index])
		if equalErr != nil {
			return FunctionAtomMismatch{parameter, otherParams[index]}
		}
	}

	return nil
}
