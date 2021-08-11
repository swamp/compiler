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
	"github.com/swamp/compiler/src/verbosity"
)

type ModuleRepository interface {
	FetchModuleInPackage(parentModuleType decorated.ModuleType, moduleName dectype.PackageRelativeModuleName, verboseFlag verbosity.Verbosity) (*decorated.Module, decshared.DecoratedError)
}

type Decorator struct {
	module              *decorated.Module
	moduleRepository    ModuleRepository
	typeLookUpAndCreate decorated.TypeAddAndReferenceMaker
	errors              []decshared.DecoratedError
}

func NewDecorator(moduleRepository ModuleRepository, module *decorated.Module, typeLookUpAndCreate decorated.TypeAddAndReferenceMaker) *Decorator {
	d := &Decorator{module: module, moduleRepository: moduleRepository, typeLookUpAndCreate: typeLookUpAndCreate}
	return d
}

func (d *Decorator) Import(importStatement *decorated.ImportStatement) error {
	return ImportModuleToModule(d.module, importStatement)
}

func (d *Decorator) TypeRepo() decorated.TypeAddAndReferenceMaker {
	return d.typeLookUpAndCreate
}

func (d *Decorator) ModuleDefinitions() *decorated.ModuleDefinitions {
	return d.module.LocalDefinitions()
}

func (d *Decorator) AddDeclaration(identifier *ast.VariableIdentifier, ofType dtype.Type) error {
	return d.module.Declarations().AddDeclaration(identifier, ofType)
}

func (d *Decorator) AddDefinition(identifier *ast.VariableIdentifier, expr decorated.Expression) error {
	return d.module.LocalDefinitions().AddDecoratedExpression(identifier, expr)
}

func (d *Decorator) AddDecoratedError(decoratedError decshared.DecoratedError) {
	d.errors = append(d.errors, decoratedError)
}

func (d *Decorator) Errors() []decshared.DecoratedError {
	return d.errors
}

func (d *Decorator) NewVariableContext() *decorator.VariableContext {
	return decorator.NewVariableContext(d.module.LocalAndImportedDefinitions())
}

func (d *Decorator) ImportModule(moduleType decorated.ModuleType, importAst *ast.Import, relativeModuleName dectype.PackageRelativeModuleName, alias dectype.SingleModuleName, exposeAll bool, verboseFlag verbosity.Verbosity) (*decorated.ImportStatement, decshared.DecoratedError) {
	moduleToImport, importErr := d.moduleRepository.FetchModuleInPackage(moduleType, relativeModuleName, verboseFlag)
	if importErr != nil {
		return nil, importErr
	}

	if moduleToImport == nil {
		panic("no module to import (DecorateImport)")
	}

	moduleRef := decorated.NewModuleReference(importAst.ModuleName(), moduleToImport)
	var moduleAliasRef *decorated.ModuleReference
	if !alias.IsEmpty() {
		relativeModuleName = dectype.MakePackageRelativeModuleName(alias.Path())
		moduleAliasRef = decorated.NewModuleReference(alias.Path(), moduleToImport)
	}

	importStatement := decorated.NewImport(importAst, moduleRef, moduleAliasRef, exposeAll)

	importModuleErr := d.Import(importStatement)
	if importModuleErr != nil {
		return nil, decorated.NewInternalError(importModuleErr)
	}

	return importStatement, nil
}

func (d *Decorator) AddExternalFunction(function *ast.ExternalFunction) (*decorated.ExternalFunctionDeclaration, decshared.DecoratedError) {
	return d.module.AddExternalFunction(function), nil
}
