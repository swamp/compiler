package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/instruction_sp"
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
	code.IntBinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator)

	return nil
}
