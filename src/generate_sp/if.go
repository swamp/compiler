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

func generateIf(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, ifExpr *decorated.If, genContext *generateContext) error {
	conditionVar, testErr := generateExpressionWithSourceVar(code, ifExpr.Condition(), genContext, "if-condition")
	if testErr != nil {
		return testErr
	}

	consequenceCode := assembler_sp.NewCode()
	consequenceContext2 := genContext.MakeScopeContext("if consequenceContext")

	consErr := generateExpression(consequenceCode, target, ifExpr.Consequence(), true, consequenceContext2)
	if consErr != nil {
		return consErr
	}

	// consequenceContext2.context.Free()

	alternativeCode := assembler_sp.NewCode()
	alternativeLabel := alternativeCode.Label("", "if-alternative")
	alternativeContext2 := genContext.MakeScopeContext("if alternative")

	altErr := generateExpression(alternativeCode, target, ifExpr.Alternative(), true, alternativeContext2)
	if altErr != nil {
		return altErr
	}

	endLabel := alternativeCode.Label("", "if-end")

	// alternativeContext2.context.Free()

	code.BranchFalse(conditionVar.Pos, alternativeLabel, opcode_sp.FilePosition{})

	consequenceCode.Jump(endLabel, opcode_sp.FilePosition{})
	code.Copy(consequenceCode)
	code.Copy(alternativeCode)

	return nil
}

func handleIf(code *assembler_sp.Code, guardExpr *decorated.If,
	genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	targetPosRange := allocMemoryForType(genContext.context.stackMemory, guardExpr.Type(), "if target")

	if err := generateIf(code, targetPosRange, guardExpr, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(targetPosRange), nil
}
