package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/opcodes/instruction_sp"
	opcode_sp_type "github.com/swamp/opcodes/type"
)

func generateStringLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.StringLiteral,
	genContext *generateContext) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.Sizeof64BitPointer) {
		panic("wrong size")
	}
	constants := genContext.context.constants
	constant := constants.AllocateStringConstant(str.Value())
	filePosition := genContext.toFilePosition(str.FetchPositionLength())
	code.LoadZeroMemoryPointer(target.Pos, constant.PosRange().Position, filePosition)
	return nil
}

func generateCharacterLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.CharacterLiteral, genContext *generateContext) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampInt) {
		panic("wrong size")
	}
	filePosition := genContext.toFilePosition(str.FetchPositionLength())
	code.LoadRune(target.Pos, instruction_sp.ShortRune(str.Value()), filePosition)
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

	filePosition := genContext.toFilePosition(typeId.FetchPositionLength())
	code.LoadInteger(target.Pos, int32(integerValue), filePosition)
	return nil
}

func generateIntLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, integer *decorated.IntegerLiteral, genContext *generateContext) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampInt) {
		panic("wrong size")
	}

	filePosition := genContext.toFilePosition(integer.FetchPositionLength())
	code.LoadInteger(target.Pos, integer.Value(), filePosition)
	return nil
}

func generateFixedLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, fixed *decorated.FixedLiteral, genContext *generateContext) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampInt) {
		panic("wrong size")
	}

	filePosition := genContext.toFilePosition(fixed.FetchPositionLength())
	code.LoadInteger(target.Pos, fixed.Value(), filePosition)
	return nil
}

func generateResourceNameLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	resourceName *decorated.ResourceNameLiteral, genContext *generateContext) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampInt) {
		panic("wrong size")
	}
	resourceId := genContext.resourceNameLookup.LookupResourceId(resourceName.Value())
	filePosition := genContext.toFilePosition(resourceName.FetchPositionLength())
	code.LoadInteger(target.Pos, int32(resourceId), filePosition)
	return nil
}

func generateBoolLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	boolLiteral *decorated.BooleanLiteral, genContext *generateContext) error {
	if target.Size != assembler_sp.StackRange(opcode_sp_type.SizeofSwampBool) {
		panic("wrong size")
	}
	filePosition := genContext.toFilePosition(boolLiteral.FetchPositionLength())
	code.LoadBool(target.Pos, boolLiteral.Value(), filePosition)
	return nil
}
