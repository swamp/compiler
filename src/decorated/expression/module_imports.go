package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
)

type ModuleImports struct {
	modules map[string]*Module
}

func NewModuleImports() *ModuleImports {
	return &ModuleImports{modules: make(map[string]*Module)}
}

func (m *ModuleImports) ImportModule(moduleName *ast.ModuleReference, module *Module) {
	m.modules[moduleName.ModuleName()] = module
}

func (m *ModuleImports) FindModule(moduleName *ast.ModuleReference) *Module {
	return m.modules[moduleName.ModuleName()]
}

func (m *ModuleImports) AllModules() map[string]*Module {
	return m.modules
}

func (m *ModuleImports) String() string {
	s := ""
	for name, module := range m.modules {
		s += fmt.Sprintf(" '%v' : %v", name, module.FullyQualifiedModuleName())
	}

	return s
}
