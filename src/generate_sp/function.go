package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
	swamppack "github.com/swamp/pack/lib"
)

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedPackageVariableName,
	f *decorated.FunctionValue, funcContext *Context,
	lookup typeinfo.TypeLookup, fileCache *assembler_sp.FileUrlCache, verboseFlag verbosity.Verbosity) (*Function, error) {
	code := assembler_sp.NewCode()

	functionType := f.Type().(*dectype.FunctionTypeReference).FunctionAtom()
	unaliasedReturnType := dectype.UnaliasWithResolveInvoker(functionType.ReturnType())
	returnValueSourcePointer, allocateVariableErr := allocateForType(funcContext.stackMemory, "__return", unaliasedReturnType)
	if allocateVariableErr != nil {
		return nil, allocateVariableErr
	}
	returnValueTargetPointer := sourceToTargetStackPosRange(returnValueSourcePointer)

	for _, parameter := range f.Parameters() {
		parameterTypeID, lookupErr := lookup.Lookup(parameter.Type())
		if lookupErr != nil {
			return nil, lookupErr
		}
		if _, err := allocateVariable(code, funcContext.scopeVariables, funcContext.stackMemory, parameter, parameter.Type(), assembler_sp.TypeID(parameterTypeID)); err != nil {
			return nil, err
		}
	}

	genContext := &generateContext{
		context: funcContext,
		// definitions: definitions,
		lookup:    lookup,
		fileCache: fileCache,
	}

	genErr := generateExpression(code, returnValueTargetPointer, f.Expression(), true, genContext)
	if genErr != nil {
		return nil, genErr
	}

	filePosition := genContext.toFilePosition(f.Expression().FetchPositionLength())

	endLabel := code.Label("end", "end of func")
	code.Return(filePosition)
	funcContext.scopeVariables.StopScope(endLabel)

	opcodes, debugLineInfos, resolveErr := code.Resolve(verboseFlag >= verbosity.Mid)
	if resolveErr != nil {
		return nil, resolveErr
	}

	if verboseFlag >= verbosity.High {
		code.PrintOut()
	}

	parameterTypes, _ := f.ForcedFunctionType().ParameterAndReturn()
	parameterCount := uint(len(parameterTypes))

	signature, lookupErr := lookup.Lookup(f.Type())
	if lookupErr != nil {
		return nil, lookupErr
	}

	functionConstant := NewFunction(fullyQualifiedVariableName, swamppack.TypeRef(signature),
		opcodes, parameterCount, debugLineInfos)

	return functionConstant, nil
}
