package generate_sp

import (
	"fmt"

	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func generateAssemblerVariable(variable *decorated.LetVariable, assignmentIndex int, tupleIndex int) assembler_sp.VariableName {
	varName := variable.Name().Name()
	if variable.IsIgnore() {
		varName = fmt.Sprintf("_%v_%v", assignmentIndex, tupleIndex)
	}

	return assembler_sp.VariableName(varName)
}

func generateLet(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, let *decorated.Let,
	genContext *generateContext) error {
	letContext := genContext.MakeScopeContext("letContext")

	var variablesInThisScope []*assembler_sp.VariableImpl
	for assignmentIndex, assignment := range let.Assignments() {
		sourceVar, sourceErr := generateExpressionWithSourceVar(code, assignment.Expression(), letContext, "let source")
		if sourceErr != nil {
			return sourceErr
		}

		if len(assignment.LetVariables()) == 1 {
			firstVar := assignment.LetVariables()[0]
			firstVarName := generateAssemblerVariable(firstVar, assignmentIndex, 0)

			newVariable := code.VariableStart(firstVarName, sourceVar)
			variablesInThisScope = append(variablesInThisScope, newVariable)

			letVariableScopeStartLabel := code.Label("scope start let", "let variable")
			typeString := assembler_sp.TypeString(assignment.Type().HumanReadable())
			if _, err := letContext.context.scopeVariables.DefineVariable(firstVarName, sourceVar, typeString, letVariableScopeStartLabel); err != nil {
				return err
			}
		} else {
			tupleType := assignment.Expression().Type().(*dectype.TupleTypeAtom)
			letVariableScopeStartLabel := code.Label("scope start let", "let variable")

			for index, tupleField := range tupleType.Fields() {
				variable := assignment.LetVariables()[index]
				fieldSourcePosRange := assembler_sp.SourceStackPosRange{
					Pos:  assembler_sp.SourceStackPos(uint(sourceVar.Pos) + uint(tupleField.MemoryOffset())),
					Size: assembler_sp.SourceStackRange(tupleField.MemorySize()),
				}

				varName := generateAssemblerVariable(variable, assignmentIndex, index)
				typeString := assembler_sp.TypeString(variable.Type().HumanReadable())
				if _, err := letContext.context.scopeVariables.DefineVariable(varName, fieldSourcePosRange, typeString, letVariableScopeStartLabel); err != nil {
					return err
				}
			}
		}
	}

	codeErr := generateExpression(code, target, let.Consequence(), true, letContext)
	if codeErr != nil {
		return codeErr
	}

	return nil
}
