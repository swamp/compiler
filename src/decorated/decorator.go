/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/parser"

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
	errors              decshared.DecoratedError
	importedModule      *decorated.ImportedModule
}

func NewDecorator(moduleRepository ModuleRepository, module *decorated.Module, importedModule *decorated.ImportedModule, typeLookUpAndCreate decorated.TypeAddAndReferenceMaker) *Decorator {
	d := &Decorator{module: module, moduleRepository: moduleRepository, importedModule: importedModule, typeLookUpAndCreate: typeLookUpAndCreate}
	return d
}

func (d *Decorator) Import(importStatement *decorated.ImportStatement) error {
	return ImportModuleToModule(d.module, importStatement)
}

func (d *Decorator) TypeReferenceMaker() decorated.TypeAddAndReferenceMaker {
	return d.typeLookUpAndCreate
}

func (d *Decorator) ModuleDefinitions() *decorated.ModuleDefinitions {
	return d.module.LocalDefinitions()
}

func (d *Decorator) AddDeclaration(identifier *ast.VariableIdentifier, ofType dtype.Type) error {
	return d.module.Declarations().AddDeclaration(identifier, ofType)
}

func (d *Decorator) FindNamedFunctionValue(identifier *ast.VariableIdentifier) *decorated.FunctionValue {
	definition := d.module.LocalDefinitions().FindDefinitionExpression(identifier)
	if definition == nil {
		log.Fatalf("couldn't find expression %v", identifier)
	}
	findNamedFunction, wasNamedFunction := definition.Expression().(*decorated.FunctionValue)
	if !wasNamedFunction {
		panic(fmt.Errorf("this was not a function value"))
	}

	return findNamedFunction
}

func (d *Decorator) AddDefinition(identifier *ast.VariableIdentifier, expr decorated.Expression) error {
	return d.module.LocalDefinitions().AddDecoratedExpression(identifier, d.importedModule, expr)
}

func (d *Decorator) AddDecoratedError(decoratedError decshared.DecoratedError) {
	d.errors = decorated.AppendError(d.errors, decoratedError)
}

func (d *Decorator) Errors() decshared.DecoratedError {
	return d.errors
}

func (d *Decorator) NewVariableContext() *decorator.VariableContext {
	return decorator.NewVariableContext(d.module.LocalAndImportedDefinitions())
}

func (d *Decorator) ImportModule(moduleType decorated.ModuleType, importAst *ast.Import, relativeModuleName dectype.PackageRelativeModuleName, alias dectype.SingleModuleName, exposeAll bool, verboseFlag verbosity.Verbosity) (*decorated.ImportStatement, decshared.DecoratedError) {
	moduleToImport, importErr := d.moduleRepository.FetchModuleInPackage(moduleType, relativeModuleName, verboseFlag)
	var appendedError decshared.DecoratedError

	if importErr != nil {
		if parser.IsCompileErr(importErr) {
			return nil, importErr
		}
		appendedError = decorated.AppendError(appendedError, importErr)
	}

	if moduleToImport == nil {
		panic("no module to import (DecorateImport)")
	}

	moduleRef := decorated.NewModuleReference(importAst.ModuleName(), moduleToImport)
	var moduleAliasRef *decorated.ModuleReference
	if !alias.IsEmpty() {
		moduleAliasRef = decorated.NewModuleReference(alias.Path(), moduleToImport)
	}

	importStatement := decorated.NewImport(importAst, moduleRef, moduleAliasRef, exposeAll)

	importModuleErr := d.Import(importStatement)
	if importModuleErr != nil {
		if parser.IsCompileErr(importErr) {
			return nil, decorated.NewInternalError(importModuleErr)
		}

		appendedError = decorated.AppendError(appendedError, importErr)
	}

	return importStatement, appendedError
}
