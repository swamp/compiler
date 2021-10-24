package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/instruction_sp"
)

func generateExpressionWithSourceVar(code *assembler_sp.Code, expr decorated.Expression,
	genContext *generateContext, debugName string) (assembler_sp.SourceStackPosRange, error) {
	switch t := expr.(type) {
	case *decorated.StringLiteral:
		constant := genContext.context.Constants().AllocateStringConstant(t.Value())

		return constantToSourceStackPosRange(code, genContext.context.stackMemory, constant)

	case *decorated.TypeIdLiteral:
		{
			intStorage := genContext.context.stackMemory.Allocate(uint(dectype.SizeofSwampInt), uint32(dectype.AlignOfSwampInt), "typeIdLiteral:"+t.String())
			integerValue, err := genContext.lookup.Lookup(t.Type())
			if err != nil {
				return assembler_sp.SourceStackPosRange{}, err
			}
			code.LoadInteger(intStorage.Pos, int32(integerValue))

			return targetToSourceStackPosRange(intStorage), nil
		}
	case *decorated.IntegerLiteral:
		{
			intStorage := genContext.context.stackMemory.Allocate(uint(dectype.SizeofSwampInt), uint32(dectype.AlignOfSwampInt), "intLiteral:"+t.String())
			code.LoadInteger(intStorage.Pos, t.Value())

			return targetToSourceStackPosRange(intStorage), nil
		}
	case *decorated.FixedLiteral:
		{
			fixedStorage := genContext.context.stackMemory.Allocate(uint(dectype.SizeofSwampInt), uint32(dectype.AlignOfSwampInt), "fixedLiteral:"+t.String())
			code.LoadInteger(fixedStorage.Pos, t.Value())

			return targetToSourceStackPosRange(fixedStorage), nil
		}
	case *decorated.CharacterLiteral:
		{
			runeStorage := genContext.context.stackMemory.Allocate(uint(dectype.SizeofSwampRune), uint32(dectype.AlignOfSwampRune), "runeLiteral"+t.String())
			code.LoadRune(runeStorage.Pos, instruction_sp.ShortRune(t.Value()))

			return targetToSourceStackPosRange(runeStorage), nil
		}
	case *decorated.BooleanLiteral:
		{
			boolStorage := genContext.context.stackMemory.Allocate(uint(dectype.SizeofSwampBool), uint32(dectype.AlignOfSwampBool), "boolLiteral"+t.String())
			code.LoadBool(boolStorage.Pos, t.Value())

			return targetToSourceStackPosRange(boolStorage), nil
		}
	case *decorated.LetVariableReference:
		letVariableReferenceName := t.LetVariable().Name().Name()

		return genContext.context.scopeVariables.FindVariable(letVariableReferenceName)
	case *decorated.FunctionParameterReference:
		parameterReferenceName := t.Identifier().Name()

		return genContext.context.scopeVariables.FindVariable(parameterReferenceName)
	case *decorated.CaseConsequenceParameterReference:
		parameterReferenceName := t.Identifier().Name()

		return genContext.context.scopeVariables.FindVariable(parameterReferenceName)
	case *decorated.FunctionReference:
		return handleFunctionReference(code, t, genContext.context.stackMemory, genContext.context.constants)
	case *decorated.FunctionCall:
		return handleFunctionCall(code, t, false, genContext)
	case *decorated.RecordLiteral:
		return handleRecordLiteral(code, t, genContext)
	case *decorated.RecordConstructorFromParameters:
		return handleRecordConstructorSortedAssignments(code, t, genContext)
	case *decorated.RecordConstructorFromRecord:
		return generateExpressionWithSourceVar(code, t.Expression(), genContext, debugName)
	case *decorated.CustomTypeVariantConstructor:
		return handleCustomTypeVariantConstructor(code, t, genContext)
	case *decorated.ListLiteral:
		return handleList(code, t, genContext)
	case *decorated.ArrayLiteral:
		return handleArray(code, t, genContext)
	case *decorated.BooleanOperator:
		return handleBinaryOperatorBooleanResult(code, t, genContext)
	case *decorated.ArithmeticOperator:
		return handleArithmeticMultiple(code, t, genContext)
	case *decorated.LogicalOperator:
		return handleLogical(code, t, genContext)
	case *decorated.CurryFunction:
		return handleCurry(code, t, genContext)
	case *decorated.RecordLookups:
		return handleRecordLookup(code, t, genContext)
	case *decorated.TupleLiteral:
		return handleTuple(code, t, genContext)
	case *decorated.CaseCustomType:
		return handleCaseCustomType(code, t, genContext)
	case *decorated.CaseForPatternMatching:
		return handleCasePatternMatchingMultiple(code, t, genContext)
	case *decorated.PipeLeftOperator:
		return handlePipeLeft(code, t, genContext)
	case *decorated.PipeRightOperator:
		return handlePipeRight(code, t, genContext)
	case *decorated.Guard:
		return handleGuard(code, t, genContext)
	case *decorated.If:
		return handleIf(code, t, genContext)

	}

	panic(fmt.Errorf("generate_sp_withSource: unknown node %T %v %v", expr, expr, genContext))
}
