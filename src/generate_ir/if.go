package generate_ir

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"log"
)

func generateIf(ifExpr *decorated.If, isLeafNode bool, genContext *generateContext) (value.Value, error) {
	conditionVar, testErr := generateExpression(ifExpr.Condition(), false, genContext)
	if testErr != nil {
		return nil, testErr
	}

	consequenceContext := genContext.NewBlock("consequence")
	consequenceValue, consErr := generateExpression(ifExpr.Consequence(), isLeafNode, consequenceContext)
	if consErr != nil {
		return nil, consErr
	}

	consequenceIncoming := ir.NewIncoming(consequenceValue, consequenceContext.block)

	alternativeContext := genContext.NewBlock("alternative")
	alternativeValue, altErr := generateExpression(ifExpr.Alternative(), isLeafNode, alternativeContext)
	if altErr != nil {
		return nil, altErr
	}

	alternativeIncoming := ir.NewIncoming(alternativeValue, alternativeContext.block)

	genContext.block.NewCondBr(conditionVar, consequenceContext.block, alternativeContext.block)
	phi := ir.NewPhi(consequenceIncoming, alternativeIncoming)

	log.Printf("\nif:%v\n", phi)

	return phi, nil
}
