package generate_sp

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func handleFunctionCall(code *assembler_sp.Code, call *decorated.FunctionCall,
	genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	log.Printf("COMPLETE FUNCTION %v\n", call.CompleteCalledFunctionType().HumanReadable())

	functionAtom := dectype.UnaliasWithResolveInvoker(call.CompleteCalledFunctionType()).(*dectype.FunctionAtom)

	if decorated.TypeHasLocalTypes(functionAtom) {
		panic(fmt.Errorf("we can not call functions that has local types %v", functionAtom))
	}

	fn := call.FunctionExpression()
	functionRegister, functionGenErr := generateExpressionWithSourceVar(code, fn, genContext, "functioncall")
	if functionGenErr != nil {
		return assembler_sp.SourceStackPosRange{}, functionGenErr
	}

	invokedReturnType := dectype.UnaliasWithResolveInvoker(functionAtom.ReturnType())
	returnValue := allocMemoryForType(genContext.context.stackMemory, invokedReturnType, "returnValue")
	if uint(returnValue.Size) == 0 {
		panic(fmt.Errorf("how can it have zero size in return? %v", returnValue))
	}
	arguments := make([]assembler_sp.TargetStackPosRange, len(call.Arguments()))
	for index, arg := range call.Arguments() {
		arguments[index] = allocMemoryForType(genContext.context.stackMemory, arg.Type(), fmt.Sprintf("arg %d", index))
		log.Printf("argument: %d: pos:%d %T %v\n", index, arguments[index].Pos, arg, arg)
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

		// arguments = append(arguments, argReg)
	}

	if call.IsExternal() {
		functionValue := fn.(*decorated.FunctionReference)
		if functionValue.FunctionValue().Annotation().Annotation().IsExternalVarFunction() {
			sizes := make([]assembler_sp.VariableArgumentPosSize, len(arguments)+1)
			startVariableArgumentPos := uint(returnValue.Pos)
			sizes[0].Offset = 0
			sizes[0].Size = uint16(returnValue.Size)
			for index, argument := range arguments {
				sizes[index+1].Offset = uint16(uint(argument.Pos) - startVariableArgumentPos)
				sizes[index+1].Size = uint16(argument.Size)
			}
			code.CallExternalWithSizes(functionRegister.Pos, returnValue.Pos, sizes)
		} else {
			code.CallExternal(functionRegister.Pos, returnValue.Pos)
		}
	} else {
		code.Call(functionRegister.Pos, returnValue.Pos)
	}

	genContext.context.stackMemory.Set(arguments[0].Pos)

	return targetToSourceStackPosRange(returnValue), nil
}

func generateFunctionCall(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, call *decorated.FunctionCall,
	genContext *generateContext) error {
	posRange, err := handleFunctionCall(code, call, genContext)
	if err != nil {
		return err
	}

	code.CopyMemory(target.Pos, posRange)

	return err
}
