package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/instruction_sp"
)

func generateExpressionWithSourceVar(code *assembler_sp.Code, expr decorated.Expression,
	genContext *generateContext, debugName string) (assembler_sp.SourceStackPosRange, error) {
	switch t := expr.(type) {
	case *decorated.StringLiteral:
		constant := genContext.context.Constants().AllocateStringConstant(t.Value())
		return constantToSourceStackPosRange(code, genContext.context.stackMemory, constant)
	case *decorated.IntegerLiteral:
		{
			intStorage := genContext.context.stackMemory.Allocate(SizeofSwampInt, AlignOfSwampInt, "intLiteral:"+t.String())
			code.LoadInteger(intStorage.Pos, t.Value())
			return targetToSourceStackPosRange(intStorage), nil
		}
	case *decorated.CharacterLiteral:
		{
			runeStorage := genContext.context.stackMemory.Allocate(SizeofSwampRune, AlignOfSwampRune, "runeLiteral"+t.String())
			code.LoadRune(runeStorage.Pos, instruction_sp.ShortRune(t.Value()))
			return targetToSourceStackPosRange(runeStorage), nil
		}
	case *decorated.BooleanLiteral:
		{
			boolStorage := genContext.context.stackMemory.Allocate(SizeofSwampBool, AlignOfSwampBool, "boolLiteral"+t.String())
			code.LoadBool(boolStorage.Pos, t.Value())
			return targetToSourceStackPosRange(boolStorage), nil
		}
	case *decorated.LetVariableReference:
		letVariableReferenceName := t.LetVariable().Name().Name()
		return genContext.context.scopeVariables.FindVariable(letVariableReferenceName)
	case *decorated.FunctionParameterReference:
		parameterReferenceName := t.Identifier().Name()
		return genContext.context.scopeVariables.FindVariable(parameterReferenceName)
	case *decorated.FunctionReference:
		return handleFunctionReference(code, t, genContext.context.stackMemory, genContext.context.constants)
	case *decorated.FunctionCall:
		return handleFunctionCall(code, t, genContext)
	case *decorated.RecordLiteral:
		return handleRecordLiteral(code, t, genContext)
	case *decorated.ListLiteral:
		return handleList(code, t, genContext)
	case *decorated.BooleanOperator:
		return handleBinaryOperatorBooleanResult(code, t, genContext)
	case *decorated.ArithmeticOperator:
		return handleArithmeticMultiple(code, t, genContext)
	case *decorated.CurryFunction:
		return handleCurry(code, t, genContext)
	}

	panic(fmt.Errorf("generate_sp_withSource: unknown node %T %v %v", expr, expr, genContext))
}
