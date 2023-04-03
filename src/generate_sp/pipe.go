/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generatePipeLeft(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.PipeLeftOperator, genContext *generateContext) error {
	//leftErr := generateExpression(code, target, operator.GenerateLeft(), false, genContext)
	//if leftErr != nil {
	//	return leftErr
	//}
	return nil
}

func generatePipeRight(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.PipeRightOperator, genContext *generateContext) error {
	//leftErr := generateExpression(code, target, operator.GenerateRight(), false, genContext)
	//if leftErr != nil {
	//	return leftErr
	//}
	return nil
}

func handlePipeRight(code *assembler_sp.Code, operator *decorated.PipeRightOperator, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := allocMemoryForType(genContext.context.stackMemory, operator.Type(), "pipeRight")

	if err := generatePipeRight(code, posRange, operator, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}

func handlePipeLeft(code *assembler_sp.Code, operator *decorated.PipeLeftOperator, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := allocMemoryForType(genContext.context.stackMemory, operator.Type(), "pipeLeft")

	if err := generatePipeLeft(code, posRange, operator, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
