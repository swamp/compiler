/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/verbosity"
)

type ModuleReader interface {
	ReadModule(repository deccy.ModuleRepository, moduleName dectype.PackageRelativeModuleName, namespacePrefix dectype.PackageRootModuleName) (*decorated.Module, decshared.DecoratedError)
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

func (l *ModuleRepository) FetchModuleInPackageEx(artifactFullyModuleName dectype.ArtifactFullyQualifiedModuleName, packageRelativeModuleName dectype.PackageRelativeModuleName, verboseFlag verbosity.Verbosity) (*decorated.Module, decshared.DecoratedError) {
	if verboseFlag >= verbosity.Mid {
		fmt.Printf("* fetching module %v artifactName:%v\n", packageRelativeModuleName, artifactFullyModuleName)
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
	// fmt.Printf("world:%v\n", l.world)

	if verboseFlag >= verbosity.Mid {
		fmt.Printf("* didn't have %v (%v), must load and parse\n", artifactFullyModuleName, packageRelativeModuleName)
	}
	if l.isReadingModule(packageRelativeModuleName) {
		return nil, decorated.NewInternalError(NewCircularDependencyDetected(l.resolutionModules, artifactFullyModuleName))
	}
	l.resolutionModules = append(l.resolutionModules, packageRelativeModuleName)
	readModule, readModuleErr := l.moduleReader.ReadModule(l, packageRelativeModuleName, l.moduleNamespace)
	if readModuleErr != nil {
		_, isModuleErrAlready := readModuleErr.(*decorated.ModuleError)
		if isModuleErrAlready {
			return nil, readModuleErr
		}
		return nil, decorated.NewModuleError(artifactFullyModuleName.String()+".swamp", readModuleErr)
	}
	l.world.AddModule(artifactFullyModuleName, readModule)

	l.resolutionModules = remove(l.resolutionModules, packageRelativeModuleName)

	return readModule, nil
}

func (l *ModuleRepository) FetchModuleInPackage(packageRelativeModuleName dectype.PackageRelativeModuleName, verboseFlag verbosity.Verbosity) (*decorated.Module, decshared.DecoratedError) {
	artifactFullyModuleName := l.moduleNamespace.Join(packageRelativeModuleName)

	return l.FetchModuleInPackageEx(artifactFullyModuleName, packageRelativeModuleName, verboseFlag)
}

func (l *ModuleRepository) FetchMainModuleInPackage(verboseFlag verbosity.Verbosity) (*decorated.Module, decshared.DecoratedError) {
	emptyPackageRelativeModuleName := dectype.NewPackageRelativeModuleName(nil)
	artifactFullyModuleName := l.moduleNamespace.Join(emptyPackageRelativeModuleName)

	x, err := l.FetchModuleInPackageEx(artifactFullyModuleName, emptyPackageRelativeModuleName, verboseFlag)
	if err != nil {
		return nil, err
	}

	x.LocalDefinitions().FindDefinitionExpression(ast.NewVariableIdentifier(token.NewVariableSymbolToken("main", token.SourceFileReference{}, 0)))
	x.LocalDefinitions().FindDefinitionExpression(ast.NewVariableIdentifier(token.NewVariableSymbolToken("init", token.SourceFileReference{}, 0)))

	return x, nil
}
