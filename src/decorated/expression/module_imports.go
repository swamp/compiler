package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
)

type ImportedModule struct {
	wasReferenced           bool
	referencedModule        *Module
	importStatementInModule *ImportStatement
	moduleName              *ast.ModuleReference
}

func NewImportedModule(referencedModule *Module, importStatementInModule *ImportStatement) *ImportedModule {
	return &ImportedModule{
		referencedModule:        referencedModule,
		importStatementInModule: importStatementInModule,
		moduleName:              referencedModule.FullyQualifiedModuleName().Path(),
	}
}

func (i *ImportedModule) MarkAsReferenced() {
	i.wasReferenced = true
	i.importStatementInModule.MarkAsReferenced()
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

func (i *ImportedModule) ImportStatementInModule() *ImportStatement {
	return i.importStatementInModule
}

func (i *ImportedModule) ModuleName() *ast.ModuleReference {
	return i.moduleName
}

type ModuleImports struct {
	modules    map[string]*ImportedModule
	allModules []*ImportedModule
}

func NewModuleImports() *ModuleImports {
	return &ModuleImports{modules: make(map[string]*ImportedModule)}
}

func (m *ModuleImports) ImportModule(moduleName *ast.ModuleReference, module *Module, importStatementInModule *ImportStatement) *ImportedModule {
	importedModule := &ImportedModule{
		moduleName: moduleName, referencedModule: module,
		importStatementInModule: importStatementInModule,
	}

	if importStatementInModule.astImport == nil {
		panic("not allowed to add this")
	}
	m.modules[moduleName.ModuleName()] = importedModule
	m.allModules = append(m.allModules, importedModule)

	return importedModule
}

func (m *ModuleImports) FindModule(moduleName *ast.ModuleReference) *ImportedModule {
	return m.modules[moduleName.ModuleName()]
}

func (m *ModuleImports) AllInOrderModules() []*ImportedModule {
	return m.allModules
}

func (m *ModuleImports) String() string {
	s := ""
	for name, module := range m.modules {
		s += fmt.Sprintf(" '%v' : %v", name, module.referencedModule.FullyQualifiedModuleName())
	}

	return s
}
