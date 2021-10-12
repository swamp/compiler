package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func generateLogical(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.LogicalOperator, genContext *generateContext) error {
	leftErr := generateExpression(code, target, operator.Left(), genContext)
	if leftErr != nil {
		return leftErr
	}

	codeAlternative := assembler_sp.NewCode()
	rightErr := generateExpression(codeAlternative, target, operator.Right(), genContext)
	if rightErr != nil {
		return rightErr
	}
	afterLabel := codeAlternative.Label("", "after-alternative")

	if operator.OperatorType() == decorated.LogicalAnd {
		code.BranchFalse(targetToSourceStackPosRange(target).Pos, afterLabel)
	} else if operator.OperatorType() == decorated.LogicalOr {
		code.BranchTrue(targetToSourceStackPosRange(target).Pos, afterLabel)
	}
	code.Copy(codeAlternative)

	return nil
}

func handleLogical(code *assembler_sp.Code,
	logical *decorated.LogicalOperator, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := genContext.context.stackMemory.Allocate(uint(dectype.SizeofSwampBool), uint32(dectype.SizeofSwampBool), "logicalOperator target")
	if err := generateLogical(code, posRange, logical, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
