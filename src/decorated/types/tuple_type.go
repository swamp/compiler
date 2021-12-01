package dectype

import (
	"fmt"
	"reflect"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

func calculateTupleFieldOffsetsAndRecordMemorySizeAndAlign(fields []*TupleTypeField) (MemorySize, MemoryAlign) {
	offset := MemoryOffset(0)
	maxMemoryAlign := MemoryAlign(0)

	for _, field := range fields {
		memorySize, memoryAlign := GetMemorySizeAndAlignment(field.fieldType)
		rest := MemoryAlign(uint32(offset) % uint32(memoryAlign))
		if rest != 0 {
			offset += MemoryOffset(memoryAlign - rest)
		}
		if memoryAlign > maxMemoryAlign {
			maxMemoryAlign = memoryAlign
		}

		field.memoryOffset = offset
		field.memorySize = memorySize

		offset += MemoryOffset(memorySize)
	}

	rest := MemoryAlign(uint32(offset) % uint32(maxMemoryAlign))
	if rest != 0 {
		offset += MemoryOffset(maxMemoryAlign - rest)
	}

	return MemorySize(offset), maxMemoryAlign
}

type TupleTypeAtom struct {
	parameterFields []*TupleTypeField
	parameterTypes  []dtype.Type
	astTupleType    *ast.TupleType
	memorySize      MemorySize
	memoryAlign     MemoryAlign
}

func NewTupleTypeAtom(astTupleType *ast.TupleType, parameterFields []*TupleTypeField) *TupleTypeAtom {
	for _, param := range parameterFields {
		if reflect.TypeOf(param) == nil {
			panic("function atom: nil parameter type")
		}
	}

	var parameterTypes []dtype.Type
	for _, param := range parameterFields {
		parameterTypes = append(parameterTypes, param.Type())
	}

	memorySize, memoryAlign := calculateTupleFieldOffsetsAndRecordMemorySizeAndAlign(parameterFields)

	return &TupleTypeAtom{
		parameterFields: parameterFields, parameterTypes: parameterTypes, astTupleType: astTupleType,
		memorySize: memorySize, memoryAlign: memoryAlign,
	}
}

func (u *TupleTypeAtom) MemorySize() MemorySize {
	return u.memorySize
}

func (u *TupleTypeAtom) Fields() []*TupleTypeField {
	return u.parameterFields
}

func (u *TupleTypeAtom) MemoryAlignment() MemoryAlign {
	return u.memoryAlign
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

func (u *TupleTypeAtom) WasReferenced() bool {
	return false // tuples are not reused
}
