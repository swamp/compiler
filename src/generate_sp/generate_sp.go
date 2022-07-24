/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/swamp/assembler/lib/assembler_sp"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
	swampdisasmsp "github.com/swamp/disassembler/lib"
	"github.com/swamp/opcodes/instruction_sp"
	"github.com/swamp/opcodes/opcode_sp"
)

type AnyPosAndRange interface {
	getPosition() uint32
	getSize() uint32
}

type TypeRef uint16

type Function struct {
	name                *decorated.FullyQualifiedPackageVariableName
	signature           TypeRef
	opcodes             []byte
	debugLines          []opcode_sp.OpcodeInfo
	debugVariables      []opcode_sp.VariableInfo
	debugParameterCount uint
}

type ExternalFunction struct {
	name           *decorated.FullyQualifiedPackageVariableName
	signature      TypeRef
	parameterCount uint
}

func NewFunction(fullyQualifiedName *decorated.FullyQualifiedPackageVariableName, signature TypeRef,
	opcodes []byte, debugParameterCount uint, debugInfos []opcode_sp.OpcodeInfo) *Function {
	f := &Function{
		name: fullyQualifiedName, signature: signature, opcodes: opcodes, debugParameterCount: debugParameterCount,
		debugLines: debugInfos,
	}

	return f
}

func (f *Function) String() string {
	return fmt.Sprintf("[function %v %v]", f.name, f.signature)
}

func (f *Function) Opcodes() []byte {
	return f.opcodes
}

func (f *Function) DebugLines() []opcode_sp.OpcodeInfo {
	return f.debugLines
}

type Generator struct {
	code              *assembler_sp.Code
	packageConstants  *assembler_sp.PackageConstants
	functionConstants []*assembler_sp.Constant
	lookup            typeinfo.TypeLookup
	chunk             *typeinfo.Chunk
	fileUrlCache      *assembler_sp.FileUrlCache
}

func NewGenerator() *Generator {
	g := &Generator{chunk: &typeinfo.Chunk{}, fileUrlCache: assembler_sp.NewFileUrlCache()}
	g.lookup = g.chunk
	return g
}

func (g *Generator) PrepareForNewPackage() {
	g.code = assembler_sp.NewCode()
	g.packageConstants = assembler_sp.NewPackageConstants()
	g.fileUrlCache = assembler_sp.NewFileUrlCache()
}

func (g *Generator) PackageConstants() *assembler_sp.PackageConstants {
	return g.packageConstants
}

func (g *Generator) LastFunctionConstants() []*assembler_sp.Constant {
	return g.functionConstants
}

func arithmeticToUnaryOperatorType(operatorType decorated.ArithmeticUnaryOperatorType) instruction_sp.UnaryOperatorType {
	switch operatorType {
	case decorated.ArithmeticUnaryMinus:
		return instruction_sp.UnaryOperatorNegate
	}

	panic("illegal unaryoperator")
}

func generateConstant(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	constant *decorated.Constant, context *generateContext) error {
	return generateExpression(code, target, constant.Expression(), true, context)
}

func allocMemoryForType(stackMemory *assembler_sp.StackMemoryMapper, typeToAlloc dtype.Type,
	debugString string) assembler_sp.TargetStackPosRange {
	memorySize, alignment := dectype.GetMemorySizeAndAlignment(typeToAlloc)
	if uint(memorySize) == 0 || uint(alignment) == 0 {
		panic(fmt.Errorf("can not allocate zero memory or align for type %T %v", typeToAlloc, typeToAlloc))
	}
	return stackMemory.Allocate(uint(memorySize), uint32(alignment), debugString)
}

func allocMemoryForTypeEx(stackMemory *assembler_sp.StackMemoryMapper, typeToAlloc dtype.Type,
	debugString string) (assembler_sp.TargetStackPosRange, dectype.MemoryAlign) {
	memorySize, alignment := dectype.GetMemorySizeAndAlignment(typeToAlloc)
	if uint(memorySize) == 0 || uint(alignment) == 0 {
		panic(fmt.Errorf("can not allocate zero memory or align for type %T %v", typeToAlloc, typeToAlloc))
	}
	return stackMemory.Allocate(uint(memorySize), uint32(alignment), debugString), alignment
}

func generateRecurCall(code *assembler_sp.Code, call *decorated.RecurCall, genContext *generateContext) error {
	filePosition := genContext.toFilePosition(call.FetchPositionLength())
	code.Recur(filePosition)

	return nil
}

func createTargetWithMemoryOffsetAndSize(target assembler_sp.TargetStackPosRange, memoryOffset uint, size uint) assembler_sp.TargetStackPosRange {
	return assembler_sp.TargetStackPosRange{
		Pos:  assembler_sp.TargetStackPos(uint(target.Pos) + memoryOffset),
		Size: assembler_sp.StackRange(size),
	}
}

const (
	PointerSize  = 8
	PointerAlign = 8
)

func targetToSourceStackPosRange(functionPointer assembler_sp.TargetStackPosRange) assembler_sp.SourceStackPosRange {
	sourcePosRange := assembler_sp.SourceStackPosRange{
		Pos:  assembler_sp.SourceStackPos(functionPointer.Pos),
		Size: assembler_sp.SourceStackRange(functionPointer.Size),
	}

	return sourcePosRange
}

func sourceToTargetStackPosRange(functionPointer assembler_sp.SourceStackPosRange) assembler_sp.TargetStackPosRange {
	targetPosRange := assembler_sp.TargetStackPosRange{
		Pos:  assembler_sp.TargetStackPos(functionPointer.Pos),
		Size: assembler_sp.StackRange(functionPointer.Size),
	}

	return targetPosRange
}

func constantToSourceStackPosRange(code *assembler_sp.Code, stackMemory *assembler_sp.StackMemoryMapper, constant *assembler_sp.Constant) (assembler_sp.SourceStackPosRange, error) {
	functionPointer := stackMemory.Allocate(PointerSize, PointerAlign, "functionReference:"+constant.String())
	code.LoadZeroMemoryPointer(functionPointer.Pos, constant.PosRange().Position, opcode_sp.FilePosition{})

	return targetToSourceStackPosRange(functionPointer), nil
}

func (g *Generator) Before(compilePackage *loader.Package) error {
	g.PrepareForNewPackage()
	err := typeinfo.GeneratePackageToChunk(compilePackage, g.chunk)
	if err != nil {
		return decorated.NewInternalError(err)
	}

	return nil
}

func (g *Generator) GenerateFromPackage(compilePackage *loader.Package, resourceNameLookup resourceid.ResourceNameLookup, absoluteOutputDirectory string, packageSubDirectory string, verboseFlag verbosity.Verbosity) error {
	g.Before(compilePackage)

	packageConstants, allConstantsErr := preparePackageConstants(compilePackage, g.lookup)
	if allConstantsErr != nil {
		return allConstantsErr
	}
	g.packageConstants = packageConstants
	for _, mod := range compilePackage.AllModules() {
		standardError := g.GenerateModule(mod, resourceNameLookup, verboseFlag)
		if standardError != nil {
			return decorated.NewInternalError(standardError)
		}
	}

	const showAssembler = true
	return g.After(resourceNameLookup, absoluteOutputDirectory, packageSubDirectory, showAssembler, verboseFlag)
}

func (g *Generator) After(resourceNameLookup resourceid.ResourceNameLookup, absoluteOutputDirectory string, packageSubDirectory string, showAssembler bool, verboseFlag verbosity.Verbosity) error {
	var allFunctions []*assembler_sp.Constant
	constants := g.packageConstants

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

	constants.AllocateDebugInfoFiles(g.fileUrlCache.FileUrls())
	typeInformationOctets, typeInformationErr := typeinfo.ChunkToOctets(g.chunk)
	if typeInformationErr != nil {
		return decorated.NewInternalError(typeInformationErr)
	}

	if verboseFlag >= verbosity.High {
		log.Printf("writing type information (%d octets)\n", len(typeInformationOctets))
		g.chunk.DebugOutput()
	}

	for _, resourceName := range resourceNameLookup.SortedResourceNames() {
		constants.AllocateResourceNameConstant(resourceName)
	}

	constants.Finalize()
	dynamicMemoryOctets := constants.DynamicMemory().Octets()

	packed, packedErr := Pack(constants.Constants(), dynamicMemoryOctets, typeInformationOctets)
	if packedErr != nil {
		return packedErr
	}

	outputFilename := path.Join(absoluteOutputDirectory, fmt.Sprintf("%s.swamp-pack", packageSubDirectory))

	if err := ioutil.WriteFile(outputFilename, packed, 0o644); err != nil {
		return decorated.NewInternalError(err)
	}

	log.Printf("wrote output file '%v'", outputFilename)

	return nil
}

func (g *Generator) GenerateModule(module *decorated.Module,
	resourceNameLookup resourceid.ResourceNameLookup, verboseFlag verbosity.Verbosity) error {
	moduleContext := NewContext(g.packageConstants, "root")

	var functionConstants []*assembler_sp.Constant

	for _, named := range module.LocalDefinitions().Definitions() {
		unknownType := named.Expression()
		_, isConstant := unknownType.(*decorated.Constant)
		if isConstant {
			continue
		}
		fullyQualifiedName := module.FullyQualifiedName(named.Identifier())
		preparedFuncConstant := moduleContext.Constants().FindFunction(assembler_sp.VariableName(fullyQualifiedName.String()))
		if preparedFuncConstant == nil {
			panic(fmt.Errorf("could not find function that should have been prepared %v", fullyQualifiedName))
		}
		functionConstants = append(functionConstants, preparedFuncConstant)
		maybeFunction, _ := unknownType.(*decorated.FunctionValue)
		if maybeFunction != nil {
			if maybeFunction.Annotation().Annotation().IsSomeKindOfExternal() {
				continue
			}
			if verboseFlag >= verbosity.Mid {
				log.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
			}

			rootContext := moduleContext.MakeFunctionContext(maybeFunction, fullyQualifiedName.String())

			if maybeFunction.Annotation().Annotation().IsSomeKindOfExternal() {
				continue
			}
			generatedFunctionInfo, genFuncErr := generateFunction(fullyQualifiedName, maybeFunction,
				rootContext, g.lookup, resourceNameLookup, g.fileUrlCache, verboseFlag)
			if genFuncErr != nil {
				return genFuncErr
			}

			if generatedFunctionInfo == nil {
				panic(fmt.Sprintf("problem %v\n", maybeFunction))
			}

			if verboseFlag >= verbosity.High {
				log.Printf("---------- generated code for '%v'", fullyQualifiedName.String())
				rootContext.scopeVariables.DebugOutput(0)
			}

			moduleContext.Constants().DefineFunctionOpcodes(preparedFuncConstant, generatedFunctionInfo.opcodes)

			debugLinesOctets, debugLinesErr := opcode_sp.SerializeDebugLines(generatedFunctionInfo.debugLines)
			if debugLinesErr != nil {
				return debugLinesErr
			}

			moduleContext.Constants().DefineFunctionDebugLines(preparedFuncConstant, uint(len(generatedFunctionInfo.debugLines)), debugLinesOctets)

			generatedFunctionInfo.debugVariables = assembler_sp.GenerateVariablesWithScope(rootContext.scopeVariables, 1)
			if verboseFlag >= verbosity.High {
				assembler_sp.VariableInfosDebugOutput(generatedFunctionInfo.debugVariables)
			}
			moduleContext.Constants().DefineFunctionDebugScopes(preparedFuncConstant, generatedFunctionInfo.debugVariables)

		} else {
			maybeConstant, _ := unknownType.(*decorated.Constant)
			if maybeConstant != nil {
				if verboseFlag >= verbosity.Mid {
					log.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
				}
			} else {
				return fmt.Errorf("generate: unknown type %T", unknownType)
			}
		}
	}

	return nil
}
