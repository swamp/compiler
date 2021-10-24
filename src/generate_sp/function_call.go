package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func align(offset dectype.MemoryOffset, memoryAlign dectype.MemoryAlign) dectype.MemoryOffset {
	rest := dectype.MemoryAlign(uint32(offset) % uint32(memoryAlign))
	if rest != 0 {
		offset += dectype.MemoryOffset(memoryAlign - rest)
	}
	return offset
}

func handleFunctionCall(code *assembler_sp.Code, call *decorated.FunctionCall, isLeafNode bool,
	genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	functionAtom := dectype.UnaliasWithResolveInvoker(call.CompleteCalledFunctionType()).(*dectype.FunctionAtom)

	if decorated.TypeIsTemplateHasLocalTypes(functionAtom) {
		panic(fmt.Errorf("we can not call functions that has local types %v", functionAtom))
	}

	fn := call.FunctionExpression()

	insideFunction := genContext.context.inFunction

	callExpressionFunctionValue, expressionIsFunctionValue := fn.(*decorated.FunctionReference)
	callSelf := false
	if expressionIsFunctionValue {
		callSelf = insideFunction == callExpressionFunctionValue.FunctionValue()
	}

	var functionRegister assembler_sp.SourceStackPosRange

	if !callSelf {
		var functionGenErr error
		functionRegister, functionGenErr = generateExpressionWithSourceVar(code, fn, genContext, "functioncall")
		if functionGenErr != nil {
			return assembler_sp.SourceStackPosRange{}, functionGenErr
		}
	} else {
		if !isLeafNode {
			return assembler_sp.SourceStackPosRange{}, fmt.Errorf("call self must be on a leafNode")
		}
	}

	genContext.context.stackMemory.AlignUpForMax()

	invokedReturnType := dectype.UnaliasWithResolveInvoker(functionAtom.ReturnType())
	returnValue, returnValueAlign := allocMemoryForTypeEx(genContext.context.stackMemory, invokedReturnType, "returnValue")
	if uint(returnValue.Size) == 0 {
		panic(fmt.Errorf("how can it have zero size in return? %v", returnValue))
	}
	// arguments := make([]assembler_sp.TargetStackPosRange, len(call.Arguments()))
	var arguments []assembler_sp.TargetStackPosRange
	var argumentsAlign []dectype.MemoryAlign
	callFnReference := call.FunctionExpression().(*decorated.FunctionReference)
	callFn := callFnReference.FunctionValue()
	originalParameters := callFn.Parameters()
	if len(originalParameters) < len(call.Arguments()) {
		panic(fmt.Errorf("wrong parameters %v %v", call.AstFunctionCall().FetchPositionLength().ToCompleteReferenceString(), call.AstFunctionCall()))
	}
	for index, arg := range call.Arguments() {
		functionArgType := originalParameters[index].Type()
		functionArgTypeUnalias := dectype.Unalias(functionArgType)
		needsTypeId := dectype.ArgumentNeedsTypeIdInsertedBefore(functionArgTypeUnalias)
		if needsTypeId {
			anySourcePosGen := genContext.context.stackMemory.Allocate(uint(dectype.SizeofSwampInt), uint32(dectype.AlignOfSwampInt), "typeid")
			arguments = append(arguments, anySourcePosGen)
			argumentsAlign = append(argumentsAlign, dectype.AlignOfSwampInt)
		}
		argPosRange, align := allocMemoryForTypeEx(genContext.context.stackMemory, arg.Type(), fmt.Sprintf("arg %d", index))
		arguments = append(arguments, argPosRange)
		argumentsAlign = append(argumentsAlign, align)
	}

	argumentIndex := 0
	for index, arg := range call.Arguments() {
		functionArgType := originalParameters[index].Type()
		functionArgTypeUnalias := dectype.Unalias(functionArgType)

		needsTypeId := dectype.ArgumentNeedsTypeIdInsertedBefore(functionArgTypeUnalias)
		if needsTypeId || dectype.IsTypeIdRef(arg.Type()) {
			typeID, err := genContext.lookup.Lookup(arg.Type())
			if err != nil {
				return assembler_sp.SourceStackPosRange{}, err
			}
			if dectype.IsTypeIdRef(arg.Type()) {
				unaliased := dectype.UnaliasWithResolveInvoker(arg.Type())
				primitiveAtom, _ := unaliased.(*dectype.PrimitiveAtom)
				pointingToType := primitiveAtom.GenericTypes()[0]
				typeID, err = genContext.lookup.Lookup(pointingToType)
				if err != nil {
					return assembler_sp.SourceStackPosRange{}, err
				}
			}

			code.LoadInteger(arguments[argumentIndex].Pos, int32(typeID))
			argumentIndex++
			if dectype.IsTypeIdRef(functionArgTypeUnalias) {
				continue
			}
		}

		argReg := arguments[argumentIndex]
		argRegErr := generateExpression(code, argReg, arg, false, genContext)
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
		} else if functionValue.FunctionValue().Annotation().Annotation().IsExternalVarExFunction() {
			sizes := make([]assembler_sp.VariableArgumentPosSizeAlign, len(arguments)+1)
			startVariableArgumentPos := uint(returnValue.Pos)
			sizes[0].Offset = 0
			sizes[0].Size = uint16(returnValue.Size)
			sizes[0].Align = uint8(returnValueAlign)
			for index, argument := range arguments {
				sizes[index+1].Offset = uint16(uint(argument.Pos) - startVariableArgumentPos)
				sizes[index+1].Size = uint16(argument.Size)
				sizes[index+1].Align = uint8(argumentsAlign[index])
			}
			code.CallExternalWithSizesAndAlign(functionRegister.Pos, returnValue.Pos, sizes)
		} else {
			code.CallExternal(functionRegister.Pos, returnValue.Pos)
		}
	} else {
		if callSelf {
			returnSize, _ := dectype.GetMemorySizeAndAlignment(insideFunction.ForcedFunctionType().ReturnType())
			pos := dectype.MemoryOffset(returnSize)
			pos = align(pos, argumentsAlign[0])
			firstArgumentStackPosition := assembler_sp.TargetStackPos(pos)
			lastArgument := arguments[len(arguments)-1]
			octetsToCopy := int(lastArgument.Pos) + int(lastArgument.Size) - int(arguments[0].Pos)
			sourcePosRange := assembler_sp.SourceStackPosRange{
				Pos:  assembler_sp.SourceStackPos(arguments[0].Pos),
				Size: assembler_sp.SourceStackRange(octetsToCopy),
			}
			code.CopyMemory(firstArgumentStackPosition, sourcePosRange)
			code.Recur()
			// Hack to notify that there is no source information left at this point
			returnValue.Pos = 0xffffffff
			returnValue.Size = 0
		} else {
			code.Call(functionRegister.Pos, returnValue.Pos)
		}
	}

	genContext.context.stackMemory.Set(returnValue.Pos + assembler_sp.TargetStackPos(returnValue.Size))

	return targetToSourceStackPosRange(returnValue), nil
}

func generateFunctionCall(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, call *decorated.FunctionCall,
	leafNode bool, genContext *generateContext) error {
	posRange, err := handleFunctionCall(code, call, leafNode, genContext)
	if err != nil {
		return err
	}
	if posRange.Pos == 0xffffffff && posRange.Size == 0 {
		return nil
	}

	code.CopyMemory(target.Pos, posRange)

	return err
}
