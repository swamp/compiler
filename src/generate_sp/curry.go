package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func generateCurry(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, call *decorated.CurryFunction,
	genContext *generateContext) error {
	if decorated.TypeIsTemplateHasLocalTypes(call.FunctionAtom()) {
		panic(fmt.Errorf("we can not call functions that has local types %v", call.FunctionAtom()))
	}

	beforePos := genContext.context.stackMemory.Tell()

	functionRegister, functionGenErr := generateExpressionWithSourceVar(code,
		call.FunctionValue(), genContext, "functioncall")
	if functionGenErr != nil {
		return functionGenErr
	}

	indexIntoTypeInformationChunk, lookupErr := genContext.lookup.Lookup(call.Type())
	if lookupErr != nil {
		return lookupErr
	}

	invokedReturnType := dectype.UnaliasWithResolveInvoker(call.FunctionAtom().ReturnType())

	genContext.context.stackMemory.AlignUpForMax()

	allocMemoryForType(genContext.context.stackMemory, invokedReturnType, "curry return")

	arguments := make([]assembler_sp.TargetStackPosRange, len(call.ArgumentsToSave()))
	for index, arg := range call.ArgumentsToSave() {
		arguments[index] = allocMemoryForType(genContext.context.stackMemory, arg.Type(), fmt.Sprintf("arg %d", index))
	}

	_, firstAlign := dectype.GetMemorySizeAndAlignment(call.ArgumentsToSave()[0].Type())

	for index, arg := range call.ArgumentsToSave() {
		argReg := arguments[index]
		argRegErr := generateExpression(code, argReg, arg, genContext)
		if argRegErr != nil {
			return argRegErr
		}
	}

	lastArgument := arguments[len(arguments)-1]
	completeArgumentRange := assembler_sp.SourceStackPosRange{
		Pos:  assembler_sp.SourceStackPos(arguments[0].Pos),
		Size: assembler_sp.SourceStackRange((uint(lastArgument.Pos) + uint(lastArgument.Size)) - uint(arguments[0].Pos)),
	}

	code.Curry(target.Pos, uint16(indexIntoTypeInformationChunk), assembler_sp.MemoryAlign(firstAlign), functionRegister.Pos, completeArgumentRange)

	genContext.context.stackMemory.Set(beforePos)

	return nil
}

func handleCurry(code *assembler_sp.Code, call *decorated.CurryFunction,
	genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	targetPosRange := genContext.context.stackMemory.Allocate(uint(dectype.Sizeof64BitPointer), uint32(dectype.Alignof64BitPointer), "")

	if err := generateCurry(code, targetPosRange, call, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(targetPosRange), nil
}
