package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/opcodes/instruction_sp"
)

func arithmeticToBinaryOperatorType(operatorType decorated.ArithmeticOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.ArithmeticPlus:
		return instruction_sp.BinaryOperatorArithmeticIntPlus
	case decorated.ArithmeticCons:
		panic("cons not handled here")
	case decorated.ArithmeticMinus:
		return instruction_sp.BinaryOperatorArithmeticIntMinus
	case decorated.ArithmeticMultiply:
		return instruction_sp.BinaryOperatorArithmeticIntMultiply
	case decorated.ArithmeticDivide:
		return instruction_sp.BinaryOperatorArithmeticIntDivide
	case decorated.ArithmeticRemainder:
		return instruction_sp.BinaryOperatorArithmeticIntRemainder
	case decorated.ArithmeticAppend:
		return instruction_sp.BinaryOperatorArithmeticListAppend
	case decorated.ArithmeticFixedMultiply:
		return instruction_sp.BinaryOperatorArithmeticFixedMultiply
	case decorated.ArithmeticFixedDivide:
		return instruction_sp.BinaryOperatorArithmeticFixedDivide
	}

	panic("unknown binary operator")
}

func generateArithmetic(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ArithmeticOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "arith-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "arit-right")
	if rightErr != nil {
		return rightErr
	}

	opcodeBinaryOperator := arithmeticToBinaryOperatorType(operator.OperatorType())

	filePosition := genContext.toFilePosition(operator.FetchPositionLength())
	code.IntBinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator, filePosition)

	return nil
}
