/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"fmt"
	"log"

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

const (
	Sizeof64BitPointer  uint   = 8
	Alignof64BitPointer uint32 = 8
)

func GetMemorySizeAndAlignment(p dtype.Type) (uint, uint32) {
	if p == nil {
		panic(fmt.Errorf("nil is not allowed"))
	}
	unaliased := dectype.Unalias(p)
	switch t := unaliased.(type) {
	case *dectype.RecordAtom:
		return t.MemorySize(), t.MemoryAlignment()
	case *dectype.PrimitiveAtom:
		{
			name := t.PrimitiveName().Name()
			switch name {
			case "List":
				{
					return Sizeof64BitPointer, Alignof64BitPointer
				}
			case "Bool":
				return SizeofSwampBool, AlignOfSwampBool
			case "Int":
				return SizeofSwampInt, AlignOfSwampInt
			case "Fixed":
				return SizeofSwampInt, AlignOfSwampInt
			case "Char":
				return SizeofSwampInt, AlignOfSwampInt
			case "String":
				return Sizeof64BitPointer, Alignof64BitPointer
			}
			panic(fmt.Errorf("do not know primitive atom of %v %T", p, unaliased))
		}
	case *dectype.InvokerType:
		switch it := t.TypeGenerator().(type) {
		case *dectype.PrimitiveTypeReference:
			return Sizeof64BitPointer, Alignof64BitPointer
		case *dectype.CustomTypeReference:
			log.Printf("this is : %T %v (%v)", it.Type(), it.HumanReadable(), t.TypeGenerator().HumanReadable())
			return GetMemorySizeAndAlignment(it.Type())
		}

	case *dectype.CustomTypeAtom:
		return t.MemorySize(), t.MemoryAlignment()
	case *dectype.FunctionAtom:
		return Sizeof64BitPointer, Alignof64BitPointer
	default:
		panic(fmt.Errorf("do not know memory size of %v %T", p, unaliased))
	}
	panic(fmt.Errorf("do not know memory size of %v %T", p, unaliased))
}

func allocMemoryForType(stackMemory *assembler_sp.StackMemoryMapper, typeToAlloc dtype.Type,
	debugString string) assembler_sp.TargetStackPosRange {
	memorySize, alignment := GetMemorySizeAndAlignment(typeToAlloc)
	return stackMemory.Allocate(memorySize, alignment, debugString)
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

const (
	SizeofSwampInt  uint = 4
	SizeofSwampRune uint = 4
	SizeofSwampBool uint = 1

	AlignOfSwampBool uint32 = uint32(SizeofSwampBool)
	AlignOfSwampRune uint32 = uint32(SizeofSwampRune)
	AlignOfSwampInt  uint32 = uint32(SizeofSwampInt)
)

func (g *Generator) GenerateAllLocalDefinedFunctions(module *decorated.Module, definitions *decorator.VariableContext,
	lookup typeinfo.TypeLookup, packageConstants *assembler_sp.PackageConstants, verboseFlag verbosity.Verbosity) (*assembler_sp.PackageConstants, []*assembler_sp.Constant, error) {
	moduleContext := NewContext(packageConstants)

	var functionConstants []*assembler_sp.Constant

	for _, named := range module.LocalDefinitions().Definitions() {
		unknownType := named.Expression()
		fullyQualifiedName := module.FullyQualifiedName(named.Identifier())
		preparedFuncConstant := moduleContext.Constants().FindFunction(assembler_sp.VariableName(fullyQualifiedName.String()))
		functionConstants = append(functionConstants, preparedFuncConstant)
		maybeFunction, _ := unknownType.(*decorated.FunctionValue)
		if maybeFunction != nil {
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

			rootContext.constants.DynamicMemory().DebugOutput()
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
