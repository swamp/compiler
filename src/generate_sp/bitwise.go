package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/opcodes/instruction_sp"
)

func bitwiseToBinaryOperatorType(operatorType decorated.BitwiseOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.BitwiseAnd:
		return instruction_sp.BinaryOperatorBitwiseIntAnd
	case decorated.BitwiseOr:
		return instruction_sp.BinaryOperatorBitwiseIntOr
	case decorated.BitwiseXor:
		return instruction_sp.BinaryOperatorBitwiseIntXor
	case decorated.BitwiseNot:
		return 0
		// return opcode_sp.BinaryOperatorBitwiseIntNot
	}

	return 0
}

func generateBitwise(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	operator *decorated.BitwiseOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "bitwise-right")
	if rightErr != nil {
		return rightErr
	}

	opcodeBinaryOperator := bitwiseToBinaryOperatorType(operator.OperatorType())
	code.IntBinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator)

	return nil
}
