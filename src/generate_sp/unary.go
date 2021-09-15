package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/instruction_sp"
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

func generateUnaryArithmetic(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ArithmeticUnaryOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := arithmeticToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target.Pos, leftVar.Pos, opcodeUnaryOperatorType)

	return nil
}
