package generate_sp

import (
	"fmt"

	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/opcodes/instruction_sp"
	opcode_sp_type "github.com/swamp/opcodes/type"
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

func booleanToBinaryEnumOperatorType(operatorType decorated.BooleanOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.BooleanEqual:
		return instruction_sp.BinaryOperatorBooleanEnumEqual
	case decorated.BooleanNotEqual:
		return instruction_sp.BinaryOperatorBooleanEnumNotEqual
	default:
		panic(fmt.Errorf("not allowed enum operator type"))
	}
}

func booleanToBinaryStringOperatorType(operatorType decorated.BooleanOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.BooleanEqual:
		return instruction_sp.BinaryOperatorBooleanStringEqual
	case decorated.BooleanNotEqual:
		return instruction_sp.BinaryOperatorBooleanStringNotEqual
	}

	return 0
}

func generateBinaryOperatorBooleanResult(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.BooleanOperator, genContext *generateContext) error {
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
	if foundPrimitive == nil {
		foundCustomType, _ := unaliasedTypeLeft.(*dectype.CustomTypeAtom)
		if foundCustomType == nil {
			panic(fmt.Errorf("not implemented binary operator boolean %v", unaliasedTypeLeft.HumanReadable()))
		} else {
			// unaliasedTypeRight := dectype.UnaliasWithResolveInvoker(operator.Right().Type())
			opcodeBinaryOperator := booleanToBinaryEnumOperatorType(operator.OperatorType())
			code.EnumBinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator)
			//			panic(fmt.Errorf("not implemented yet %v", unaliasedTypeRight))
		}
	} else if foundPrimitive.AtomName() == "String" {
		opcodeBinaryOperator := booleanToBinaryStringOperatorType(operator.OperatorType())
		code.StringBinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator)
	} else if foundPrimitive.AtomName() == "Int" || foundPrimitive.AtomName() == "Char" {
		opcodeBinaryOperator := booleanToBinaryIntOperatorType(operator.OperatorType())

		code.IntBinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator)
	} else {
		panic(fmt.Errorf("what operator is this for %v", foundPrimitive.AtomName()))
	}

	return nil
}

func handleBinaryOperatorBooleanResult(code *assembler_sp.Code, operator *decorated.BooleanOperator, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	target := genContext.context.stackMemory.Allocate(uint(opcode_sp_type.SizeofSwampBool), uint32(opcode_sp_type.AlignOfSwampBool), "booleanOperatorTarget")
	if err := generateBinaryOperatorBooleanResult(code, target, operator, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(target), nil
}
