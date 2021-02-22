package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type TypeReference struct {
	referencedType dtype.Type
	typeIdentfier  *ast.TypeIdentifier
}

func NewTypeReference(typeIdentfier *ast.TypeIdentifier, referencedType dtype.Type) *TypeReference {
	return &TypeReference{
		referencedType: referencedType,
		typeIdentfier:  typeIdentfier,
	}
}

func (t *TypeReference) FetchPositionLength() token.SourceFileReference {
	return t.typeIdentfier.FetchPositionLength()
}

func (t *TypeReference) HumanReadable() string {
	return t.referencedType.HumanReadable()
}

func (t *TypeReference) ShortString() string {
	return fmt.Sprintf("typeref %v %v", t.typeIdentfier, t.referencedType)
}

func (t *TypeReference) ShortName() string {
	return fmt.Sprintf("typeref %v %v", t.typeIdentfier, t.referencedType)
}

func (t *TypeReference) String() string {
	return fmt.Sprintf("typeref %v %v", t.typeIdentfier, t.referencedType)
}

func (t *TypeReference) Resolve() (dtype.Atom, error) {
	return t.referencedType.Resolve()
}

func (t *TypeReference) Next() dtype.Type {
	return t.referencedType
}

func (t *TypeReference) DecoratedName() string {
	return fmt.Sprintf("typeref %v %v", t.typeIdentfier, t.referencedType)
}

func (t *TypeReference) ParameterCount() int {
	return t.referencedType.ParameterCount()
}

func (t *TypeReference) Generate(params []dtype.Type) (dtype.Type, error) {
	return t.referencedType.Generate(params)
}
