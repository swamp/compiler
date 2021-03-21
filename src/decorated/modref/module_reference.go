package modref

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type ModuleReferencer interface {
	FetchPositionLength() token.SourceFileReference
	AstModuleReference() *ast.ModuleReference
	String() string
}
