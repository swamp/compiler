package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type NamedCustomType struct {
	customTypeAtom dtype.Type
	inclusive      token.SourceFileReference
}

func NewNamedCustomType(identifier *ast.TypeIdentifier, value dtype.Type) *NamedCustomType {
	inclusive := token.MakeInclusiveSourceFileReference(identifier.FetchPositionLength(), value.FetchPositionLength())
	return &NamedCustomType{
		customTypeAtom: value,
		inclusive:      inclusive,
	}
}

func (n *NamedCustomType) String() string {
	return fmt.Sprintf("named function customTypeAtom  %v", n.customTypeAtom)
}

func (n *NamedCustomType) StatementString() string {
	return fmt.Sprintf("named function customTypeAtom  %v", n.customTypeAtom)
}

func (n *NamedCustomType) Type() dtype.Type {
	return n.customTypeAtom
}

func (n *NamedCustomType) FetchPositionLength() token.SourceFileReference {
	return n.inclusive
}
