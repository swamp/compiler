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
	name      *decorated.FullyQualifiedPackageVariableName
	signature swamppack.TypeRef
	opcodes   []byte
}

type ExternalFunction struct {
	name           *decorated.FullyQualifiedPackageVariableName
	signature      swamppack.TypeRef
	parameterCount uint
}

func NewFunction(fullyQualifiedName *decorated.FullyQualifiedPackageVariableName, signature swamppack.TypeRef,
	opcodes []byte) *Function {
	f := &Function{
		name: fullyQualifiedName, signature: signature, opcodes: opcodes,
	}

	return f
}

func NewExternalFunction(fullyQualifiedName *decorated.FullyQualifiedPackageVariableName,
	signature swamppack.TypeRef, parameterCount uint) *ExternalFunction {
	f := &ExternalFunction{name: fullyQualifiedName, signature: signature, parameterCount: parameterCount}

	return f
}

func (f *Function) String() string {
	return fmt.Sprintf("[function %v %v %v %v]", f.name, f.signature)
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

func getMemorySizeAndAlignment(p dtype.Type) (uint, uint32) {
	unaliased := dectype.Unalias(p)
	switch t := unaliased.(type) {
	case *dectype.RecordAtom:
	}
}

func allocMemoryForType(stackMemory *assembler_sp.StackMemoryMapper, typeToAlloc dtype.Type,
	debugString string) assembler_sp.TargetStackPosRange {
	memorySize, alignment := getMemorySizeAndAlignment(typeToAlloc)
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

func constantToSourceStackPosRange(code *assembler_sp.Code, stackMemory *assembler_sp.StackMemoryMapper, constant *assembler_sp.Constant) (assembler_sp.SourceStackPosRange, error) {
	functionPointer := stackMemory.Allocate(PointerSize, PointerAlign, "functionReference")
	code.LoadZeroMemoryPointer(functionPointer.Pos, constant.PosRange().Position)

	return targetToSourceStackPosRange(functionPointer), nil
}

const (
	SizeofSwampInt  = 4
	SizeofSwampRune = 2
	SizeofSwampBool = 1
)

func (g *Generator) GenerateAllLocalDefinedFunctions(module *decorated.Module, definitions *decorator.VariableContext,
	lookup typeinfo.TypeLookup, verboseFlag verbosity.Verbosity) ([]*Function, error) {
	var functionConstants []*Function

	for _, named := range module.LocalDefinitions().Definitions() {
		unknownType := named.Expression()

		fullyQualifiedName := module.FullyQualifiedName(named.Identifier())

		maybeFunction, _ := unknownType.(*decorated.FunctionValue)
		if maybeFunction != nil {
			if verboseFlag >= verbosity.Mid {
				fmt.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
			}

			rootContext := assembler_sp.NewFunctionRootContext()

			functionConstant, genFuncErr := generateFunction(fullyQualifiedName, maybeFunction,
				rootContext, definitions, lookup, verboseFlag)
			if genFuncErr != nil {
				return nil, genFuncErr
			}

			if functionConstant == nil {
				panic(fmt.Sprintf("problem %v\n", maybeFunction))
			}

			functionConstants = append(functionConstants, functionConstant)
		} else {
			maybeConstant, _ := unknownType.(*decorated.Constant)
			if maybeConstant != nil {
				if verboseFlag >= verbosity.Mid {
					fmt.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
				}
			} else {
				return nil, fmt.Errorf("generate: unknown type %T", unknownType)
			}
		}
	}

	return functionConstants, nil
}
