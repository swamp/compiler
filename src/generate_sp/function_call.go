package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func generateFunctionCall(code *assembler_sp.Code, call *decorated.FunctionCall,
	genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	functionType := dectype.Unalias(call.FunctionExpression().Type())
	functionAtom, wasFunctionAtom := functionType.(*dectype.FunctionAtom)

	if !wasFunctionAtom {
		return assembler_sp.SourceStackPosRange{}, fmt.Errorf("this is not a function atom %T", functionType)
	}

	fn := call.FunctionExpression()
	functionRegister, functionGenErr := generateExpressionWithSourceVar(code, fn, genContext, "functioncall")
	if functionGenErr != nil {
		return assembler_sp.SourceStackPosRange{}, functionGenErr
	}

	returnValue := allocMemoryForType(genContext.context.stackMemory, functionAtom.ReturnType(), "returnValue")
	var arguments []assembler_sp.TargetStackPosRange
	for index, arg := range call.Arguments() {
		arguments[index] = allocMemoryForType(genContext.context.stackMemory, arg.Type(), fmt.Sprintf("arg %d", index))
	}

	for index, arg := range call.Arguments() {
		functionArgType := functionAtom.FunctionParameterTypes()[index]
		functionArgTypeUnalias := dectype.Unalias(functionArgType)

		argReg := arguments[index]
		argRegErr := generateExpression(code, argReg, arg, genContext)
		if argRegErr != nil {
			return assembler_sp.SourceStackPosRange{}, argRegErr
		}

		isAny := dectype.IsAny(functionArgTypeUnalias)
		if isAny { // arg.NeedsTypeId() {
			/*
				constant, err := generateTypeIdConstant(arg.Type(), genContext)
				if err != nil {
					return err
				}

				tempAnyConstructor := genContext.context.AllocateTempVariable("anyConstructor")
				code.Constructor(tempAnyConstructor, []assembler_sp.SourceVariable{constant, argReg})

				argReg = tempAnyConstructor

				tempVariables = append(tempVariables, tempAnyConstructor)

			*/
		}

		arguments = append(arguments, argReg)
	}

	code.Call(functionRegister.Pos, returnValue.Pos)

	return targetToSourceStackPosRange(returnValue), nil
}
