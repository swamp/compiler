package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/opcodes/instruction_sp"
	opcode_sp_type "github.com/swamp/opcodes/type"
)

func generateStringLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.StringLiteral,
	constants *assembler_sp.PackageConstants) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.Sizeof64BitPointer) {
		panic("wrong size")
	}

	constant := constants.AllocateStringConstant(str.Value())
	code.LoadZeroMemoryPointer(target.Pos, constant.PosRange().Position)
	return nil
}

func generateCharacterLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.CharacterLiteral) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampInt) {
		panic("wrong size")
	}

	code.LoadRune(target.Pos, instruction_sp.ShortRune(str.Value()))
	return nil
}

func generateTypeIdLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, typeId *decorated.TypeIdLiteral,
	genContext *generateContext) error {
	integerValue, err := genContext.lookup.Lookup(typeId.Type())
	if err != nil {
		return err
	}
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampInt) {
		panic("wrong size")
	}

	code.LoadInteger(target.Pos, int32(integerValue))
	return nil
}

func generateIntLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, integer *decorated.IntegerLiteral) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampInt) {
		panic("wrong size")
	}

	code.LoadInteger(target.Pos, integer.Value())
	return nil
}

func generateFixedLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, fixed *decorated.FixedLiteral) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampInt) {
		panic("wrong size")
	}

	code.LoadInteger(target.Pos, fixed.Value())
	return nil
}

func generateResourceNameLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	resourceName *decorated.ResourceNameLiteral, context *assembler_sp.PackageConstants) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampInt) {
		panic("wrong size")
	}
	constant := context.AllocateResourceNameConstant(resourceName.Value())
	code.LoadInteger(target.Pos, int32(constant.ResourceID()))
	return nil
}

func generateBoolLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	boolLiteral *decorated.BooleanLiteral) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampBool) {
		panic("wrong size")
	}
	code.LoadBool(target.Pos, boolLiteral.Value())
	return nil
}
