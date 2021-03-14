/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package swampcompiler

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/swamp/compiler/src/ast"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/file"
	"github.com/swamp/compiler/src/generate"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/typeinfo"

	swampdisasm "github.com/swamp/disassembler/lib"
)

func CheckUnused(world *loader.Package) {
	for _, module := range world.AllModules() {
		if module.IsInternal() {
			continue
		}
		for _, def := range module.Definitions().Definitions() {
			if !def.WasReferenced() {
				module := def.ParentDefinitions().OwnedByModule()
				warning := decorated.NewUnusedWarning(def)
				module.AddWarning(warning)
			}
		}
	}
}

func CompileMain(mainSourceFile string, documentProvider loader.DocumentProvider, enforceStyle bool, verboseFlag bool) (*loader.Package, *decorated.Module, decshared.DecoratedError) {
	mainPrefix := mainSourceFile
	if file.IsDir(mainSourceFile) {
	} else {
		mainPrefix = filepath.Dir(mainSourceFile)
	}
	world := loader.NewPackage(loader.LocalFileSystemRoot(mainPrefix))

	worldDecorator, worldDecoratorErr := loader.NewWorldDecorator(enforceStyle, verboseFlag)
	if worldDecoratorErr != nil {
		return nil, nil, worldDecoratorErr
	}
	for _, rootModule := range worldDecorator.ImportModules() {
		world.AddModule(rootModule.FullyQualifiedModuleName(), rootModule)
	}

	mainNamespace := dectype.MakePackageRootModuleName(nil)

	rootPackage := NewPackageLoader(mainPrefix, documentProvider, mainNamespace, world, worldDecorator)

	libraryReader := loader.NewLibraryReaderAndDecorator()
	libraryModule, libErr := libraryReader.ReadLibraryModule(world, rootPackage.repository, mainSourceFile, mainNamespace, documentProvider)
	if libErr != nil {
		return nil, nil, libErr
	}
	// color.Cyan(fmt.Sprintf("=> importing package %v as top package", mainPrefix))

	CheckUnused(world)

	return world, libraryModule, nil
}

func CompileMainFindLibraryRoot(mainSource string, documentProvider loader.DocumentProvider, enforceStyle bool, verboseFlag bool) (*loader.Package, *decorated.Module, error) {
	if !file.IsDir(mainSource) {
		mainSource = filepath.Dir(mainSource)
	}
	libraryDirectory, libraryErr := loader.FindSettingsDirectory(mainSource)
	if libraryErr != nil {
		return nil, nil, libraryErr
	}

	return CompileMain(libraryDirectory, documentProvider, enforceStyle, verboseFlag)
}

type CoreFunctionInfo struct {
	Name       string
	ParamCount uint
}

func GenerateAndLink(world *loader.Package, outputFilename string, verboseFlag bool) decshared.DecoratedError {
	gen := generate.NewGenerator()
	var allFunctions []*generate.Function
	var allExternalFunctions []*generate.ExternalFunction
	fakeMod := decorated.NewModule(dectype.MakeArtifactFullyQualifiedModuleName(nil), nil)

	typeInformationChunk, err := typeinfo.Generate(world)
	if err != nil {
		return decorated.NewInternalError(err)
	}

	for _, module := range world.AllModules() {
		if verboseFlag {
			fmt.Printf(">>> has module %v\n", module.FullyQualifiedModuleName())
		}
	}

	for _, module := range world.AllModules() {
		if verboseFlag {
			fmt.Printf("============================================== generating for module %v\n", module)
		}
		context := decorator.NewVariableContext(module.LocalAndImportedDefinitions())
		functions, genErr := gen.GenerateAllLocalDefinedFunctions(module, context, typeInformationChunk, verboseFlag)
		if genErr != nil {
			return decorated.NewInternalError(genErr)
		}
		allFunctions = append(allFunctions, functions...)
		externalFunctions := module.ExternalFunctions()
		for _, externalFunction := range externalFunctions {
			fakeName := decorated.NewFullyQualifiedVariableName(fakeMod, ast.NewVariableIdentifier(token.NewVariableSymbolToken(externalFunction.AstExternalFunction.FunctionName(), token.SourceFileReference{}, 0)))
			fakeFunc := generate.NewExternalFunction(fakeName, 0, externalFunction.AstExternalFunction.ParameterCount())
			allExternalFunctions = append(allExternalFunctions, fakeFunc)
		}
	}

	if verboseFlag {
		var assemblerOutput string
		for _, f := range allFunctions {
			lines := swampdisasm.Disassemble(f.Opcodes())
			assemblerOutput = assemblerOutput + fmt.Sprintf("func %v\n%s\n\n", f, strings.Join(lines[:], "\n"))
		}
		fmt.Println(assemblerOutput)
	}

	typeInformationOctets, typeInformationErr := typeinfo.ChunkToOctets(typeInformationChunk)
	if typeInformationErr != nil {
		return decorated.NewInternalError(typeInformationErr)
	}

	packed, packedErr := generate.Pack(allFunctions, allExternalFunctions, typeInformationOctets, typeInformationChunk)
	if packedErr != nil {
		return decorated.NewInternalError(packedErr)
	}

	if err := ioutil.WriteFile(outputFilename, packed, 0o644); err != nil {
		return decorated.NewInternalError(err)
	}

	// color.Cyan("wrote output file '%v'\n", outputFilename)
	return nil
}

func CompileAndLink(filename string, outputFilename string, enforceStyle bool, verboseFlag bool) decshared.DecoratedError {
	defaultDocumentProvider := loader.NewFileSystemDocumentProvider()
	world, _, moduleErr := CompileMain(filename, defaultDocumentProvider, enforceStyle, verboseFlag)
	if moduleErr != nil {
		return moduleErr
	}

	return GenerateAndLink(world, outputFilename, verboseFlag)
}
