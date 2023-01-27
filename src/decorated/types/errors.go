package dectype

import (
	"fmt"
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type CouldNotFindLocalTypeName struct {
	err  error
	name *ast.LocalTypeNameReference
}

func NewCouldNotFindLocalTypeName(name *ast.LocalTypeNameReference, err error) *CouldNotFindLocalTypeName {
	return &CouldNotFindLocalTypeName{name: name, err: err}
}

func (e *CouldNotFindLocalTypeName) Error() string {
	return fmt.Sprintf("could not find local type name %v %v", e.name, e.err)
}

func (e *CouldNotFindLocalTypeName) FetchPositionLength() token.SourceFileReference {
	return e.name.FetchPositionLength()
}
