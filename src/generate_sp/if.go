package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateIf(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, ifExpr *decorated.If, genContext *generateContext) error {
	conditionVar, testErr := generateExpressionWithSourceVar(code, ifExpr.Condition(), genContext, "if-condition")
	if testErr != nil {
		return testErr
	}

	consequenceCode := assembler_sp.NewCode()
	consequenceContext2 := genContext.MakeScopeContext()

	consErr := generateExpression(consequenceCode, target, ifExpr.Consequence(), consequenceContext2)
	if consErr != nil {
		return consErr
	}

	// consequenceContext2.context.Free()

	alternativeCode := assembler_sp.NewCode()
	alternativeLabel := alternativeCode.Label("", "if-alternative")
	alternativeContext2 := genContext.MakeScopeContext()

	altErr := generateExpression(alternativeCode, target, ifExpr.Alternative(), alternativeContext2)
	if altErr != nil {
		return altErr
	}

	endLabel := alternativeCode.Label("", "if-end")

	// alternativeContext2.context.Free()

	code.BranchFalse(conditionVar.Pos, alternativeLabel)

	consequenceCode.Jump(endLabel)
	code.Copy(consequenceCode)
	code.Copy(alternativeCode)

	return nil
}
