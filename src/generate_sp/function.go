package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
	swamppack "github.com/swamp/pack/lib"
)

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedPackageVariableName,
	f *decorated.FunctionValue, funcContext *Context, definitions *decorator.VariableContext,
	lookup typeinfo.TypeLookup, verboseFlag verbosity.Verbosity) (*Function, error) {
	code := assembler_sp.NewCode()

	functionType := f.Type().(*dectype.FunctionTypeReference).FunctionAtom()
	unaliasedReturnType := dectype.UnaliasWithResolveInvoker(functionType.ReturnType())
	returnValueSourcePointer := allocateVariable(funcContext.scopeVariables, funcContext.stackMemory, "__return", unaliasedReturnType)
	returnValueTargetPointer := sourceToTargetStackPosRange(returnValueSourcePointer)

	for _, parameter := range f.Parameters() {
		paramVarName := parameter.Identifier().Name()
		allocateVariable(funcContext.scopeVariables, funcContext.stackMemory, paramVarName, parameter.Type())
	}

	genContext := &generateContext{
		context: funcContext,
		// definitions: definitions,
		lookup: lookup,
	}

	genErr := generateExpression(code, returnValueTargetPointer, f.Expression(), genContext)
	if genErr != nil {
		return nil, genErr
	}

	code.Return()

	opcodes, resolveErr := code.Resolve(verboseFlag >= verbosity.Mid)
	if resolveErr != nil {
		return nil, resolveErr
	}

	if verboseFlag >= verbosity.Mid {
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
