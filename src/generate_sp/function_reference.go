/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"fmt"

	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func handleFunctionReference(code *assembler_sp.Code,
	t *decorated.FunctionReference,
	stackMemory *assembler_sp.StackMemoryMapper,
	constants *assembler_sp.PackageConstants) (assembler_sp.SourceStackPosRange, error) {
	ident := t.NameReference().FullyQualifiedName()
	functionReferenceName := assembler_sp.VariableName(ident)
	foundConstant := constants.FindFunction(functionReferenceName)
	if foundConstant == nil {
		return assembler_sp.SourceStackPosRange{}, fmt.Errorf("generatesp: %v couldn't find function reference '%s' %v", t.FetchPositionLength().ToReferenceString(), functionReferenceName, t)
	}

	return constantToSourceStackPosRange(code, stackMemory, foundConstant)
}

func generateFunctionReference(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	getVar *decorated.FunctionReference, genContext *generateContext) error {
	ident := getVar.NameReference().FullyQualifiedName()
	varName := assembler_sp.VariableName(ident)
	constants := genContext.context.Constants()
	functionConstant := constants.FindFunction(varName)
	if functionConstant == nil {
		panic(fmt.Errorf("can not find function:%v", varName))
	}

	filePosition := genContext.toFilePosition(getVar.FetchPositionLength())

	code.LoadZeroMemoryPointer(target.Pos, functionConstant.PosRange().Position, filePosition)

	return nil
}
