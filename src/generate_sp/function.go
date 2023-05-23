/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
)

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedPackageVariableName,
	f *decorated.FunctionValue, funcContext *Context,
	lookup typeinfo.TypeLookup, resourceNameLookup resourceid.ResourceNameLookup, fileCache *assembler_sp.FileUrlCache,
	verboseFlag verbosity.Verbosity) (*Function, error) {
	code := assembler_sp.NewCode()

	//functionType := f.Type().(*dectype.FunctionTypeReference).FunctionAtom()
	functionType := f.Type().(*dectype.FunctionAtom)
	unaliasedReturnType := dectype.Unalias(functionType.ReturnType())
	returnValueSourcePointer, allocateVariableErr := allocateForType(funcContext.stackMemory, "__return",
		unaliasedReturnType)
	if allocateVariableErr != nil {
		return nil, allocateVariableErr
	}
	returnValueTargetPointer := sourceToTargetStackPosRange(returnValueSourcePointer)

	for _, parameter := range f.Parameters() {
		parameterTypeID, lookupErr := lookup.Lookup(parameter.Type())
		if lookupErr != nil {
			return nil, lookupErr
		}
		if _, err := allocateVariable(code, funcContext.scopeVariables, funcContext.stackMemory, parameter,
			parameter.Type(), assembler_sp.TypeID(parameterTypeID)); err != nil {
			return nil, err
		}
	}

	genContext := &generateContext{
		context: funcContext,
		// definitions: definitions,
		lookup:             lookup,
		resourceNameLookup: resourceNameLookup,
		fileCache:          fileCache,
	}

	genErr := generateExpression(code, returnValueTargetPointer, f.Expression(), true, genContext)
	if genErr != nil {
		return nil, genErr
	}

	filePosition := genContext.toFilePosition(f.Expression().FetchPositionLength())

	endLabel := code.Label("end", "end of func")
	code.Return(filePosition)
	funcContext.scopeVariables.StopScope(endLabel)

	opcodes, debugLineInfos, resolveErr := code.Resolve(false)
	if resolveErr != nil {
		return nil, resolveErr
	}

	if verboseFlag >= verbosity.High {
		//code.PrintOut()
	}

	parameterTypes, _ := f.DeclaredFunctionTypeAtom2().ParameterAndReturn()
	parameterCount := uint(len(parameterTypes))

	signature, lookupErr := lookup.Lookup(f.Type())
	if lookupErr != nil {
		return nil, lookupErr
	}

	functionConstant := NewFunction(fullyQualifiedVariableName, TypeRef(signature),
		opcodes, parameterCount, debugLineInfos)

	return functionConstant, nil
}
