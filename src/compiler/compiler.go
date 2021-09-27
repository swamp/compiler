/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package swampcompiler

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/swamp/compiler/src/assembler_sp"
	"github.com/swamp/compiler/src/ast"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	swampdisasm_sp "github.com/swamp/compiler/src/disassemble_sp"
	"github.com/swamp/compiler/src/environment"
	"github.com/swamp/compiler/src/file"
	"github.com/swamp/compiler/src/generate_sp"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/solution"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
)

func CheckUnused(world *loader.Package) {
	for _, module := range world.AllModules() {
		if module.IsInternal() {
			continue
		}
		for _, def := range module.LocalDefinitions().Definitions() {
			if !def.WasReferenced() {
				module := def.OwnedByModule()
				warning := decorated.NewUnusedWarning(def)
				module.AddWarning(warning)
			}
		}
	}
}

func BuildMain(mainSourceFile string, absoluteOutputDirectory string, enforceStyle bool, showAssembler bool, verboseFlag verbosity.Verbosity) ([]*loader.Package, error) {
	statInfo, statErr := os.Stat(mainSourceFile)
	if statErr != nil {
		return nil, statErr
	}

	config, _, configErr := environment.LoadFromConfig()
	if configErr != nil {
		return nil, configErr
	}

	if statInfo.IsDir() {
		typeInformationChunk := &typeinfo.Chunk{}
		if solutionSettings, err := solution.LoadIfExists(mainSourceFile); err == nil {
			var packages []*loader.Package
			for _, packageSubDirectoryName := range solutionSettings.Packages {
				outputFilename := path.Join(absoluteOutputDirectory, fmt.Sprintf("%s.swamp-pack", packageSubDirectoryName))
				absoluteSubDirectory := path.Join(mainSourceFile, packageSubDirectoryName)
				compiledPackage, err := CompileAndLink(typeInformationChunk, config, packageSubDirectoryName, absoluteSubDirectory, outputFilename, enforceStyle, showAssembler, verboseFlag)
				if err != nil {
					return packages, err
				}
				packages = append(packages, compiledPackage)
			}
			return packages, nil
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
				if err != nil {
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
	for _, rootModule := range worldDecorator.ImportModules() {
		world.AddModule(rootModule.FullyQualifiedModuleName(), rootModule)
	}

	mainNamespace := dectype.MakePackageRootModuleName(nil)

	rootPackage := NewPackageLoader(mainPrefix, documentProvider, mainNamespace, world, worldDecorator)

	libraryReader := loader.NewLibraryReaderAndDecorator()
	libraryModule, libErr := libraryReader.ReadLibraryModule(decorated.ModuleTypeNormal, world, rootPackage.repository, mainSourceFile, mainNamespace, documentProvider, configuration)
	if libErr != nil {
		return nil, nil, libErr
	}
	// color.Cyan(fmt.Sprintf("=> importing package %v as top package", mainPrefix))

	CheckUnused(world)

	return world, libraryModule, nil
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

func GenerateAndLink(typeInformationChunk *typeinfo.Chunk, compiledPackage *loader.Package, outputFilename string, showAssembler bool, verboseFlag verbosity.Verbosity) decshared.DecoratedError {
	gen := generate_sp.NewGenerator()

	var allFunctions []*assembler_sp.Constant

	var allExternalFunctions []*generate_sp.ExternalFunction

	fakeMod := decorated.NewModule(decorated.ModuleTypeNormal, dectype.MakeArtifactFullyQualifiedModuleName(nil), nil)

	err := typeinfo.GeneratePackageToChunk(compiledPackage, typeInformationChunk)
	if err != nil {
		return decorated.NewInternalError(err)
	}

	for _, module := range compiledPackage.AllModules() {
		if verboseFlag >= verbosity.Mid {
			fmt.Printf(">>> has module %v\n", module.FullyQualifiedModuleName())
		}
	}

	packageConstants := assembler_sp.NewPackageConstants()
	for _, module := range compiledPackage.AllModules() {
		for _, named := range module.LocalDefinitions().Definitions() {
			unknownType := named.Expression()
			maybeFunction, _ := unknownType.(*decorated.FunctionValue)
			if maybeFunction != nil {
				fullyQualifiedName := module.FullyQualifiedName(named.Identifier())
				isExternal := maybeFunction.Annotation().Annotation().IsSomeKindOfExternal()
				if isExternal {
					var paramPosRanges []assembler_sp.SourceStackPosRange
					hasLocalTypes := decorated.TypeIsTemplateHasLocalTypes(maybeFunction.ForcedFunctionType())
					// parameterCount := len(maybeFunction.Parameters())
					pos := dectype.MemoryOffset(0)
					log.Printf("this has local types:'%s'\n", named.FullyQualifiedVariableName())
					if hasLocalTypes {
						returnPosRange := assembler_sp.SourceStackPosRange{
							Pos:  assembler_sp.SourceStackPos(0),
							Size: assembler_sp.SourceStackRange(0),
						}
						paramPosRanges = make([]assembler_sp.SourceStackPosRange, len(maybeFunction.Parameters()))
						if _, err := packageConstants.AllocatePrepareExternalFunctionConstant(fullyQualifiedName.String(), returnPosRange, paramPosRanges); err != nil {
							return decorated.NewInternalError(err)
						}
						continue
					}
					returnSize, _ := dectype.GetMemorySizeAndAlignment(maybeFunction.ForcedFunctionType().ReturnType())
					returnPosRange := assembler_sp.SourceStackPosRange{
						Pos:  assembler_sp.SourceStackPos(pos),
						Size: assembler_sp.SourceStackRange(returnSize),
					}

					pos += dectype.MemoryOffset(returnSize)

					log.Printf("generating for function:'%s'\n", named.FullyQualifiedVariableName())
					for index, param := range maybeFunction.Parameters() {
						unaliased := dectype.Unalias(maybeFunction.Parameters()[index].Type())
						if dectype.ArgumentNeedsTypeIdInsertedBefore(unaliased) {
							pos = align(pos, dectype.AlignOfSwampInt)
							typeIndexPosRange := assembler_sp.SourceStackPosRange{
								Pos:  assembler_sp.SourceStackPos(pos),
								Size: assembler_sp.SourceStackRange(dectype.SizeofSwampInt),
							}
							paramPosRanges = append(paramPosRanges, typeIndexPosRange)
							pos += dectype.MemoryOffset(typeIndexPosRange.Size)
						}
						size, alignment := dectype.GetMemorySizeAndAlignment(param.Type())
						pos = align(pos, alignment)
						posRange := assembler_sp.SourceStackPosRange{
							Pos:  assembler_sp.SourceStackPos(pos),
							Size: assembler_sp.SourceStackRange(size),
						}
						paramPosRanges = append(paramPosRanges, posRange)
						pos += dectype.MemoryOffset(size)

					}
					if _, err := packageConstants.AllocatePrepareExternalFunctionConstant(fullyQualifiedName.String(), returnPosRange, paramPosRanges); err != nil {
						return decorated.NewInternalError(err)
					}
				} else {
					returnSize, returnAlign := dectype.GetMemorySizeAndAlignment(maybeFunction.ForcedFunctionType().ReturnType())
					parameterCount := uint(len(maybeFunction.Parameters()))

					functionTypeIndex, lookupErr := typeInformationChunk.Lookup(maybeFunction.ForcedFunctionType())
					if lookupErr != nil {
						return decorated.NewInternalError(lookupErr)
					}

					if _, err := packageConstants.AllocatePrepareFunctionConstant(fullyQualifiedName.String(), returnSize, returnAlign, parameterCount, 0, uint(functionTypeIndex)); err != nil {
						return decorated.NewInternalError(err)
					}
				}
			}
		}
	}

	var constants *assembler_sp.PackageConstants
	for _, module := range compiledPackage.AllModules() {
		if verboseFlag >= verbosity.High {
			fmt.Printf("============================================== generating for module %v\n", module)
		}

		context := decorator.NewVariableContext(module.LocalAndImportedDefinitions())

		createdConstants, functions, genErr := gen.GenerateAllLocalDefinedFunctions(module, context, typeInformationChunk, packageConstants, verboseFlag)
		if genErr != nil {
			return decorated.NewInternalError(genErr)
		}
		constants = createdConstants

		allFunctions = append(allFunctions, functions...)
		externalFunctions := module.ExternalFunctions()

		for _, externalFunction := range externalFunctions {
			fakeName := decorated.NewFullyQualifiedVariableName(fakeMod,
				ast.NewVariableIdentifier(token.NewVariableSymbolToken(externalFunction.AstExternalFunction.FunctionName(),
					token.SourceFileReference{}, 0))) //nolint:exhaustivestruct
			fakeFunc := generate_sp.NewExternalFunction(fakeName, 0, externalFunction.AstExternalFunction.ParameterCount())
			allExternalFunctions = append(allExternalFunctions, fakeFunc)
		}
	}

	if verboseFlag >= verbosity.Mid || showAssembler {
		constants.DynamicMemory().DebugOutput()
	}

	if verboseFlag >= verbosity.Mid || showAssembler {
		var assemblerOutput string

		for _, f := range allFunctions {
			if f.ConstantType() == assembler_sp.ConstantTypeFunction {
				opcodes := constants.FetchOpcodes(f)
				lines := swampdisasm_sp.Disassemble(opcodes)

				assemblerOutput += fmt.Sprintf("func %v\n%s\n\n", f, strings.Join(lines[:], "\n"))
			}
		}

		fmt.Println(assemblerOutput)
	}

	typeInformationOctets, typeInformationErr := typeinfo.ChunkToOctets(typeInformationChunk)
	if typeInformationErr != nil {
		return decorated.NewInternalError(typeInformationErr)
	}

	if verboseFlag >= verbosity.High {
		fmt.Printf("writing type information (%d octets)\n", len(typeInformationOctets))
		typeInformationChunk.DebugOutput()
	}

	dynamicMemoryOctets := constants.DynamicMemory().Octets()

	packed, packedErr := generate_sp.Pack(constants.Constants(), dynamicMemoryOctets, typeInformationOctets)
	if packedErr != nil {
		return decorated.NewInternalError(packedErr)
	}

	if err := ioutil.WriteFile(outputFilename, packed, 0o644); err != nil {
		return decorated.NewInternalError(err)
	}

	// color.Cyan("wrote output file '%v'\n", outputFilename)
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

func CompileAndLink(typeInformationChunk *typeinfo.Chunk, configuration environment.Environment, name string,
	filename string, outputFilename string, enforceStyle bool, showAssembler bool, verboseFlag verbosity.Verbosity) (*loader.Package, decshared.DecoratedError) {
	compiledPackage, compileErr := CompileMainDefaultDocumentProvider(name, filename, configuration, enforceStyle, verboseFlag)
	if compileErr != nil {
		return nil, compileErr
	}

	if err := GenerateAndLink(typeInformationChunk, compiledPackage, outputFilename, showAssembler, verboseFlag); err != nil {
		return compiledPackage, err
	}

	return compiledPackage, nil
}
