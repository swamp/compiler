/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/opcodes/opcode_sp"
)

func generateGuard(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, guardExpr *decorated.Guard,
	genContext *generateContext) error {
	type codeItem struct {
		ConditionVariable     assembler_sp.SourceStackPosRange
		ConditionCode         *assembler_sp.Code
		ConsequenceCode       *assembler_sp.Code
		EndOfConsequenceLabel *assembler_sp.Label
	}

	defaultCode := assembler_sp.NewCode()
	defaultContext := genContext.MakeScopeContext("guard")

	altErr := generateExpression(defaultCode, target, guardExpr.DefaultGuard().Expression(), false, defaultContext)
	if altErr != nil {
		return altErr
	}

	endLabel := defaultCode.Label("", "guard-end")

	var codeItems []codeItem

	for _, item := range guardExpr.Items() {
		conditionCode := assembler_sp.NewCode()
		conditionCodeContext := genContext.MakeScopeContext("guard condition")

		conditionVar, testErr := generateExpressionWithSourceVar(conditionCode,
			item.Condition(), conditionCodeContext, "guard-condition")
		if testErr != nil {
			return testErr
		}

		consequenceCode := assembler_sp.NewCode()
		consequenceContext := genContext.MakeScopeContext("guard consequence")

		consErr := generateExpression(consequenceCode, target, item.Expression(), false, consequenceContext)
		if consErr != nil {
			return consErr
		}

		consequenceCode.Jump(endLabel, opcode_sp.FilePosition{})
		endOfConsequenceLabel := consequenceCode.Label("", "guard-end")

		// consequenceContext.context.Free()

		codeItem := codeItem{
			ConditionCode: conditionCode, ConditionVariable: conditionVar, ConsequenceCode: consequenceCode,
			EndOfConsequenceLabel: endOfConsequenceLabel,
		}

		codeItems = append(codeItems, codeItem)
	}

	for _, codeItem := range codeItems {
		code.Copy(codeItem.ConditionCode)
		code.BranchFalse(codeItem.ConditionVariable.Pos, codeItem.EndOfConsequenceLabel, opcode_sp.FilePosition{})

		code.Copy(codeItem.ConsequenceCode)
	}

	code.Copy(defaultCode)

	return nil
}

func handleGuard(code *assembler_sp.Code, guardExpr *decorated.Guard,
	genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	targetPosRange := allocMemoryForType(genContext.context.stackMemory, guardExpr.Type(), "guard target")

	if err := generateGuard(code, targetPosRange, guardExpr, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(targetPosRange), nil
}
