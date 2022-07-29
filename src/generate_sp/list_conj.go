package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	opcode_sp_type "github.com/swamp/opcodes/type"
)

func generateListCons(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ConsOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "cons-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "cons-right")
	if rightErr != nil {
		return rightErr
	}

	itemSize, itemAlign := dectype.GetMemorySizeAndAlignment(operator.Left().Type())

	filePosition := genContext.toFilePosition(operator.FetchPositionLength())
	code.ListConj(target.Pos, leftVar.Pos, assembler_sp.StackItemSize(itemSize), opcode_sp_type.MemoryAlign(itemAlign), rightVar.Pos, filePosition)

	return nil
}

func handleListCons(code *assembler_sp.Code, operator *decorated.ConsOperator, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := genContext.context.stackMemory.Allocate(uint(opcode_sp_type.Sizeof64BitPointer), uint32(opcode_sp_type.Alignof64BitPointer), "list struct pointer")
	if err := generateListCons(code, posRange, operator, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
