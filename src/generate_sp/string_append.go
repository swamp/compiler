/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateStringAppend(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ArithmeticOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "string-append-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "string-append-right")
	if rightErr != nil {
		return rightErr
	}

	filePosition := genContext.toFilePosition(operator.FetchPositionLength())

	code.StringAppend(target.Pos, leftVar.Pos, rightVar.Pos, filePosition)

	return nil
}
