package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateExpressionWithSourceVar(code *assembler_sp.Code, expr decorated.Expression,
	genContext *generateContext, debugName string) (assembler_sp.SourceStackPosRange, error) {
	switch t := expr.(type) {
	case *decorated.StringLiteral:
		constant := genContext.context.Constants().AllocateStringConstant(t.Value())
		return constantToSourceStackPosRange(code, genContext.context.stackMemory, constant)
	case *decorated.IntegerLiteral:
		{
			intStorage := genContext.context.stackMemory.Allocate(SizeofSwampInt, SizeofSwampInt, "intLiteral")
			code.LoadInteger(intStorage.Pos, t.Value())
			return targetToSourceStackPosRange(intStorage), nil
		}
	case *decorated.CharacterLiteral:
		{
			runeStorage := genContext.context.stackMemory.Allocate(SizeofSwampRune, SizeofSwampRune, "runeLiteral")
			code.LoadRune(runeStorage.Pos, uint8(t.Value()))
			return targetToSourceStackPosRange(runeStorage), nil
		}
	case *decorated.BooleanLiteral:
		{
			boolStorage := genContext.context.stackMemory.Allocate(SizeofSwampBool, SizeofSwampBool, "boolLiteral")
			code.LoadBool(boolStorage.Pos, t.Value())
			return targetToSourceStackPosRange(boolStorage), nil
		}
	case *decorated.LetVariableReference:
		letVariableReferenceName := t.LetVariable().Name().Name()
		return genContext.context.functionVariables.FindVariable(letVariableReferenceName)
	case *decorated.FunctionParameterReference:
		parameterReferenceName := t.Identifier().Name()
		return genContext.context.functionVariables.FindVariable(parameterReferenceName)
	case *decorated.FunctionReference:
		return handleFunctionReference(code, t, genContext.context.stackMemory, genContext.context.constants)
	}

	return assembler_sp.SourceStackPosRange{}, nil
}
