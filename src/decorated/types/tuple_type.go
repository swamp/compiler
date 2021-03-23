package dectype

import (
	"fmt"
	"reflect"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type TupleTypeAtom struct {
	parameterTypes []dtype.Type
	astTupleType   *ast.TupleType
}

func NewTupleTypeAtom(astTupleType *ast.TupleType, parameterTypes []dtype.Type) *TupleTypeAtom {
	for _, param := range parameterTypes {
		if reflect.TypeOf(param) == nil {
			panic("function atom: nil parameter type")
		}
	}
	return &TupleTypeAtom{parameterTypes: parameterTypes, astTupleType: astTupleType}
}

func (u *TupleTypeAtom) ParameterTypes() []dtype.Type {
	return u.parameterTypes
}

func (u *TupleTypeAtom) ParameterAndReturn() ([]dtype.Type, dtype.Type) {
	count := len(u.parameterTypes)
	ret := u.parameterTypes[count-1]
	params := u.parameterTypes[:count-1]
	return params, ret
}

func (u *TupleTypeAtom) ReturnType() dtype.Type {
	_, returnType := u.ParameterAndReturn()
	return returnType
}

func (u *TupleTypeAtom) Resolve() (dtype.Atom, error) {
	return u, nil
}

func (u *TupleTypeAtom) Next() dtype.Type {
	return nil
}

func (u *TupleTypeAtom) ParameterCount() int {
	return len(u.parameterTypes)
}

func (u *TupleTypeAtom) String() string {
	return fmt.Sprintf("[tupletype %v]", u.parameterTypes)
}

func (u *TupleTypeAtom) HumanReadable() string {
	str := "("
	for index, parameterType := range u.parameterTypes {
		if index > 0 {
			str += ", "
		}
		str += parameterType.HumanReadable()
	}
	str += ")"
	str += " "
	return str
}

func (u *TupleTypeAtom) AtomName() string {
	s := "tuple("
	for index, param := range u.parameterTypes {
		if index > 0 {
			s += ", "
		}
		s += param.String()
	}
	s += ")"
	return s
}

func (u *TupleTypeAtom) FetchPositionLength() token.SourceFileReference {
	return u.astTupleType.FetchPositionLength()
}

func (u *TupleTypeAtom) IsEqual(other_ dtype.Atom) error {
	other, wasFunctionAtom := other_.(*TupleTypeAtom)
	if !wasFunctionAtom {
		return fmt.Errorf("wasnt a tuple %v", other)
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
