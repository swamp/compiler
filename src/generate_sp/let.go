package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
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
			if err := genContext.context.scopeVariables.DefineVariable(firstVar.Name().Name(), sourceVar); err != nil {
				return err
			}
		} else {
			tupleType := assignment.Expression().Type().(*dectype.TupleTypeAtom)
			for index, tupleField := range tupleType.Fields() {
				variable := assignment.LetVariables()[index]
				fieldSourcePosRange := assembler_sp.SourceStackPosRange{
					Pos:  assembler_sp.SourceStackPos(uint(sourceVar.Pos) + uint(tupleField.MemoryOffset())),
					Size: assembler_sp.SourceStackRange(tupleField.MemorySize()),
				}

				if err := genContext.context.scopeVariables.DefineVariable(variable.Name().Name(), fieldSourcePosRange); err != nil {
					return err
				}
			}
		}
	}

	codeErr := generateExpression(code, target, let.Consequence(), genContext)
	if codeErr != nil {
		return codeErr
	}

	return nil
}
