package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/instruction_sp"
)

func generateStringLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.StringLiteral,
	constants *assembler_sp.PackageConstants) error {
	constant := constants.AllocateStringConstant(str.Value())
	code.LoadZeroMemoryPointer(target.Pos, constant.PosRange().Position)
	return nil
}

func generateCharacterLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.CharacterLiteral) error {
	code.LoadRune(target.Pos, instruction_sp.ShortRune(str.Value()))
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
	resourceName *decorated.ResourceNameLiteral, context *assembler_sp.PackageConstants) error {
	constant := context.AllocateResourceNameConstant(resourceName.Value())
	code.LoadInteger(target.Pos, int32(constant.ResourceID()))
	return nil
}

func generateBoolLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	boolLiteral *decorated.BooleanLiteral) error {
	code.LoadBool(target.Pos, boolLiteral.Value())
	return nil
}
