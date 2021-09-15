package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
	swamppack "github.com/swamp/pack/lib"
)

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedPackageVariableName, f *decorated.FunctionValue, root *assembler_sp.FunctionRootContext, definitions *decorator.VariableContext, lookup typeinfo.TypeLookup, verboseFlag verbosity.Verbosity) (*Function, error) {
	code := assembler_sp.NewCode()
	funcContext := root.ScopeContext()
	tempVar := root.ReturnVariable()

	for _, parameter := range f.Parameters() {
		paramVarName := assembler_sp.NewVariableName(parameter.Identifier().Name())
		funcContext.AllocateKeepParameterVariable(paramVarName)
	}

	genContext := &generateContext{
		context:     funcContext,
		definitions: definitions,
		lookup:      lookup,
	}

	genErr := generateExpression(code, tempVar, f.Expression(), genContext)
	if genErr != nil {
		return nil, genErr
	}

	code.Return(0)

	opcodes, resolveErr := code.Resolve(root, verboseFlag >= verbosity.Mid)
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
		parameterCount, root.Constants().Constants(), opcodes)

	return functionConstant, nil
}
