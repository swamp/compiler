/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"github.com/swamp/compiler/src/ast"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

type ModuleRepository interface {
	FetchModuleInPackage(moduleName dectype.PackageRelativeModuleName, verboseFlag bool) (*decorated.Module, decshared.DecoratedError)
}

type Decorator struct {
	module           *decorated.Module
	moduleRepository ModuleRepository
}

func NewDecorator(moduleRepository ModuleRepository, module *decorated.Module) *Decorator {
	d := &Decorator{module: module, moduleRepository: moduleRepository}
	return d
}

func (d *Decorator) Import(source *decorated.Module, relativeName dectype.PackageRelativeModuleName, exposeAll bool) error {
	return ImportModuleToModule(d.module, source, relativeName, exposeAll)
}

func (d *Decorator) TypeRepo() *dectype.TypeRepo {
	return d.module.TypeRepo()
}

func (d *Decorator) ModuleDefinitions() *decorated.ModuleDefinitions {
	return d.module.Definitions()
}

func (d *Decorator) AddDeclaration(identifier *ast.VariableIdentifier, ofType dtype.Type) error {
	return d.module.Declarations().AddDeclaration(identifier, ofType)
}

func (d *Decorator) AddDefinition(identifier *ast.VariableIdentifier, expr decorated.Expression) error {
	return d.module.Definitions().AddDecoratedExpression(identifier, expr)
}

func (d *Decorator) NewVariableContext() *decorator.VariableContext {
	return decorator.NewVariableContext(d.module.LocalAndImportedDefinitions())
}

func (d *Decorator) ImportModule(importAst *ast.Import, relativeModuleName dectype.PackageRelativeModuleName, alias dectype.SingleModuleName, exposeAll bool, verboseFlag bool) (*decorated.ImportStatement, decshared.DecoratedError) {
	moduleToImport, importErr := d.moduleRepository.FetchModuleInPackage(relativeModuleName, verboseFlag)
	if importErr != nil {
		return nil, importErr
	}

	if moduleToImport == nil {
		panic("no module to import (DecorateImport)")
	}

	if !alias.IsEmpty() {
		relativeModuleName = dectype.MakePackageRelativeModuleName(alias.Path())
	}

	importStatement := decorated.NewImport(importAst, moduleToImport)

	importModuleErr := d.Import(moduleToImport, relativeModuleName, exposeAll)
	if importModuleErr != nil {
		return nil, decorated.NewInternalError(importModuleErr)
	}

	return importStatement, nil
}

func (d *Decorator) AddExternalFunction(function *ast.ExternalFunction) (*decorated.ExternalFunctionDeclaration, decshared.DecoratedError) {
	return d.module.AddExternalFunction(function), nil
}
