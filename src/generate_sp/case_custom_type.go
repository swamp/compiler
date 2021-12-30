package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/opcodes/opcode_sp"
)

func allocateForType(stackMemory *assembler_sp.StackMemoryMapper, debugName string, variableType dtype.Type) (assembler_sp.SourceStackPosRange, error) {
	targetPosRange := allocMemoryForType(stackMemory, variableType, debugName)
	sourcePosRange := targetToSourceStackPosRange(targetPosRange)
	return sourcePosRange, nil
}

func allocateVariable(code *assembler_sp.Code, scopeVariables *assembler_sp.ScopeVariables, stackMemory *assembler_sp.StackMemoryMapper, variableName *decorated.FunctionParameterDefinition, variableType dtype.Type) (assembler_sp.SourceStackPosRange, error) {
	sourcePosRange, allocErr := allocateForType(stackMemory, "variable:"+variableName.Identifier().Name(), variableType)
	if allocErr != nil {
		return assembler_sp.SourceStackPosRange{}, allocErr
	}

	if !variableName.Identifier().IsIgnore() {
		startLabel := code.Label("dfsfdj", "fjskdjf")
		variableTypeString := assembler_sp.TypeString(variableType.HumanReadable())
		if _, err := scopeVariables.DefineVariable(assembler_sp.VariableName(variableName.Identifier().Name()), sourcePosRange, variableTypeString, startLabel); err != nil {
			return assembler_sp.SourceStackPosRange{}, err
		}
	}

	return sourcePosRange, nil
}

func generateCaseCustomType(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, caseExpr *decorated.CaseCustomType, genContext *generateContext) error {
	testVar, testErr := generateExpressionWithSourceVar(code, caseExpr.Test(), genContext, "cast-test")
	if testErr != nil {
		return testErr
	}

	var consequences []*assembler_sp.CaseConsequence

	var consequencesCodes []*assembler_sp.Code

	for _, consequence := range caseExpr.Consequences() {
		consequenceContext := genContext.MakeScopeContext("case")

		consequencesCode := assembler_sp.NewCode()

		fields := consequence.VariantReference().CustomTypeVariant().Fields()
		for index, param := range consequence.Parameters() {
			field := fields[index]
			consequenceLabelVariableName := assembler_sp.VariableName(param.Identifier().Name())

			paramVariable := assembler_sp.SourceStackPosRange{
				Pos:  assembler_sp.SourceStackPos(uint(testVar.Pos) + uint(field.MemoryOffset())),
				Size: assembler_sp.SourceStackRange(field.MemorySize()),
			}
			paramValidFromLabel := code.Label(assembler_sp.VariableName(param.Identifier().Name()+"_scope"), "scopeStart")
			typeString := assembler_sp.TypeString(param.Type().HumanReadable())
			if _, err := consequenceContext.context.scopeVariables.DefineVariable(consequenceLabelVariableName, paramVariable, typeString, paramValidFromLabel); err != nil {
				return err
			}
		}

		labelVariableName := assembler_sp.VariableName(
			consequence.VariantReference().AstIdentifier().SomeTypeIdentifier().Name())

		caseLabel := consequencesCode.Label(labelVariableName, "case")

		caseExprErr := generateExpression(consequencesCode, target, consequence.Expression(), false, consequenceContext)
		if caseExprErr != nil {
			return caseExprErr
		}

		asmConsequence := assembler_sp.NewCaseConsequence(uint8(consequence.InternalIndex()), caseLabel)

		consequences = append(consequences, asmConsequence)

		consequencesCodes = append(consequencesCodes, consequencesCode)

		endOfScopeLabel := consequencesCode.Label("end", "end of consequence scope")
		consequenceContext.context.scopeVariables.StopScope(endOfScopeLabel)
		// consequenceContext.context.Free()
	}

	var defaultCase *assembler_sp.CaseConsequence

	if caseExpr.DefaultCase() != nil {
		consequencesCode := assembler_sp.NewCode()
		defaultContext := genContext.MakeScopeContext("case default")

		decoratedDefault := caseExpr.DefaultCase()
		defaultLabel := consequencesCode.Label("default", "default")
		caseExprErr := generateExpression(consequencesCode, target, decoratedDefault, true, defaultContext)
		if caseExprErr != nil {
			return caseExprErr
		}
		defaultCase = assembler_sp.NewCaseConsequence(0xff, defaultLabel)
		consequencesCodes = append(consequencesCodes, consequencesCode)
		//		endLabel := consequencesBlockCode.Label(nil, "if-end")
		// defaultContext.context.Free()
		endOfScopeLabel := consequencesCode.Label("enddefault", "end of consequence scope")
		defaultContext.context.scopeVariables.StopScope(endOfScopeLabel)
	}

	consequencesBlockCode := assembler_sp.NewCode()

	lastConsequnce := consequencesCodes[len(consequencesCodes)-1]

	labelVariableEndName := assembler_sp.VariableName("case end")

	endLabel := lastConsequnce.Label(labelVariableEndName, "caseend")

	for index, consequenceCode := range consequencesCodes {
		if index != len(consequencesCodes)-1 {
			consequenceCode.Jump(endLabel, opcode_sp.FilePosition{})
		}
	}

	for _, consequenceCode := range consequencesCodes {
		consequencesBlockCode.Copy(consequenceCode)
	}

	filePosition := genContext.toFilePosition(caseExpr.Test().FetchPositionLength())
	code.CaseEnum(testVar.Pos, consequences, defaultCase, filePosition)

	code.Copy(consequencesBlockCode)

	return nil
}

func handleCaseCustomType(code *assembler_sp.Code,
	caseCustomType *decorated.CaseCustomType, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := allocMemoryForType(genContext.context.stackMemory, caseCustomType.Type(), "caseCustomTypeLiteral")
	if err := generateCaseCustomType(code, posRange, caseCustomType, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
