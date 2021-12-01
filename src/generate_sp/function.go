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
	lookup typeinfo.TypeLookup, verboseFlag verbosity.Verbosity) (*Function, error) {
	code := assembler_sp.NewCode()

	functionType := f.Type().(*dectype.FunctionTypeReference).FunctionAtom()
	unaliasedReturnType := dectype.UnaliasWithResolveInvoker(functionType.ReturnType())
	returnValueSourcePointer, allocateVariableErr := allocateForType(funcContext.stackMemory, "__return", unaliasedReturnType)
	if allocateVariableErr != nil {
		return nil, allocateVariableErr
	}
	returnValueTargetPointer := sourceToTargetStackPosRange(returnValueSourcePointer)

	for _, parameter := range f.Parameters() {
		if _, err := allocateVariable(funcContext.scopeVariables, funcContext.stackMemory, parameter, parameter.Type()); err != nil {
			return nil, err
		}
	}

	genContext := &generateContext{
		context: funcContext,
		// definitions: definitions,
		lookup: lookup,
	}

	genErr := generateExpression(code, returnValueTargetPointer, f.Expression(), true, genContext)
	if genErr != nil {
		return nil, genErr
	}

	code.Return()

	opcodes, resolveErr := code.Resolve(verboseFlag >= verbosity.Mid)
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
		opcodes, parameterCount)

	return functionConstant, nil
}
