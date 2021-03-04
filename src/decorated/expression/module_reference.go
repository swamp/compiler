package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type ModuleReference struct {
	module     *Module
	identifier []*ast.TypeIdentifier
	inclusive  token.SourceFileReference
}

func NewModuleReference(identifier []*ast.TypeIdentifier, module *Module) *ModuleReference {
	inclusive := token.MakeInclusiveSourceFileReference(identifier[0].FetchPositionLength(), identifier[len(identifier)-1].FetchPositionLength())
	ref := &ModuleReference{module: module, identifier: identifier, inclusive: inclusive}

	module.AddReference(ref)

	return ref
}

func (m *ModuleReference) String() string {
	return fmt.Sprintf("moduleref %v", m.module)
}

func (m *ModuleReference) Module() *Module {
	return m.module
}

func (m *ModuleReference) Identifiers() []*ast.TypeIdentifier {
	return m.identifier
}

func (m *ModuleReference) HumanReadable() string {
	return "Module Reference"
}

func (m *ModuleReference) FetchPositionLength() token.SourceFileReference {
	return m.inclusive
}
