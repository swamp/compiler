package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func handleFunctionReference(code *assembler_sp.Code,
	t *decorated.FunctionReference,
	stackMemory *assembler_sp.StackMemoryMapper,
	constants *assembler_sp.Constants) (assembler_sp.SourceStackPosRange, error) {
	ident := t.Identifier()
	functionReferenceName := assembler_sp.VariableName(ident.Name())
	foundConstant := constants.FindFunction(functionReferenceName)
	if foundConstant == nil {
		return assembler_sp.SourceStackPosRange{}, fmt.Errorf("couldn't find it %v", t)
	}

	return constantToSourceStackPosRange(code, stackMemory, foundConstant)
}

func generateFunctionReference(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	getVar *decorated.FunctionReference, context *assembler_sp.Context) error {
	varName := assembler_sp.VariableName(getVar.Identifier().Name())
	functionConstant := context.Constants().FindFunction(varName)
	code.LoadZeroMemoryPointer(target.Pos, functionConstant.PosRange().Position)
	return nil
}
