/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package swampcompiler

import (
	"fmt"
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
	"log"
	"os"
	"path"
	"path/filepath"
)

func CheckUnused(world *loader.Package) []decshared.DecoratedError {
	var errors []decshared.DecoratedError
	for _, module := range world.AllModules() {
		if module.IsInternal() {
			continue
		}
		for _, def := range module.LocalDefinitions().Definitions() {
			if !def.WasReferenced() {
				warning := decorated.NewUnusedWarning(def)
				errors = append(errors, warning)
			}
		}
		for _, def := range module.LocalTypes().AllTypes() {
			if !def.WasReferenced() {
				warning := decorated.NewUnusedTypeWarning(def)
				errors = append(errors, warning)
			}
		}

		for _, importModule := range module.ImportedModules().AllModules() {
			if !importModule.ImportStatementInModule().WasReferenced() && importModule.ImportStatementInModule().AstImport().ModuleName().ModuleName() != "" {
				//warning := decorated.NewUnusedImportStatementWarning(importModule.ImportStatementInModule())
				//module.AddWarning(warning)
			}
		}
	}
	return errors
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
			var errors []decshared.DecoratedError
			var gen generate.Generator

			if target == LlvmIr {
				gen = generate_ir.NewGenerator()
			} else {
				gen = generate_sp.NewGenerator()
			}

			for _, packageSubDirectoryName := range solutionSettings.Packages {
				absoluteSubDirectory := path.Join(mainSourceFile, packageSubDirectoryName)
				compiledPackage, err := CompileAndLink(gen, resourceNameLookup, config, packageSubDirectoryName, absoluteSubDirectory, absoluteOutputDirectory, enforceStyle, verboseFlag)
				if err != nil {
					errors = append(errors, err)
				}
				if parser.IsCompileError(err) {
					return packages, decorated.NewMultiErrors(errors)
				}

				packages = append(packages, compiledPackage)
			}
			var returnErr decshared.DecoratedError
			if len(errors) > 0 {
				returnErr = decorated.NewMultiErrors(errors)
			}
			return packages, returnErr
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
	if worldDecoratorErr != nil {
		return nil, nil, worldDecoratorErr
	}
	/*
		for _, rootModule := range worldDecorator.ImportModules() {
			world.AddModule(rootModule.FullyQualifiedModuleName(), rootModule)
		}
	*/
	mainNamespace := dectype.MakePackageRootModuleName(nil)

	rootPackage := NewPackageLoader(mainPrefix, documentProvider, mainNamespace, world, worldDecorator)

	libraryReader := loader.NewLibraryReaderAndDecorator()
	libraryModule, libErr := libraryReader.ReadLibraryModule(decorated.ModuleTypeNormal, world, rootPackage.repository, mainSourceFile, mainNamespace, documentProvider, configuration)
	if libErr != nil {
		return nil, nil, libErr
	}
	// color.Cyan(fmt.Sprintf("=> importing package %v as top package", mainPrefix))

	unusedErrors := CheckUnused(world)
	var returnErr decshared.DecoratedError
	if len(unusedErrors) > 0 {
		returnErr = decorated.NewMultiErrors(unusedErrors)
	}

	rootModule, err := deccy.CreateDefaultRootModule(true)
	if err != nil {
		return nil, nil, err
	}
	for _, importedRootSubModule := range rootModule.ImportedModules().AllModules() {
		world.AddModule(importedRootSubModule.ReferencedModule().FullyQualifiedModuleName(), importedRootSubModule.ReferencedModule())
	}
	// world.AddModule(dectype.MakeArtifactFullyQualifiedModuleName(nil), rootModule)

	return world, libraryModule, returnErr
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

func GenerateAndLink(gen generate.Generator, resourceNameLookup resourceid.ResourceNameLookup, compiledPackage *loader.Package, outputDirectory string, packageSubDirectory string, verboseFlag verbosity.Verbosity) decshared.DecoratedError {

	for _, module := range compiledPackage.AllModules() {
		if verboseFlag >= verbosity.High {
			log.Printf(">>> has module %v\n", module.FullyQualifiedModuleName())
		}
	}

	genErr := gen.GenerateFromPackage(compiledPackage, resourceNameLookup, outputDirectory, packageSubDirectory, verboseFlag)
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
	filename string, outputFilename string, enforceStyle bool, verboseFlag verbosity.Verbosity) (*loader.Package, decshared.DecoratedError) {
	var errors []decshared.DecoratedError
	compiledPackage, compileErr := CompileMainDefaultDocumentProvider(name, filename, configuration, enforceStyle, verboseFlag)
	if parser.IsCompileError(compileErr) {
		return nil, compileErr
	}
	if compileErr != nil {
		errors = append(errors, compileErr)
	}

	linkErr := GenerateAndLink(gen, resourceNameLookup, compiledPackage, outputFilename, name, verboseFlag)
	if linkErr != nil {
		errors = append(errors, linkErr)
	}

	var returnErr decshared.DecoratedError
	if len(errors) > 0 {
		returnErr = decorated.NewMultiErrors(errors)
	}

	return compiledPackage, returnErr
}
