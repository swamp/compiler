/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/instruction_sp"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
	swamppack "github.com/swamp/pack/lib"
)

type AnyPosAndRange interface {
	getPosition() uint32
	getSize() uint32
}

type Function struct {
	name                *decorated.FullyQualifiedPackageVariableName
	signature           swamppack.TypeRef
	opcodes             []byte
	debugParameterCount uint
}

type ExternalFunction struct {
	name           *decorated.FullyQualifiedPackageVariableName
	signature      swamppack.TypeRef
	parameterCount uint
}

func NewFunction(fullyQualifiedName *decorated.FullyQualifiedPackageVariableName, signature swamppack.TypeRef,
	opcodes []byte, debugParameterCount uint) *Function {
	f := &Function{
		name: fullyQualifiedName, signature: signature, opcodes: opcodes, debugParameterCount: debugParameterCount,
	}

	return f
}

func NewExternalFunction(fullyQualifiedName *decorated.FullyQualifiedPackageVariableName,
	signature swamppack.TypeRef, parameterCount uint) *ExternalFunction {
	f := &ExternalFunction{name: fullyQualifiedName, signature: signature, parameterCount: parameterCount}

	return f
}

func (f *Function) String() string {
	return fmt.Sprintf("[function %v %v]", f.name, f.signature)
}

func (f *Function) Opcodes() []byte {
	return f.opcodes
}

type Generator struct {
	code *assembler_sp.Code
}

func NewGenerator() *Generator {
	return &Generator{code: assembler_sp.NewCode()}
}

func arithmeticToUnaryOperatorType(operatorType decorated.ArithmeticUnaryOperatorType) instruction_sp.UnaryOperatorType {
	switch operatorType {
	case decorated.ArithmeticUnaryMinus:
		return instruction_sp.UnaryOperatorNegate
	}

	panic("illegal unaryoperator")
}

func generateAsm(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, asm *decorated.AsmConstant, context *Context) error {
	// compileErr := asmcompile.CompileToCodeAndContext(asm.Asm().Asm(), code, context)
	// return compileErr

	return nil
}

func generateConstant(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	constant *decorated.Constant, context *generateContext) error {
	return generateExpression(code, target, constant.Expression(), context)
}

func allocMemoryForType(stackMemory *assembler_sp.StackMemoryMapper, typeToAlloc dtype.Type,
	debugString string) assembler_sp.TargetStackPosRange {
	memorySize, alignment := dectype.GetMemorySizeAndAlignment(typeToAlloc)
	if uint(memorySize) == 0 || uint(alignment) == 0 {
		panic(fmt.Errorf("can not allocate zero memory or align for type %T %v", typeToAlloc, typeToAlloc))
	}
	return stackMemory.Allocate(uint(memorySize), uint32(alignment), debugString)
}

func generateRecurCall(code *assembler_sp.Code, call *decorated.RecurCall, genContext *generateContext) error {
	code.Recur()

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
	code.LoadZeroMemoryPointer(functionPointer.Pos, constant.PosRange().Position)

	return targetToSourceStackPosRange(functionPointer), nil
}

func (g *Generator) GenerateAllLocalDefinedFunctions(module *decorated.Module, definitions *decorator.VariableContext,
	lookup typeinfo.TypeLookup, packageConstants *assembler_sp.PackageConstants, verboseFlag verbosity.Verbosity) (*assembler_sp.PackageConstants, []*assembler_sp.Constant, error) {
	moduleContext := NewContext(packageConstants)

	var functionConstants []*assembler_sp.Constant

	for _, named := range module.LocalDefinitions().Definitions() {
		unknownType := named.Expression()
		fullyQualifiedName := module.FullyQualifiedName(named.Identifier())
		preparedFuncConstant := moduleContext.Constants().FindFunction(assembler_sp.VariableName(fullyQualifiedName.String()))
		if preparedFuncConstant == nil {
			// panic(fmt.Errorf("could not find function that should have been prepared %v", fullyQualifiedName))
			continue
		}
		functionConstants = append(functionConstants, preparedFuncConstant)
		maybeFunction, _ := unknownType.(*decorated.FunctionValue)
		if maybeFunction != nil {
			if maybeFunction.Annotation().Annotation().IsExternal() {
				continue
			}
			if verboseFlag >= verbosity.Mid {
				fmt.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
			}

			rootContext := moduleContext.MakeFunctionContext()

			if maybeFunction.Annotation().Annotation().IsExternal() {
				continue
			}
			generatedFunctionInfo, genFuncErr := generateFunction(fullyQualifiedName, maybeFunction,
				rootContext, definitions, lookup, verboseFlag)
			if genFuncErr != nil {
				return nil, nil, genFuncErr
			}

			if generatedFunctionInfo == nil {
				panic(fmt.Sprintf("problem %v\n", maybeFunction))
			}

			moduleContext.Constants().DefineFunctionOpcodes(preparedFuncConstant, generatedFunctionInfo.opcodes)

		} else {
			maybeConstant, _ := unknownType.(*decorated.Constant)
			if maybeConstant != nil {
				if verboseFlag >= verbosity.Mid {
					fmt.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
				}
			} else {
				return nil, nil, fmt.Errorf("generate: unknown type %T", unknownType)
			}
		}
	}

	return moduleContext.constants, functionConstants, nil
}
