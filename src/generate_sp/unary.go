package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/opcodes/instruction_sp"
)

func bitwiseToUnaryOperatorType(operatorType decorated.BitwiseUnaryOperatorType) instruction_sp.UnaryOperatorType {
	switch operatorType {
	case decorated.BitwiseUnaryNot:
		return instruction_sp.UnaryOperatorBitwiseNot
	}

	panic("illegal unaryoperator")
}

func generateUnaryBitwise(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.BitwiseUnaryOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := bitwiseToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target.Pos, leftVar.Pos, opcodeUnaryOperatorType)

	return nil
}

func logicalToUnaryOperatorType(operatorType decorated.LogicalUnaryOperatorType) instruction_sp.UnaryOperatorType {
	switch operatorType {
	case decorated.LogicalUnaryNot:
		return instruction_sp.UnaryOperatorNot
	}

	panic("illegal unaryoperator")
}

func generateUnaryLogical(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.LogicalUnaryOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := logicalToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target.Pos, leftVar.Pos, opcodeUnaryOperatorType)
	return nil
}

func handleUnaryLogical(code *assembler_sp.Code, operator *decorated.LogicalUnaryOperator, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	itemSize, itemAlign := dectype.GetMemorySizeAndAlignment(operator.Type())
	unaryPointer := genContext.context.stackMemory.Allocate(uint(itemSize), uint32(itemAlign), "unary")

	if err := generateUnaryLogical(code, unaryPointer, operator, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(unaryPointer), nil
}

func generateUnaryArithmetic(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ArithmeticUnaryOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := arithmeticToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target.Pos, leftVar.Pos, opcodeUnaryOperatorType)

	return nil
}
