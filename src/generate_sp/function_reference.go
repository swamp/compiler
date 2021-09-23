package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func handleFunctionReference(code *assembler_sp.Code,
	t *decorated.FunctionReference,
	stackMemory *assembler_sp.StackMemoryMapper,
	constants *assembler_sp.PackageConstants) (assembler_sp.SourceStackPosRange, error) {
	ident := t.NameReference().FullyQualifiedName()
	functionReferenceName := assembler_sp.VariableName(ident)
	foundConstant := constants.FindFunction(functionReferenceName)
	if foundConstant == nil {
		/*
			targetPosRange := stackMemory.Allocate(Sizeof64BitPointer, Alignof64BitPointer, "Hackptr")
			fake := assembler_sp.SourceDynamicMemoryPos(9494)
			code.LoadZeroMemoryPointer(targetPosRange.Pos, fake)
			return targetToSourceStackPosRange(targetPosRange), nil

		*/
		return assembler_sp.SourceStackPosRange{}, fmt.Errorf("%v couldn't find function reference '%s' %v", t.FetchPositionLength().ToReferenceString(), functionReferenceName, t)
	}

	return constantToSourceStackPosRange(code, stackMemory, foundConstant)
}

func generateFunctionReference(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	getVar *decorated.FunctionReference, constants *assembler_sp.PackageConstants) error {
	varName := assembler_sp.VariableName(getVar.Identifier().Name())
	functionConstant := constants.FindFunction(varName)
	code.LoadZeroMemoryPointer(target.Pos, functionConstant.PosRange().Position)
	return nil
}
