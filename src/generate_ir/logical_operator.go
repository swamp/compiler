/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_ir

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateLogical(operator *decorated.LogicalOperator, genContext *generateContext) (value.Value, error) {
	leftContext := genContext.NewBlock("left")
	leftValue, leftErr := generateExpression(operator.Left(), false, leftContext)
	if leftErr != nil {
		return nil, leftErr
	}

	rightContext := genContext.NewBlock("alternative")
	rightValue, rightErr := generateExpression(operator.Right(), true, rightContext)
	if rightErr != nil {
		return nil, rightErr
	}

	falseBlock := ir.NewBlock("empty")
	falseValue := constant.NewBool(false)
	falseBlock.NewRet(falseValue)

	if operator.OperatorType() == decorated.LogicalAnd {
		genContext.block.NewCondBr(leftValue, rightContext.block, falseBlock)
	} else if operator.OperatorType() == decorated.LogicalOr {
		// inverted := ir.NewXor(leftValue, constant.NewBool(true))
		genContext.block.NewCondBr(leftValue, falseBlock, rightContext.block)
	}

	var incomingInfos []*ir.Incoming

	incomingInfos = append(incomingInfos, ir.NewIncoming(leftValue, leftContext.block))
	incomingInfos = append(incomingInfos, ir.NewIncoming(rightValue, rightContext.block))
	incomingInfos = append(incomingInfos, ir.NewIncoming(falseValue, falseBlock))

	return ir.NewPhi(incomingInfos...), nil
}
