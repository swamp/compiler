package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateListAppend(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ArithmeticOperator, genContext *generateContext) error {
	// leftVar := context.AllocateTempVariable("arit-left")
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "list-append-left")
	if leftErr != nil {
		return leftErr
	}

	// rightVar := context.AllocateTempVariable("arit-right")
	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "list-append-right")
	if rightErr != nil {
		return rightErr
	}

	filePosition := genContext.toFilePosition(operator.FetchPositionLength())
	code.ListAppend(target.Pos, leftVar.Pos, rightVar.Pos, filePosition)

	return nil
}
