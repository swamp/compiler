package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
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
	afterLabel := codeAlternative.Label(nil, "after-alternative")

	if operator.OperatorType() == decorated.LogicalAnd {
		code.BranchFalse(targetToSourceStackPosRange(target).Pos, afterLabel)
	} else if operator.OperatorType() == decorated.LogicalOr {
		code.BranchTrue(targetToSourceStackPosRange(target).Pos, afterLabel)
	}
	code.Copy(codeAlternative)

	return nil
}
