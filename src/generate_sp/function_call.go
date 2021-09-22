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
	// arguments := make([]assembler_sp.TargetStackPosRange, len(call.Arguments()))
	var arguments []assembler_sp.TargetStackPosRange
	originalParameters := call.FunctionExpression().(*decorated.FunctionReference).FunctionValue().Parameters()
	for index, arg := range call.Arguments() {
		functionArgType := originalParameters[index].Type()
		functionArgTypeUnalias := dectype.Unalias(functionArgType)
		isAny := dectype.IsAny(functionArgTypeUnalias)
		if isAny { // arg.NeedsTypeId() {
			anySourcePosGen := genContext.context.stackMemory.Allocate(uint(dectype.SizeofSwampInt), uint32(dectype.AlignOfSwampInt), "typeid")
			arguments = append(arguments, anySourcePosGen)
		} else {
			log.Printf("wasn't any %T %v\n", functionArgTypeUnalias, functionArgTypeUnalias)
		}
		arguments = append(arguments, allocMemoryForType(genContext.context.stackMemory, arg.Type(), fmt.Sprintf("arg %d", index)))
		log.Printf("argument: %d: pos:%d %T %v\n", index, arguments[index].Pos, arg, arg)
	}

	argumentIndex := 0
	for index, arg := range call.Arguments() {
		functionArgType := originalParameters[index].Type()
		functionArgTypeUnalias := dectype.Unalias(functionArgType)

		isAny := dectype.IsAny(functionArgTypeUnalias)
		if isAny { // arg.NeedsTypeId() {
			typeID, err := genContext.lookup.Lookup(arg.Type())
			if err != nil {
				return assembler_sp.SourceStackPosRange{}, err
			}
			code.LoadInteger(arguments[argumentIndex].Pos, int32(typeID))
			argumentIndex++
		}

		argReg := arguments[argumentIndex]
		argRegErr := generateExpression(code, argReg, arg, genContext)
		if argRegErr != nil {
			return assembler_sp.SourceStackPosRange{}, argRegErr
		}
		argumentIndex++
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
