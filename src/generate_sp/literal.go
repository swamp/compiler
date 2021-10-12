package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/instruction_sp"
)

func generateStringLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.StringLiteral,
	constants *assembler_sp.PackageConstants) error {
	if target.Size != assembler_sp.StackRange(dectype.Sizeof64BitPointer) {
		panic("wrong size")
	}

	constant := constants.AllocateStringConstant(str.Value())
	code.LoadZeroMemoryPointer(target.Pos, constant.PosRange().Position)
	return nil
}

func generateCharacterLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.CharacterLiteral) error {
	if target.Size != assembler_sp.StackRange(dectype.SizeofSwampInt) {
		panic("wrong size")
	}

	code.LoadRune(target.Pos, instruction_sp.ShortRune(str.Value()))
	return nil
}

func generateTypeIdLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, typeId *decorated.TypeIdLiteral,
	genContext *generateContext) error {
	return nil
}

func generateIntLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, integer *decorated.IntegerLiteral) error {
	if target.Size != assembler_sp.StackRange(dectype.SizeofSwampInt) {
		panic("wrong size")
	}

	code.LoadInteger(target.Pos, integer.Value())
	return nil
}

func generateFixedLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, fixed *decorated.FixedLiteral) error {
	if target.Size != assembler_sp.StackRange(dectype.SizeofSwampInt) {
		panic("wrong size")
	}

	code.LoadInteger(target.Pos, fixed.Value())
	return nil
}

func generateResourceNameLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	resourceName *decorated.ResourceNameLiteral, context *assembler_sp.PackageConstants) error {
	if target.Size != assembler_sp.StackRange(dectype.SizeofSwampInt) {
		panic("wrong size")
	}
	constant := context.AllocateResourceNameConstant(resourceName.Value())
	code.LoadInteger(target.Pos, int32(constant.ResourceID()))
	return nil
}

func generateBoolLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	boolLiteral *decorated.BooleanLiteral) error {
	if target.Size != assembler_sp.StackRange(dectype.SizeofSwampBool) {
		panic("wrong size")
	}
	code.LoadBool(target.Pos, boolLiteral.Value())
	return nil
}
