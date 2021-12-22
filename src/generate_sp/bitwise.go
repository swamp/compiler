package generate_sp

import (
	"fmt"

	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/opcodes/instruction_sp"
	opcode_sp_type "github.com/swamp/opcodes/type"
)

func bitwiseToBinaryOperatorType(operatorType decorated.BitwiseOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.BitwiseAnd:
		return instruction_sp.BinaryOperatorBitwiseIntAnd
	case decorated.BitwiseOr:
		return instruction_sp.BinaryOperatorBitwiseIntOr
	case decorated.BitwiseXor:
		return instruction_sp.BinaryOperatorBitwiseIntXor
	case decorated.BitwiseShiftLeft:
		return instruction_sp.BinaryOperatorBitwiseShiftLeft
	case decorated.BitwiseShiftRight:
		return instruction_sp.BinaryOperatorBitwiseShiftRight
	default:
		panic(fmt.Errorf("not a binary operator %v", operatorType))
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

	filePosition := genContext.toFilePosition(operator.FetchPositionLength())
	code.IntBinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator, filePosition)

	return nil
}

func handleBitwise(code *assembler_sp.Code,
	bitwise *decorated.BitwiseOperator, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := genContext.context.stackMemory.Allocate(uint(opcode_sp_type.SizeofSwampInt), uint32(opcode_sp_type.SizeofSwampInt), "bitwiseOperator target")
	if err := generateBitwise(code, posRange, bitwise, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
