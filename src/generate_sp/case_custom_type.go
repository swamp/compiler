package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func allocateVariable(scopeVariables *assembler_sp.ScopeVariables, stackMemory *assembler_sp.StackMemoryMapper, variableName string, variableType dtype.Type) assembler_sp.SourceStackPosRange {
	targetPosRange := allocMemoryForType(stackMemory, variableType, "variable:"+variableName)
	sourcePosRange := targetToSourceStackPosRange(targetPosRange)
	scopeVariables.DefineVariable(variableName, sourcePosRange)
	return sourcePosRange
}

func generateCaseCustomType(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, caseExpr *decorated.CaseCustomType, genContext *generateContext) error {
	testVar, testErr := generateExpressionWithSourceVar(code, caseExpr.Test(), genContext, "cast-test")
	if testErr != nil {
		return testErr
	}

	var consequences []*assembler_sp.CaseConsequence

	var consequencesCodes []*assembler_sp.Code

	for _, consequence := range caseExpr.Consequences() {
		consequenceContext := genContext.MakeScopeContext()

		consequencesCode := assembler_sp.NewCode()

		fields := consequence.VariantReference().CustomTypeVariant().Fields()
		for index, param := range consequence.Parameters() {
			field := fields[index]
			consequenceLabelVariableName := param.Identifier().Name()

			paramVariable := assembler_sp.SourceStackPosRange{
				Pos:  assembler_sp.SourceStackPos(uint(testVar.Pos) + uint(field.MemoryOffset())),
				Size: assembler_sp.SourceStackRange(field.MemorySize()),
			}
			consequenceContext.context.scopeVariables.DefineVariable(consequenceLabelVariableName, paramVariable)
		}

		labelVariableName := assembler_sp.VariableName(
			consequence.VariantReference().AstIdentifier().SomeTypeIdentifier().Name())

		caseLabel := consequencesCode.Label(labelVariableName, "case")

		caseExprErr := generateExpression(consequencesCode, target, consequence.Expression(), consequenceContext)
		if caseExprErr != nil {
			return caseExprErr
		}

		asmConsequence := assembler_sp.NewCaseConsequence(uint8(consequence.InternalIndex()), caseLabel)

		consequences = append(consequences, asmConsequence)

		consequencesCodes = append(consequencesCodes, consequencesCode)

		// consequenceContext.context.Free()
	}

	var defaultCase *assembler_sp.CaseConsequence

	if caseExpr.DefaultCase() != nil {
		consequencesCode := assembler_sp.NewCode()
		defaultContext := genContext.MakeScopeContext()

		decoratedDefault := caseExpr.DefaultCase()
		defaultLabel := consequencesCode.Label("default", "default")
		caseExprErr := generateExpression(consequencesCode, target, decoratedDefault, defaultContext)
		if caseExprErr != nil {
			return caseExprErr
		}
		defaultCase = assembler_sp.NewCaseConsequence(0xff, defaultLabel)
		consequencesCodes = append(consequencesCodes, consequencesCode)
		//		endLabel := consequencesBlockCode.Label(nil, "if-end")
		// defaultContext.context.Free()
	}

	consequencesBlockCode := assembler_sp.NewCode()

	lastConsequnce := consequencesCodes[len(consequencesCodes)-1]

	labelVariableEndName := assembler_sp.VariableName("case end")
	endLabel := lastConsequnce.Label(labelVariableEndName, "caseend")

	for index, consequenceCode := range consequencesCodes {
		if index != len(consequencesCodes)-1 {
			consequenceCode.Jump(endLabel)
		}
	}

	for _, consequenceCode := range consequencesCodes {
		consequencesBlockCode.Copy(consequenceCode)
	}

	code.CaseEnum(testVar.Pos, consequences, defaultCase)

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
