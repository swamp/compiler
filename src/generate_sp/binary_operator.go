package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/instruction_sp"
)

func booleanToBinaryIntOperatorType(operatorType decorated.BooleanOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.BooleanEqual:
		return instruction_sp.BinaryOperatorBooleanIntEqual
	case decorated.BooleanNotEqual:
		return instruction_sp.BinaryOperatorBooleanIntNotEqual
	case decorated.BooleanLess:
		return instruction_sp.BinaryOperatorBooleanIntLess
	case decorated.BooleanLessOrEqual:
		return instruction_sp.BinaryOperatorBooleanIntLessOrEqual
	case decorated.BooleanGreater:
		return instruction_sp.BinaryOperatorBooleanIntGreater
	case decorated.BooleanGreaterOrEqual:
		return instruction_sp.BinaryOperatorBooleanIntGreaterOrEqual
	}

	return 0
}

func booleanToBinaryValueOperatorType(operatorType decorated.BooleanOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.BooleanEqual:
		return instruction_sp.BinaryOperatorBooleanValueEqual
	case decorated.BooleanNotEqual:
		return instruction_sp.BinaryOperatorBooleanValueNotEqual
	}

	return 0
}

func generateBoolean(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.BooleanOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bool-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "bool-right")
	if rightErr != nil {
		return rightErr
	}

	unaliasedTypeLeft := dectype.UnaliasWithResolveInvoker(operator.Left().Type())
	foundPrimitive, _ := unaliasedTypeLeft.(*dectype.PrimitiveAtom)

	opcodeBinaryOperator := booleanToBinaryIntOperatorType(operator.OperatorType())
	if foundPrimitive == nil || foundPrimitive.AtomName() != "Int" {
		opcodeBinaryOperator = booleanToBinaryValueOperatorType(operator.OperatorType())
	}

	code.IntBinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator)

	return nil
}
