package generate_ir

import (
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"log"
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

	for consequenceIndex, consequence := range caseExpr.Consequences() {
		customTypeCaseConsequence := customTypeCaseConsequence{}

		fields := consequence.VariantReference().CustomTypeVariant().Fields()
		for index, param := range consequence.Parameters() {
			field := fields[index]

			consequenceParameter := ir.NewGetElementPtr(GetIrType(field.Type()), testVar, constant.NewInt(types.I32, int64(index)))
			consequenceParameter.LocalIdent.SetName(param.Identifier().Name())
		}

		consequenceContext := genContext.NewBlock(fmt.Sprintf("consequnce%d", consequenceIndex))
		consequenceValue, caseExprErr := generateExpression(consequence.Expression(), false, consequenceContext)
		if caseExprErr != nil {
			return nil, caseExprErr
		}
		if consequenceValue == nil {
			panic(fmt.Errorf("this should have been reported as error"))
		}
		customTypeCaseConsequence.Value = consequenceValue
		customTypeCaseConsequence.Block = consequenceContext.block
		customTypeCaseConsequence.Incoming = ir.NewIncoming(customTypeCaseConsequence.Value, customTypeCaseConsequence.Block)
		consequenceBlocks = append(consequenceBlocks, customTypeCaseConsequence)
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

	var incomingInfos []*ir.Incoming

	if defaultIncoming != nil {
		incomingInfos = append(incomingInfos, defaultIncoming)
	}

	for _, consequenceBlock := range consequenceBlocks {
		incoming := ir.NewIncoming(consequenceBlock.Value, consequenceBlock.Block)
		if consequenceBlock.Value == nil {
			panic("consequenceBlock.Value is nil")
		}
		if consequenceBlock.Block == nil {
			panic("consequenceBlock.Block is nil")
		}
		incomingInfos = append(incomingInfos, incoming)
	}

	log.Printf("incoming infos %v", incomingInfos)

	phiInstruction := genContext.block.NewPhi(incomingInfos...)

	return phiInstruction, nil
}
