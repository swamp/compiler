package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func handleNormalVariableLookup(functionVariables *assembler_sp.ScopeVariables, varName string) (assembler_sp.SourceStackPosRange, error) {
	sourcePosRange, err := functionVariables.FindVariable(varName)
	if err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}
	return sourcePosRange, nil
}

func handleLocalFunctionParameterReference(getVar *decorated.FunctionParameterReference, functionVariables *assembler_sp.ScopeVariables) (assembler_sp.SourceStackPosRange, error) {
	varName := getVar.Identifier().Name()
	return handleNormalVariableLookup(functionVariables, varName)
}

func generateLocalFunctionParameterReference(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	getVar *decorated.FunctionParameterReference, context *Context) error {
	sourcePosRange, err := handleLocalFunctionParameterReference(getVar, context.scopeVariables)
	if err != nil {
		return err
	}

	code.CopyMemory(target.Pos, sourcePosRange)

	return nil
}

func handleLocalConsequenceParameterReference(getVar *decorated.CaseConsequenceParameterReference,
	functionVariables *assembler_sp.ScopeVariables) (assembler_sp.SourceStackPosRange, error) {
	varName := getVar.Identifier().Name()
	return handleNormalVariableLookup(functionVariables, varName)
}

func generateLocalConsequenceParameterReference(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	getVar *decorated.CaseConsequenceParameterReference, context *Context) error {
	sourcePosRange, err := handleLocalConsequenceParameterReference(getVar, context.scopeVariables)
	if err != nil {
		return err
	}
	code.CopyMemory(target.Pos, sourcePosRange)

	return nil
}

func handleLetVariableReference(getVar *decorated.LetVariableReference,
	functionVariables *assembler_sp.ScopeVariables) (assembler_sp.SourceStackPosRange, error) {
	varName := getVar.LetVariable().Name().Name()
	return handleNormalVariableLookup(functionVariables, varName)
}

func generateLetVariableReference(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	getVar *decorated.LetVariableReference, context *Context) error {
	sourcePosRange, err := handleLetVariableReference(getVar, context.scopeVariables)
	if err != nil {
		return err
	}
	code.CopyMemory(target.Pos, sourcePosRange)
	return nil
}
