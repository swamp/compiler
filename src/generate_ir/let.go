package generate_ir

import (
	"fmt"
	"github.com/llir/llvm/ir/value"
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateAssemblerVariable(variable *decorated.LetVariable, assignmentIndex int, tupleIndex int) assembler_sp.VariableName {
	varName := variable.Name().Name()
	if variable.IsIgnore() {
		varName = fmt.Sprintf("_%v_%v", assignmentIndex, tupleIndex)
	}

	return assembler_sp.VariableName(varName)
}

/*
func generateExpressionVariable(expression decorated.Expression, genContext *generateContext, p ir.Func) (ir.Func, error) {
	block := ir.NewBlock("let")
	localIdent := ir.NewLocalIdent("someThing")
	result := block.NewAnd(localIdent, localIdent)
	ir.Instruction()
}
*/

func generateLet(let *decorated.Let, genContext *generateContext) (value.Value, error) {
	//letContext := genContext.MakeScopeContext("letContext")
	a := genContext.parameterContext.Find("a")
	if a == nil {
		return nil, nil
	}

	result := genContext.block.NewAdd(a, a)

	return result, nil

	/*
			var variablesInThisScope []*assembler_sp.VariableImpl
			for assignmentIndex, assignment := range let.Assignments() {
				sourceVar, sourceErr := generateExpressionVariable(assignment.Expression(), letContext)
				if sourceErr != nil {
					return sourceErr
				}
			}

		   		if assignment.WasRecordDestructuring() {
		   			recordType := dectype.UnaliasWithResolveInvoker(assignment.Expression().Type()).(*dectype.RecordAtom)
		   			letVariableScopeStartLabel := code.Label("scope start let", "let variable")

		   			for index, letVariable := range assignment.LetVariables() {
		   				recordField := recordType.FindField(letVariable.Name().Name())

		   				fieldSourcePosRange := assembler_sp.SourceStackPosRange{
		   					Pos:  assembler_sp.SourceStackPos(uint(sourceVar.Pos) + uint(recordField.MemoryOffset())),
		   					Size: assembler_sp.SourceStackRange(recordField.MemorySize()),
		   				}

		   				varName := generateAssemblerVariable(letVariable, assignmentIndex, index)
		   				typeString := assembler_sp.TypeString(letVariable.Type().HumanReadable())
		   				letVariableTypeID, lookupErr := genContext.lookup.Lookup(letVariable.Type())
		   				if lookupErr != nil {
		   					return lookupErr
		   				}
		   				if _, err := letContext.context.scopeVariables.DefineVariable(varName, fieldSourcePosRange, assembler_sp.TypeID(letVariableTypeID), typeString, letVariableScopeStartLabel); err != nil {
		   					return err
		   				}
		   			}
		   		} else {
		   *
	*/

	/*
				if len(assignment.LetVariables()) == 1 {
					firstVar := assignment.LetVariables()[0]
					firstVarName := generateAssemblerVariable(firstVar, assignmentIndex, 0)

					newVariable := code.VariableStart(firstVarName, sourceVar)
					variablesInThisScope = append(variablesInThisScope, newVariable)

					letVariableScopeStartLabel := code.Label("scope start let", "let variable")
					typeString := assembler_sp.TypeString(assignment.Type().HumanReadable())
					letVariableTypeID, lookupErr := genContext.lookup.Lookup(assignment.Type())
					if lookupErr != nil {
						return lookupErr
					}
					if _, err := letContext.context.scopeVariables.DefineVariable(firstVarName, sourceVar, assembler_sp.TypeID(letVariableTypeID), typeString, letVariableScopeStartLabel); err != nil {
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
						letVariableTypeID, lookupErr := genContext.lookup.Lookup(variable.Type())
						if lookupErr != nil {
							return lookupErr
						}
						if _, err := letContext.context.scopeVariables.DefineVariable(varName, fieldSourcePosRange, assembler_sp.TypeID(letVariableTypeID), typeString, letVariableScopeStartLabel); err != nil {
							return err
						}
					}
				}
			}
		}

		codeErr := generateExpression(code, target, let.Consequence(), true, letContext)
		if codeErr != nil {
			return codeErr
		}

		endLabel := code.Label("end of let", "end of let")

		letContext.context.scopeVariables.StopScope(endLabel)
	*/
}
