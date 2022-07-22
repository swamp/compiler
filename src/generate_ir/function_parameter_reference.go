package generate_ir

import (
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func handleNormalVariableLookup(varName string, ctx *parameterContext) (value.Value, error) {
	irParameter := ctx.Find(varName)
	return irParameter, nil
}

func handleLocalFunctionParameterReference(getVar *decorated.FunctionParameterReference, ctx *parameterContext) (value.Value, error) {
	varName := getVar.Identifier().Name()
	return handleNormalVariableLookup(varName, ctx)
}

func generateLocalFunctionParameterReference(getVar *decorated.FunctionParameterReference, genContext *generateContext) (value.Value, error) {
	context := genContext.parameterContext
	sourcePosRange, err := handleLocalFunctionParameterReference(getVar, context)
	if err != nil {
		return nil, err
	}

	return sourcePosRange, nil
}
