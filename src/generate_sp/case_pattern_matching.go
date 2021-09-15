package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateCasePatternMatching(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, caseExpr *decorated.CaseForPatternMatching, genContext *generateContext) error {
	/*
		testVar, testErr := generateExpressionWithSourceVar(code, caseExpr.Test(), genContext, "cast-test")
		if testErr != nil {
			return testErr
		}

		var consequences []*assembler_sp.CaseConsequencePatternMatching

		var consequencesCodes []*assembler_sp.Code

		for _, consequence := range caseExpr.Consequences() {
			consequenceContext := *genContext
			consequenceContext.context = genContext.context.MakeScopeContext()

			consequencesCode := assembler_sp.NewCode()

			literalVariable, literalVariableErr := generateExpressionWithSourceVar(consequencesCode,
				consequence.Literal(), genContext, "literal")
			if literalVariableErr != nil {
				return literalVariableErr
			}

			labelVariableName := assembler_sp.NewVariableName("a1")
			caseLabel := consequencesCode.Label(labelVariableName, "case")

			caseExprErr := generateExpression(consequencesCode, target, consequence.Expression(), &consequenceContext)
			if caseExprErr != nil {
				return caseExprErr
			}

			asmConsequence := assembler_sp.NewCaseConsequencePatternMatching(literalVariable, caseLabel)
			consequences = append(consequences, asmConsequence)

			consequencesCodes = append(consequencesCodes, consequencesCode)

			consequenceContext.context.Free()
		}

		var defaultCase *assembler_sp.CaseConsequencePatternMatching

		if caseExpr.DefaultCase() != nil {
			consequencesCode := assembler_sp.NewCode()
			defaultContext := *genContext
			defaultContext.context = genContext.context.MakeScopeContext()

			decoratedDefault := caseExpr.DefaultCase()
			defaultLabel := consequencesCode.Label(nil, "default")
			caseExprErr := generateExpression(consequencesCode, target, decoratedDefault, &defaultContext)
			if caseExprErr != nil {
				return caseExprErr
			}
			defaultCase = assembler_sp.NewCaseConsequencePatternMatching(nil, defaultLabel)
			consequencesCodes = append(consequencesCodes, consequencesCode)
			//		endLabel := consequencesBlockCode.Label(nil, "if-end")
			defaultContext.context.Free()
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

		code.CasePatternMatching(testVar, consequences, defaultCase)

		code.Copy(consequencesBlockCode)
	*/
	return nil
}
