package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type TypeReferenceScoped struct {
	referencedType dtype.Type
	typeIdentfier  *ast.TypeIdentifierScoped
}

func NewTypeReferenceScoped(typeIdentfier *ast.TypeIdentifierScoped, referencedType dtype.Type) *TypeReferenceScoped {
	return &TypeReferenceScoped{
		referencedType: referencedType,
		typeIdentfier:  typeIdentfier,
	}
}

func (t *TypeReferenceScoped) FetchPositionLength() token.SourceFileReference {
	return t.typeIdentfier.FetchPositionLength()
}

func (t *TypeReferenceScoped) HumanReadable() string {
	return t.referencedType.HumanReadable()
}

func (t *TypeReferenceScoped) ShortString() string {
	return fmt.Sprintf("typeref %v %v", t.typeIdentfier, t.referencedType)
}

func (t *TypeReferenceScoped) ShortName() string {
	return fmt.Sprintf("typeref %v %v", t.typeIdentfier, t.referencedType)
}

func (t *TypeReferenceScoped) String() string {
	return fmt.Sprintf("typeref %v %v", t.typeIdentfier, t.referencedType)
}

func (t *TypeReferenceScoped) Resolve() (dtype.Atom, error) {
	return t.referencedType.Resolve()
}

func (t *TypeReferenceScoped) Next() dtype.Type {
	return t.referencedType
}

func (t *TypeReferenceScoped) DecoratedName() string {
	return fmt.Sprintf("typeref %v %v", t.typeIdentfier, t.referencedType)
}

func (t *TypeReferenceScoped) ParameterCount() int {
	return t.referencedType.ParameterCount()
}

func (t *TypeReferenceScoped) Generate(params []dtype.Type) (dtype.Type, error) {
	return t.referencedType.Generate(params)
}
