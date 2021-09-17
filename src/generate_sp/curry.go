package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateCurry(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, call *decorated.CurryFunction,
	genContext *generateContext) error {
	functionRegister, functionGenErr := generateExpressionWithSourceVar(code,
		call.FunctionValue(), genContext, "functioncall")
	if functionGenErr != nil {
		return functionGenErr
	}

	indexIntoTypeInformationChunk, lookupErr := genContext.lookup.Lookup(call.Type())
	if lookupErr != nil {
		return lookupErr
	}

	arguments := make([]assembler_sp.TargetStackPosRange, len(call.ArgumentsToSave()))
	for index, arg := range call.ArgumentsToSave() {
		arguments[index] = allocMemoryForType(genContext.context.stackMemory, arg.Type(), fmt.Sprintf("arg %d", index))
	}

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

	code.Curry(target.Pos, uint16(indexIntoTypeInformationChunk), functionRegister.Pos, completeArgumentRange)

	return nil
}

func handleCurry(code *assembler_sp.Code, call *decorated.CurryFunction,
	genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	targetPosRange := genContext.context.stackMemory.Allocate(Sizeof64BitPointer, Alignof64BitPointer, "")

	if err := generateCurry(code, targetPosRange, call, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(targetPosRange), nil
}
