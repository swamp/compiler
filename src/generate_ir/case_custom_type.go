package generate_ir

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func GetIrType(p dtype.Type) types.Type {
	return types.I32
}

type customTypeCaseConsequence struct {
	Block    *ir.Block
	Value    value.Value
	Incoming *ir.Incoming
}

func generateCaseCustomType(caseExpr *decorated.CaseCustomType, genContext *generateContext) (value.Value, error) {
	testVar, testErr := generateExpression(caseExpr.Test(), false, genContext)
	if testErr != nil {
		return nil, testErr
	}

	var consequenceBlocks []customTypeCaseConsequence

	for _, consequence := range caseExpr.Consequences() {
		caseBlock := ir.NewBlock("case")
		customTypeCaseConsequence := customTypeCaseConsequence{Block: caseBlock}
		consequenceBlocks = append(consequenceBlocks, customTypeCaseConsequence)

		fields := consequence.VariantReference().CustomTypeVariant().Fields()
		for index, param := range consequence.Parameters() {
			field := fields[index]

			consequenceParameter := ir.NewGetElementPtr(GetIrType(field.Type()), testVar, constant.NewInt(types.I32, int64(index)))
			consequenceParameter.LocalIdent.SetName(param.Identifier().Name())
		}

		consequenceValue, caseExprErr := generateExpression(consequence.Expression(), false, genContext)
		if caseExprErr != nil {
			return nil, caseExprErr
		}
		customTypeCaseConsequence.Value = consequenceValue
		customTypeCaseConsequence.Incoming = ir.NewIncoming(customTypeCaseConsequence.Value, customTypeCaseConsequence.Block)
	}

	defaultCase := customTypeCaseConsequence{Block: nil}
	var defaultIncoming *ir.Incoming
	if caseExpr.DefaultCase() != nil {
		defaultBlock := ir.NewBlock("default")
		newContext := genContext
		newContext.block = defaultBlock
		defaultValue, testErr := generateExpression(caseExpr.DefaultCase(), false, genContext)
		if testErr != nil {
			return nil, testErr
		}
		defaultCase.Block = defaultBlock
		defaultIncoming = ir.NewIncoming(defaultValue, defaultBlock)
	}

	cases := make([]*ir.Case, len(consequenceBlocks))
	for index, consequenceBlock := range consequenceBlocks {
		cases[index] = ir.NewCase(constant.NewInt(types.I32, int64(index)), consequenceBlock.Block)
	}

	genContext.block.NewSwitch(testVar, defaultCase.Block, cases...)

	phiInstruction := genContext.block.NewPhi(defaultIncoming)

	return phiInstruction, nil
}
