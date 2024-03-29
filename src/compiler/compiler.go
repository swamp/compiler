/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package swampcompiler

import (
	"fmt"
	"github.com/swamp/compiler/src/semantic"
	"log"
	"os"
	"path"
	"path/filepath"

	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/environment"
	"github.com/swamp/compiler/src/file"
	"github.com/swamp/compiler/src/generate"
	"github.com/swamp/compiler/src/generate_ir"
	"github.com/swamp/compiler/src/generate_sp"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/parser"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/solution"
	"github.com/swamp/compiler/src/verbosity"
)

func CheckUnused(world *loader.Package) decshared.DecoratedError {
	var appendedError decshared.DecoratedError
	for _, module := range world.AllModules() {
		if module.IsInternal() {
			continue
		}
		for _, def := range module.LocalDefinitions().Definitions() {
			if !def.WasReferenced() {
				warning := decorated.NewUnusedWarning(def)
				appendedError = decorated.AppendError(appendedError, warning)
			}
		}
		for _, def := range module.LocalTypes().AllInOrderTypes() {
			if !def.RealType().WasReferenced() {
				warning := decorated.NewUnusedTypeWarning(def.RealType())
				appendedError = decorated.AppendError(appendedError, warning)
			}
		}

		for _, importModule := range module.ImportedModules().AllInOrderModules() {
			if !importModule.ImportStatementInModule().WasReferenced() && importModule.ImportStatementInModule().AstImport().ModuleName().ModuleName() != "" {
				//warning := decorated.NewUnusedImportStatementWarning(importModule.ImportStatementInModule())
				//module.AddWarning(warning)
			}
		}
	}
	return appendedError
}

type Target uint8

const (
	SwampOpcode Target = iota
	LlvmIr
)

func BuildMain(mainSourceFile string, absoluteOutputDirectory string, enforceStyle bool, showAssembler bool, target Target, verboseFlag verbosity.Verbosity) ([]*loader.Package, error) {
	statInfo, statErr := os.Stat(mainSourceFile)
	if statErr != nil {
		return nil, statErr
	}

	config, _, configErr := environment.LoadFromConfig()
	if configErr != nil {
		return nil, configErr
	}

	if statInfo.IsDir() {
		resourceNameLookup := resourceid.NewResourceNameLookupImpl()
		if solutionSettings, err := solution.LoadIfExists(mainSourceFile); err == nil {
			var packages []*loader.Package
			var errors decshared.DecoratedError
			var gen generate.Generator

			if target == LlvmIr {
				gen = generate_ir.NewGenerator()
			} else {
				gen = generate_sp.NewGenerator()
			}

			for _, packageSubDirectoryName := range solutionSettings.Packages {
				absoluteSubDirectory := path.Join(mainSourceFile, packageSubDirectoryName)
				compiledPackage, compileAndLinkErr := CompileAndLink(gen, resourceNameLookup, config, packageSubDirectoryName, absoluteSubDirectory, absoluteOutputDirectory, enforceStyle, verboseFlag, showAssembler)
				errors = decorated.AppendError(errors, compileAndLinkErr)
				if parser.IsCompileError(compileAndLinkErr) {
					return packages, errors
				}

				packages = append(packages, compiledPackage)
			}
			return packages, errors
		} else {
			return nil, fmt.Errorf("must have a solution file in this version")
		}
	}

	return nil, fmt.Errorf("must be directory in this version %v %v", mainSourceFile, absoluteOutputDirectory)
}

func BuildMainOnlyCompile(mainSourceFile string, enforceStyle bool, verboseFlag verbosity.Verbosity) ([]*loader.Package, error) {
	statInfo, statErr := os.Stat(mainSourceFile)
	if statErr != nil {
		return nil, statErr
	}

	config, _, configErr := environment.LoadFromConfig()
	if configErr != nil {
		return nil, configErr
	}

	if !statInfo.IsDir() {
		return nil, fmt.Errorf("must have a solution file in this version")
	} else {
		if solutionSettings, err := solution.LoadIfExists(mainSourceFile); err == nil {
			var packages []*loader.Package
			for _, packageSubDirectoryName := range solutionSettings.Packages {
				absoluteSubDirectory := path.Join(mainSourceFile, packageSubDirectoryName)
				compiledPackage, err := CompileMainDefaultDocumentProvider(packageSubDirectoryName, absoluteSubDirectory, config, enforceStyle, verboseFlag)
				if parser.IsCompileError(err) {
					return packages, err
				}
				packages = append(packages, compiledPackage)
			}
			return packages, nil
		}
	}

	return nil, fmt.Errorf("must be directory in this version %v", mainSourceFile)
}

func CompileMain(name string, mainSourceFile string, documentProvider loader.DocumentProvider, configuration environment.Environment, enforceStyle bool, verboseFlag verbosity.Verbosity) (*loader.Package, *decorated.Module, decshared.DecoratedError) {
	mainPrefix := mainSourceFile
	if file.IsDir(mainSourceFile) {
	} else {
		mainPrefix = path.Dir(mainSourceFile)
	}
	world := loader.NewPackage(loader.LocalFileSystemRoot(mainPrefix), name)

	worldDecorator, worldDecoratorErr := loader.NewWorldDecorator(enforceStyle, verboseFlag)
	if parser.IsCompileErr(worldDecoratorErr) {
		return nil, nil, worldDecoratorErr
	}

	var appendedError decshared.DecoratedError

	appendedError = decorated.AppendError(appendedError, worldDecoratorErr)
	/*
		for _, rootModule := range worldDecorator.ImportModules() {
			world.AddModule(rootModule.FullyQualifiedModuleName(), rootModule)
		}
	*/
	mainNamespace := dectype.MakePackageRootModuleName(nil)

	rootPackage := NewPackageLoader(mainPrefix, documentProvider, mainNamespace, world, worldDecorator)

	libraryReader := loader.NewLibraryReaderAndDecorator()
	libraryModule, libErr := libraryReader.ReadLibraryModule(decorated.ModuleTypeNormal, world, rootPackage.repository, mainSourceFile, mainNamespace, documentProvider, configuration)
	if parser.IsCompileErr(libErr) {
		return nil, nil, libErr
	}
	appendedError = decorated.AppendError(appendedError, libErr)

	unusedErrors := CheckUnused(world)
	appendedError = decorated.AppendError(appendedError, unusedErrors)

	rootModule, err := deccy.CreateDefaultRootModule(true)
	if parser.IsCompileError(err) {
		return nil, nil, err
	}
	appendedError = decorated.AppendError(appendedError, err)

	for _, importedRootSubModule := range rootModule.ImportedModules().AllInOrderModules() {
		world.AddModule(importedRootSubModule.ReferencedModule().FullyQualifiedModuleName(), importedRootSubModule.ReferencedModule())
	}

	_, semanticErr := semantic.GenerateTokensEncodedValues(rootModule.Nodes())
	if semanticErr != nil {
		panic(semanticErr)
		return nil, nil, decorated.NewInternalError(semanticErr)
	}

	return world, libraryModule, appendedError
}

func CompileMainFindLibraryRoot(mainSource string, documentProvider loader.DocumentProvider, configuration environment.Environment, enforceStyle bool, verboseFlag verbosity.Verbosity) (*loader.Package, *decorated.Module, error) {
	if !file.IsDir(mainSource) {
		mainSource = filepath.Dir(mainSource)
	}

	libraryDirectory, libraryErr := loader.FindSettingsDirectory(mainSource)
	if libraryErr != nil {
		return nil, nil, fmt.Errorf("couldn't find settings directory when compiling %w", libraryErr)
	}

	return CompileMain(mainSource, libraryDirectory, documentProvider, configuration, enforceStyle, verboseFlag)
}

type CoreFunctionInfo struct {
	Name       string
	ParamCount uint
}

func align(offset dectype.MemoryOffset, memoryAlign dectype.MemoryAlign) dectype.MemoryOffset {
	rest := dectype.MemoryAlign(uint32(offset) % uint32(memoryAlign))
	if rest != 0 {
		offset += dectype.MemoryOffset(memoryAlign - rest)
	}
	return offset
}

func GenerateAndLink(gen generate.Generator, resourceNameLookup resourceid.ResourceNameLookup, compiledPackage *loader.Package, outputDirectory string, packageSubDirectory string, verboseFlag verbosity.Verbosity, showAssembler bool) decshared.DecoratedError {
	for _, module := range compiledPackage.AllModules() {
		if verboseFlag >= verbosity.High {
			log.Printf(">>> has module %v\n", module.FullyQualifiedModuleName())
		}
	}

	genErr := gen.GenerateFromPackageAndWriteOutput(compiledPackage, resourceNameLookup, outputDirectory, packageSubDirectory, verboseFlag, showAssembler)
	if genErr != nil {
		return decorated.NewInternalError(genErr)
	}

	return nil
}

func CompileMainDefaultDocumentProvider(name string, filename string, configuration environment.Environment,
	enforceStyle bool, verboseFlag verbosity.Verbosity) (*loader.Package, decshared.DecoratedError) {
	defaultDocumentProvider := loader.NewFileSystemDocumentProvider()

	compiledPackage, _, moduleErr := CompileMain(name, filename, defaultDocumentProvider, configuration, enforceStyle, verboseFlag)
	if moduleErr != nil {
		return compiledPackage, moduleErr
	}

	return compiledPackage, nil
}

func CompileAndLink(gen generate.Generator, resourceNameLookup resourceid.ResourceNameLookup, configuration environment.Environment, name string,
	filename string, outputFilename string, enforceStyle bool, verboseFlag verbosity.Verbosity, showAssembler bool) (*loader.Package, decshared.DecoratedError) {
	var errors decshared.DecoratedError
	compiledPackage, compileErr := CompileMainDefaultDocumentProvider(name, filename, configuration, enforceStyle, verboseFlag)
	if parser.IsCompileError(compileErr) {
		return nil, compileErr
	}
	errors = decorated.AppendError(errors, compileErr)
	if compiledPackage == nil {
		panic("not possible")
	}

	linkErr := GenerateAndLink(gen, resourceNameLookup, compiledPackage, outputFilename, name, verboseFlag, showAssembler)
	errors = decorated.AppendError(errors, linkErr)

	return compiledPackage, errors
}
