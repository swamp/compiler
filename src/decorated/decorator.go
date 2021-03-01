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
	nodes            []decorated.Node
}

func NewDecorator(moduleRepository ModuleRepository, module *decorated.Module) *Decorator {
	d := &Decorator{module: module, moduleRepository: moduleRepository}
	return d
}

func (d *Decorator) InternalAddNode(node decorated.Node) {
	_, isRef := node.(*dectype.TypeReference)
	if isRef {
		isRef = false
	}
	d.nodes = append(d.nodes, node)
}

func (d *Decorator) RootNodes() []decorated.Node {
	return d.nodes
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

func (d *Decorator) AddImport(importAst *ast.Import, relativeModuleName dectype.PackageRelativeModuleName, alias dectype.SingleModuleName, exposeAll bool, verboseFlag bool) decshared.DecoratedError {
	moduleToImport, importErr := d.moduleRepository.FetchModuleInPackage(relativeModuleName, verboseFlag)
	if importErr != nil {
		return importErr
	}
	if moduleToImport == nil {
		panic("no module to import (AddImport)")
	}
	if !alias.IsEmpty() {
		relativeModuleName = dectype.MakePackageRelativeModuleName(alias.Path())
	}

	importStatement := decorated.NewImport(importAst, moduleToImport)
	d.InternalAddNode(importStatement)

	importModuleErr := d.Import(moduleToImport, relativeModuleName, exposeAll)
	if importModuleErr != nil {
		return decorated.NewInternalError(importModuleErr)
	}

	return nil
}

func (d *Decorator) AddExternalFunction(name string, parameterCount uint) decshared.DecoratedError {
	d.module.AddExternalFunction(name, parameterCount)
	return nil
}
