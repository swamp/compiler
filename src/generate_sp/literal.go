package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateStringLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.StringLiteral,
	constants *assembler_sp.Constants) error {
	constant := constants.AllocateStringConstant(str.Value())
	code.LoadZeroMemoryPointer(target.Pos, constant.PosRange().Position)
	return nil
}

func generateCharacterLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.CharacterLiteral) error {
	code.LoadRune(target.Pos, uint8(str.Value()))
	return nil
}

func generateTypeIdLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, typeId *decorated.TypeIdLiteral,
	genContext *generateContext) error {
	return nil
}

func generateIntLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, integer *decorated.IntegerLiteral) error {
	code.LoadInteger(target.Pos, integer.Value())
	return nil
}

func generateFixedLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, fixed *decorated.FixedLiteral) error {
	code.LoadInteger(target.Pos, fixed.Value())
	return nil
}

func generateResourceNameLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	resourceName *decorated.ResourceNameLiteral, context *assembler_sp.Constants) error {
	return nil
}

func generateBoolLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	boolLiteral *decorated.BooleanLiteral) error {
	code.LoadBool(target.Pos, boolLiteral.Value())
	return nil
}
