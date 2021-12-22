package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	opcode_sp_type "github.com/swamp/opcodes/type"
)

func generateLogical(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.LogicalOperator, genContext *generateContext) error {
	leftErr := generateExpression(code, target, operator.Left(), false, genContext)
	if leftErr != nil {
		return leftErr
	}

	codeAlternative := assembler_sp.NewCode()
	rightErr := generateExpression(codeAlternative, target, operator.Right(), false, genContext)
	if rightErr != nil {
		return rightErr
	}
	afterLabel := codeAlternative.Label("", "after-alternative")

	filePosition := genContext.toFilePosition(operator.FetchPositionLength())
	if operator.OperatorType() == decorated.LogicalAnd {
		code.BranchFalse(targetToSourceStackPosRange(target).Pos, afterLabel, filePosition)
	} else if operator.OperatorType() == decorated.LogicalOr {
		code.BranchTrue(targetToSourceStackPosRange(target).Pos, afterLabel, filePosition)
	}
	code.Copy(codeAlternative)

	return nil
}

func handleLogical(code *assembler_sp.Code,
	logical *decorated.LogicalOperator, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := genContext.context.stackMemory.Allocate(uint(opcode_sp_type.SizeofSwampBool), uint32(opcode_sp_type.SizeofSwampBool), "logicalOperator target")
	if err := generateLogical(code, posRange, logical, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
