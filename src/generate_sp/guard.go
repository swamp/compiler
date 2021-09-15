package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateGuard(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, guardExpr *decorated.Guard, genContext *generateContext) error {
	type codeItem struct {
		ConditionVariable     assembler_sp.SourceStackPosRange
		ConditionCode         *assembler_sp.Code
		ConsequenceCode       *assembler_sp.Code
		EndOfConsequenceLabel *assembler_sp.Label
	}

	defaultCode := assembler_sp.NewCode()
	// defaultLabel := defaultCode.Label(nil, "guard-default")
	defaultContext := genContext.MakeScopeContext()

	altErr := generateExpression(defaultCode, target, guardExpr.DefaultGuard().Expression(), defaultContext)
	if altErr != nil {
		return altErr
	}

	endLabel := defaultCode.Label(nil, "guard-end")

	var codeItems []codeItem

	for _, item := range guardExpr.Items() {
		conditionCode := assembler_sp.NewCode()
		conditionCodeContext := genContext.MakeScopeContext()

		conditionVar, testErr := generateExpressionWithSourceVar(conditionCode,
			item.Condition(), conditionCodeContext, "guard-condition")
		if testErr != nil {
			return testErr
		}

		consequenceCode := assembler_sp.NewCode()
		consequenceContext := genContext.MakeScopeContext()

		consErr := generateExpression(consequenceCode, target, item.Expression(), consequenceContext)
		if consErr != nil {
			return consErr
		}

		consequenceCode.Jump(endLabel)
		endOfConsequenceLabel := consequenceCode.Label(nil, "guard-end")

		// consequenceContext.context.Free()

		codeItem := codeItem{
			ConditionCode: conditionCode, ConditionVariable: conditionVar, ConsequenceCode: consequenceCode,
			EndOfConsequenceLabel: endOfConsequenceLabel,
		}

		codeItems = append(codeItems, codeItem)
	}

	for _, codeItem := range codeItems {
		code.Copy(codeItem.ConditionCode)
		code.BranchFalse(codeItem.ConditionVariable.Pos, codeItem.EndOfConsequenceLabel)

		code.Copy(codeItem.ConsequenceCode)
	}

	code.Copy(defaultCode)

	return nil
}