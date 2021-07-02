package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
)

type ImportedModule struct {
	wasReferenced           bool
	referencedModule        *Module
	importStatementInModule *Module
	moduleName              *ast.ModuleReference
}

func NewImportedModule(referencedModule *Module, importStatementInModule *Module) *ImportedModule {
	return &ImportedModule{
		referencedModule:        referencedModule,
		importStatementInModule: importStatementInModule,
		moduleName:              referencedModule.FullyQualifiedModuleName().Path(),
	}
}

func (i *ImportedModule) MarkAsReferenced() {
	i.wasReferenced = true
}

func (i *ImportedModule) WasReferenced() bool {
	return i.wasReferenced
}

func (i *ImportedModule) String() string {
	return fmt.Sprintf("%v (%v)\n", i.referencedModule.FullyQualifiedModuleName(), i.wasReferenced)
}

func (i *ImportedModule) ReferencedModule() *Module {
	return i.referencedModule
}

func (i *ImportedModule) ImportStatementInModule() *Module {
	return i.importStatementInModule
}

func (i *ImportedModule) ModuleName() *ast.ModuleReference {
	return i.moduleName
}

type ModuleImports struct {
	modules map[string]*ImportedModule
}

func NewModuleImports() *ModuleImports {
	return &ModuleImports{modules: make(map[string]*ImportedModule)}
}

func (m *ModuleImports) ImportModule(moduleName *ast.ModuleReference, module *Module, importStatementInModule *Module) *ImportedModule {
	importedModule := &ImportedModule{
		moduleName: moduleName, referencedModule: module,
		importStatementInModule: importStatementInModule,
	}

	m.modules[moduleName.ModuleName()] = importedModule

	return importedModule
}

func (m *ModuleImports) FindModule(moduleName *ast.ModuleReference) *ImportedModule {
	return m.modules[moduleName.ModuleName()]
}

func (m *ModuleImports) AllModules() map[string]*ImportedModule {
	return m.modules
}

func (m *ModuleImports) String() string {
	s := ""
	for name, module := range m.modules {
		s += fmt.Sprintf(" '%v' : %v", name, module.referencedModule.FullyQualifiedModuleName())
	}

	return s
}
