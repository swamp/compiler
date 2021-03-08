package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type FunctionTypeReference struct {
	referencedType  *FunctionAtom
	astFunctionType *ast.FunctionType
}

func NewFunctionTypeReference(astFunctionType *ast.FunctionType, referencedType *FunctionAtom) *FunctionTypeReference {
	return &FunctionTypeReference{
		referencedType:  referencedType,
		astFunctionType: astFunctionType,
	}
}

func (t *FunctionTypeReference) FetchPositionLength() token.SourceFileReference {
	return t.astFunctionType.FetchPositionLength()
}

func (t *FunctionTypeReference) HumanReadable() string {
	return t.referencedType.HumanReadable()
}

func (t *FunctionTypeReference) ShortString() string {
	return fmt.Sprintf("fntyperef %v %v", t.astFunctionType, t.referencedType)
}

func (t *FunctionTypeReference) ShortName() string {
	return fmt.Sprintf("fntyperef %v %v", t.astFunctionType, t.referencedType)
}

func (t *FunctionTypeReference) String() string {
	return fmt.Sprintf("fntyperef %v %v", t.astFunctionType, t.referencedType)
}

func (t *FunctionTypeReference) Resolve() (dtype.Atom, error) {
	return t.referencedType.Resolve()
}

func (t *FunctionTypeReference) Next() dtype.Type {
	return t.referencedType
}

func (t *FunctionTypeReference) ReturnType() dtype.Type {
	return t.referencedType.ReturnType()
}

func (t *FunctionTypeReference) FunctionAtom() *FunctionAtom {
	return t.referencedType
}

func (t *FunctionTypeReference) ParameterAndReturn() ([]dtype.Type, dtype.Type) {
	return t.referencedType.ParameterAndReturn()
}

func (t *FunctionTypeReference) DecoratedName() string {
	return fmt.Sprintf("fntyperef %v %v", t.astFunctionType, t.referencedType)
}

func (t *FunctionTypeReference) ParameterCount() int {
	return t.referencedType.ParameterCount()
}

func (t *FunctionTypeReference) Generate(params []dtype.Type) (dtype.Type, error) {
	return t.referencedType.Generate(params)
}