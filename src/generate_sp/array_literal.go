/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"fmt"

	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	opcode_sp_type "github.com/swamp/opcodes/type"
)

func generateArray(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	array *decorated.ArrayLiteral, genContext *generateContext) error {
	variables := make([]assembler_sp.SourceStackPos, len(array.Expressions()))
	for index, expr := range array.Expressions() {
		debugName := fmt.Sprintf("arrayliteral%v", index)
		exprVar, genErr := generateExpressionWithSourceVar(code, expr, genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar.Pos
	}
	primitiveReference, _ := array.Type().(*dectype.PrimitiveTypeReference)
	firstPrimitiveType := primitiveReference.PrimitiveAtom().ParameterTypes()[0]
	itemSize, itemAlign := dectype.GetMemorySizeAndAlignment(firstPrimitiveType)

	filePosition := genContext.toFilePosition(array.FetchPositionLength())
	code.ArrayLiteral(target.Pos, variables, assembler_sp.StackRange(itemSize), opcode_sp_type.MemoryAlign(itemAlign), filePosition)

	return nil
}

func handleArray(code *assembler_sp.Code,
	array *decorated.ArrayLiteral, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := genContext.context.stackMemory.Allocate(uint(opcode_sp_type.Sizeof64BitPointer), uint32(opcode_sp_type.Alignof64BitPointer), "listLiteral")
	if err := generateArray(code, posRange, array, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
