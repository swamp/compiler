package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generatePipeLeft(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.PipeLeftOperator, genContext *generateContext) error {
	leftErr := generateExpression(code, target, operator.GenerateLeft(), genContext)
	if leftErr != nil {
		return leftErr
	}
	return nil
}

func generatePipeRight(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.PipeRightOperator, genContext *generateContext) error {
	leftErr := generateExpression(code, target, operator.GenerateRight(), genContext)
	if leftErr != nil {
		return leftErr
	}
	return nil
}
