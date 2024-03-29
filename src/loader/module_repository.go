/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"log"

	"github.com/swamp/compiler/src/parser"

	"github.com/swamp/compiler/src/ast"
	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/verbosity"
)

type ModuleReader interface {
	ReadModule(moduleType decorated.ModuleType, repository deccy.ModuleRepository, moduleName dectype.PackageRelativeModuleName, namespacePrefix dectype.PackageRootModuleName) (*decorated.Module, decshared.DecoratedError)
}

type ModuleRepository struct {
	moduleReader      ModuleReader
	world             *Package
	moduleNamespace   dectype.PackageRootModuleName
	resolutionModules []dectype.PackageRelativeModuleName
}

func NewModuleRepository(world *Package, moduleNamespace dectype.PackageRootModuleName, moduleReader ModuleReader) *ModuleRepository {
	return &ModuleRepository{world: world, moduleNamespace: moduleNamespace, moduleReader: moduleReader}
}

func (l *ModuleRepository) isReadingModule(moduleName dectype.PackageRelativeModuleName) bool {
	for _, readingModuleName := range l.resolutionModules {
		if readingModuleName.String() == moduleName.String() {
			return true
		}
	}
	return false
}

func (l *ModuleRepository) InternalReader() ModuleReader {
	return l.moduleReader
}

func remove(s []dectype.PackageRelativeModuleName, r dectype.PackageRelativeModuleName) []dectype.PackageRelativeModuleName {
	for i, v := range s {
		if v.String() == r.String() {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (l *ModuleRepository) FetchModuleInPackageEx(moduleType decorated.ModuleType, artifactFullyModuleName dectype.ArtifactFullyQualifiedModuleName, packageRelativeModuleName dectype.PackageRelativeModuleName, verboseFlag verbosity.Verbosity) (*decorated.Module, decshared.DecoratedError) {
	if verboseFlag >= verbosity.Mid {
		log.Printf("* fetching module '%v' artifactName:'%v'\n", packageRelativeModuleName, artifactFullyModuleName)
	}

	module := l.world.FindModule(artifactFullyModuleName)
	if module != nil {
		return module, nil
	}

	secondTry := dectype.MakeArtifactFullyQualifiedModuleName(packageRelativeModuleName.Path())
	module = l.world.FindModule(secondTry)
	if module != nil {
		return module, nil
	}

	if verboseFlag >= verbosity.Mid {
		log.Printf("* didn't have %v (%v), must load and parse\n", artifactFullyModuleName, packageRelativeModuleName)
	}
	if l.isReadingModule(packageRelativeModuleName) {
		return nil, decorated.NewCircularDependencyDetected(packageRelativeModuleName, l.resolutionModules, artifactFullyModuleName)
	}

	var errors decshared.DecoratedError
	l.resolutionModules = append(l.resolutionModules, packageRelativeModuleName)
	readModule, readModuleErr := l.moduleReader.ReadModule(moduleType, l, packageRelativeModuleName, l.moduleNamespace)
	if readModuleErr != nil {
		_, isModuleErrAlready := readModuleErr.(*decorated.ModuleError)
		if isModuleErrAlready {
			return nil, readModuleErr
		}
		if parser.IsCompileErr(readModuleErr) {
			return nil, decorated.NewModuleError(artifactFullyModuleName.String()+".swamp", readModuleErr)
		}
		errors = decorated.AppendError(errors, readModuleErr)
	}
	l.world.AddModule(artifactFullyModuleName, readModule)

	l.resolutionModules = remove(l.resolutionModules, packageRelativeModuleName)

	return readModule, errors
}

func (l *ModuleRepository) FetchModuleInPackage(parentModuleType decorated.ModuleType, packageRelativeModuleName dectype.PackageRelativeModuleName, verboseFlag verbosity.Verbosity) (*decorated.Module, decshared.DecoratedError) {
	artifactFullyModuleName := l.moduleNamespace.Join(packageRelativeModuleName)

	return l.FetchModuleInPackageEx(parentModuleType, artifactFullyModuleName, packageRelativeModuleName, verboseFlag)
}

func (l *ModuleRepository) FetchMainModuleInPackage(moduleType decorated.ModuleType, verboseFlag verbosity.Verbosity) (*decorated.Module, decshared.DecoratedError) {
	emptyPackageRelativeModuleName := dectype.NewPackageRelativeModuleName(nil)
	artifactFullyModuleName := l.moduleNamespace.Join(emptyPackageRelativeModuleName)

	x, err := l.FetchModuleInPackageEx(moduleType, artifactFullyModuleName, emptyPackageRelativeModuleName, verboseFlag)
	if parser.IsCompileError(err) {
		return nil, err
	}

	x.LocalDefinitions().FindDefinitionExpression(ast.NewVariableIdentifier(token.NewVariableSymbolToken("main", token.SourceFileReference{}, 0)))
	x.LocalDefinitions().FindDefinitionExpression(ast.NewVariableIdentifier(token.NewVariableSymbolToken("init", token.SourceFileReference{}, 0)))

	return x, err
}
