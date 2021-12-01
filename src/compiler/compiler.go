/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package swampcompiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/swamp/assembler/lib/assembler_sp"
	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/environment"
	"github.com/swamp/compiler/src/file"
	"github.com/swamp/compiler/src/generate_sp"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/solution"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
	swampdisasmsp "github.com/swamp/disassembler/lib"
	opcode_sp_type "github.com/swamp/opcodes/type"
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
		for _, def := range module.LocalTypes().AllTypes() {
			if !def.WasReferenced() {
				warning := decorated.NewUnusedTypeWarning(def)
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

	CheckUnused(world)

	rootModule, err := deccy.CreateDefaultRootModule(true)
	if err != nil {
		return nil, nil, err
	}
	for _, importedRootSubModule := range rootModule.ImportedModules().AllModules() {
		world.AddModule(importedRootSubModule.ReferencedModule().FullyQualifiedModuleName(), importedRootSubModule.ReferencedModule())
	}
	// world.AddModule(dectype.MakeArtifactFullyQualifiedModuleName(nil), rootModule)

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

	err := typeinfo.GeneratePackageToChunk(compiledPackage, typeInformationChunk)
	if err != nil {
		return decorated.NewInternalError(err)
	}

	for _, module := range compiledPackage.AllModules() {
		if verboseFlag >= verbosity.High {
			fmt.Printf(">>> has module %v\n", module.FullyQualifiedModuleName())
		}
	}

	for _, module := range compiledPackage.AllModules() {
		for _, named := range module.LocalDefinitions().Definitions() {
			unknownType := named.Expression()
			maybeFunction, _ := unknownType.(*decorated.FunctionValue)
			if maybeFunction == nil {
				// log.Printf("WHAT IS THIS!_ %T %v", unknownType, unknownType)
			}
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

					for index, param := range maybeFunction.Parameters() {
						unaliased := dectype.Unalias(maybeFunction.Parameters()[index].Type())
						if dectype.ArgumentNeedsTypeIdInsertedBefore(unaliased) || dectype.IsTypeIdRef(unaliased) {
							pos = align(pos, dectype.MemoryAlign(opcode_sp_type.AlignOfSwampInt))
							typeIndexPosRange := assembler_sp.SourceStackPosRange{
								Pos:  assembler_sp.SourceStackPos(pos),
								Size: assembler_sp.SourceStackRange(opcode_sp_type.SizeofSwampInt),
							}
							paramPosRanges = append(paramPosRanges, typeIndexPosRange)
							pos += dectype.MemoryOffset(typeIndexPosRange.Size)
							if dectype.IsTypeIdRef(unaliased) {
								continue
							}
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

					pos := dectype.MemoryOffset(0)
					for _, param := range maybeFunction.Parameters() {
						paramSize, paramAlign := dectype.GetMemorySizeAndAlignment(param.Type())
						pos = align(pos, paramAlign)
						pos += dectype.MemoryOffset(paramSize)
					}
					parameterOctetSize := dectype.MemorySize(pos)
					if _, err := packageConstants.AllocatePrepareFunctionConstant(fullyQualifiedName.String(), opcode_sp_type.MemorySize(returnSize), opcode_sp_type.MemoryAlign(returnAlign), parameterCount, opcode_sp_type.MemorySize(parameterOctetSize), uint(functionTypeIndex)); err != nil {
						return decorated.NewInternalError(err)
					}
				}
			} else {
				// log.Printf("Unknown thing here:%T\n", unknownType)
			}
		}
	}

	var constants *assembler_sp.PackageConstants
	for _, module := range compiledPackage.AllModules() {
		if verboseFlag >= verbosity.High {
			fmt.Printf("============================================== generating for module %v\n", module)
		}

		createdConstants, functions, genErr := gen.GenerateAllLocalDefinedFunctions(module, typeInformationChunk, packageConstants, verboseFlag)
		if genErr != nil {
			return decorated.NewInternalError(genErr)
		}
		constants = createdConstants

		allFunctions = append(allFunctions, functions...)
	}

	if verboseFlag >= verbosity.Mid || showAssembler {
		constants.DynamicMemory().DebugOutput()
	}

	if verboseFlag >= verbosity.Mid || showAssembler {
		var assemblerOutput string

		for _, f := range allFunctions {
			if f.ConstantType() == assembler_sp.ConstantTypeFunction {
				opcodes := constants.FetchOpcodes(f)
				lines := swampdisasmsp.Disassemble(opcodes)

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

	constants.Finalize()
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
