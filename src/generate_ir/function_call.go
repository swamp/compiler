package generate_ir

import (
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"log"
)

func generateFunctionCall(call *decorated.FunctionCall,
	leafNode bool, genContext *generateContext) (value.Value, error) {
	functionValue, funcErr := handleFunctionCall(call, leafNode, genContext)
	if funcErr != nil {
		return nil, funcErr
	}

	log.Printf("func %v", functionValue)

	return functionValue, nil
}

func handleFunctionCall(call *decorated.FunctionCall, isLeafNode bool,
	genContext *generateContext) (value.Value, error) {
	functionAtom := dectype.UnaliasWithResolveInvoker(call.SmashedFunctionType()).(*dectype.FunctionAtom)
	maybeOriginalFunctionType := dectype.UnaliasWithResolveInvoker(call.FunctionExpression().Type())
	originalFunctionType, _ := maybeOriginalFunctionType.(*dectype.FunctionAtom)
	if decorated.TypeIsTemplateHasLocalTypes(functionAtom) {
		panic(fmt.Errorf("we can not call functions that has local types %v %v", call.AstFunctionCall().FetchPositionLength().ToCompleteReferenceString(), functionAtom))
	}

	fn := call.FunctionExpression()

	callExpressionFunctionValue, expressionIsFunctionValue := fn.(*decorated.FunctionReference)

	insideFunction := genContext.inFunction

	callSelf := false
	if expressionIsFunctionValue {
		callSelf = insideFunction == callExpressionFunctionValue.FunctionValue()
	}

	var functionRegister value.Value

	if !callSelf {
		var functionGenErr error
		functionRegister, functionGenErr = generateExpression(fn, false, genContext)
		if functionGenErr != nil {
			return nil, functionGenErr
		}
	} else {
		if !isLeafNode {
			return nil, fmt.Errorf("call self must be on a leafNode")
		}
	}

	expectedParameters, _ := originalFunctionType.ParameterAndReturn()
	if len(expectedParameters) < len(call.Arguments()) {
		panic(fmt.Errorf("wrong parameters %v %v", call.AstFunctionCall().FetchPositionLength().ToCompleteReferenceString(), call.AstFunctionCall()))
	}

	var argumentValues []value.Value

	for index, arg := range call.Arguments() {
		functionArgType := expectedParameters[index]
		functionArgTypeUnalias := dectype.Unalias(functionArgType)

		needsTypeId := dectype.ArgumentNeedsTypeIdInsertedBefore(functionArgTypeUnalias)
		if needsTypeId || dectype.IsTypeIdRef(arg.Type()) {
			typeID, err := genContext.lookup.Lookup(arg.Type())
			if err != nil {
				return nil, err
			}
			if dectype.IsTypeIdRef(arg.Type()) {
				unaliased := dectype.UnaliasWithResolveInvoker(arg.Type())
				primitiveAtom, _ := unaliased.(*dectype.PrimitiveAtom)
				typeID, err = genContext.lookup.Lookup(primitiveAtom)
				if err != nil {
					return nil, err
				}
			}

			argumentValues = append(argumentValues, constant.NewInt(types.I32, int64(typeID)))
			if dectype.IsTypeIdRef(functionArgTypeUnalias) {
				continue
			}
		}

		argVal, argRegErr := generateExpression(arg, false, genContext)
		if argRegErr != nil {
			return nil, argRegErr
		}
		argumentValues = append(argumentValues, argVal)
	}

	_, isExternal := decorated.CallIsExternal(fn)
	if callSelf && isExternal {
		panic(fmt.Errorf("can not be external and self"))
	}
	if isExternal {
		/*
			if annotation.IsExternalVarFunction() {
				sizes := make([]assembler_sp.VariableArgumentPosSize, len(arguments)+1)
				startVariableArgumentPos := uint(returnValue.Pos)
				sizes[0].Offset = 0
				sizes[0].Size = uint16(returnValue.Size)
				for index, argument := range arguments {
					sizes[index+1].Offset = uint16(uint(argument.Pos) - startVariableArgumentPos)
					sizes[index+1].Size = uint16(argument.Size)
				}

				code.CallExternalWithSizes(functionRegister.Pos, returnValue.Pos, sizes, filePosition)
			} else if annotation.IsExternalVarExFunction() {
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
				code.CallExternalWithSizesAndAlign(functionRegister.Pos, returnValue.Pos, sizes, filePosition)
			} else {
				code.CallExternal(functionRegister.Pos, returnValue.Pos, filePosition)
			}

		*/
	} else {
		if callSelf {
			/*
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
				code.CopyMemory(firstArgumentStackPosition, sourcePosRange, filePosition)
				code.Recur(filePosition)
				// Hack to notify that there is no source information left at this point
				returnValue.Pos = 0xffffffff
				returnValue.Size = 0

			*/
		} else {
			return ir.NewCall(functionRegister, argumentValues...), nil
		}
	}

	// This doesn't work for multiple recur / callself
	// genContext.context.stackMemory.Set(returnValue.Pos + assembler_sp.TargetStackPos(returnValue.Size))

	return nil, fmt.Errorf("not implemented yet")
}
