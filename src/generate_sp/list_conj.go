package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateListCons(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ConsOperator, genContext *generateContext) error {
	// leftVar := context.AllocateTempVariable("arit-left")
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "cons-left")
	if leftErr != nil {
		return leftErr
	}

	// rightVar := context.AllocateTempVariable("arit-right")
	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "cons-right")
	if rightErr != nil {
		return rightErr
	}

	code.ListConj(target.Pos, leftVar.Pos, rightVar.Pos)

	return nil
}
