package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateLet(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, let *decorated.Let,
	genContext *generateContext) error {
	for _, assignment := range let.Assignments() {
		sourceVar, sourceErr := generateExpressionWithSourceVar(code, assignment.Expression(), genContext, "let source")
		if sourceErr != nil {
			return sourceErr
		}

		if len(assignment.LetVariables()) == 1 {
			firstVar := assignment.LetVariables()[0]
			genContext.context.scopeVariables.DefineVariable(firstVar.Name().Name(), sourceVar)
		} else {
		}
	}

	codeErr := generateExpression(code, target, let.Consequence(), genContext)
	if codeErr != nil {
		return codeErr
	}

	return nil
}
