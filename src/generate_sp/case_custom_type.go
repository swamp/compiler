package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

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

		var parameters []assembler_sp.SourceVariable

		for _, param := range consequence.Parameters() {
			consequenceLabelVariableName := assembler_sp.NewVariableName(param.Identifier().Name())
			paramVariable := consequenceContext.context.AllocateVariable(consequenceLabelVariableName)
			parameters = append(parameters, paramVariable)
		}

		labelVariableName := assembler_sp.NewVariableName(
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
		defaultLabel := consequencesCode.Label(nil, "default")
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

	labelVariableEndName := assembler_sp.NewVariableName("case end")
	endLabel := lastConsequnce.Label(labelVariableEndName, "caseend")

	for index, consequenceCode := range consequencesCodes {
		if index != len(consequencesCodes)-1 {
			consequenceCode.Jump(endLabel)
		}
	}

	for _, consequenceCode := range consequencesCodes {
		consequencesBlockCode.Copy(consequenceCode)
	}

	code.Case(testVar, consequences, defaultCase)

	code.Copy(consequencesBlockCode)

	return nil
}
