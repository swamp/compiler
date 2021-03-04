package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type ModuleReferenceFull struct {
	module     *Module
	identifier *ast.TypeIdentifier
}

func NewModuleReferenceFull(identifier *ast.TypeIdentifier, module *Module) *ModuleReferenceFull {
	ref := &ModuleReferenceFull{module: module, identifier: identifier}

	return ref
}

func (m *ModuleReferenceFull) String() string {
	return fmt.Sprintf("modulereffull %v", m.module)
}

func (m *ModuleReferenceFull) Module() *Module {
	return m.module
}

func (m *ModuleReferenceFull) Ident() *ast.TypeIdentifier {
	return m.identifier
}

func (m *ModuleReferenceFull) HumanReadable() string {
	return "Module Reference"
}

func (m *ModuleReferenceFull) FetchPositionLength() token.SourceFileReference {
	return m.identifier.FetchPositionLength()
}
