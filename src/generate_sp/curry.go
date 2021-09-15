package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateCurry(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, call *decorated.CurryFunction,
	genContext *generateContext) error {
	var arguments []assembler_sp.SourceStackPos

	for _, arg := range call.ArgumentsToSave() {
		argReg, argRegErr := generateExpressionWithSourceVar(code, arg, genContext, "sourceSave")
		if argRegErr != nil {
			return argRegErr
		}
		arguments = append(arguments, argReg.Pos)
	}

	functionRegister, functionGenErr := generateExpressionWithSourceVar(code,
		call.FunctionValue(), genContext, "functioncall")
	if functionGenErr != nil {
		return functionGenErr
	}

	indexIntoTypeInformationChunk, lookupErr := genContext.lookup.Lookup(call.Type())
	if lookupErr != nil {
		return lookupErr
	}

	code.Curry(target.Pos, uint16(indexIntoTypeInformationChunk), functionRegister.Pos, arguments)

	return nil
}
