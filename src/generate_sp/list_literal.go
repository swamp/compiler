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

func generateList(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	list *decorated.ListLiteral, genContext *generateContext) error {
	variables := make([]assembler_sp.SourceStackPos, len(list.Expressions()))
	for index, expr := range list.Expressions() {
		debugName := fmt.Sprintf("listliteral%v", index)
		exprVar, genErr := generateExpressionWithSourceVar(code, expr, genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar.Pos
	}
	primitive, _ := list.Type().(*dectype.PrimitiveAtom)
	firstPrimitiveType := primitive.ParameterTypes()[0]
	itemSize, itemAlign := dectype.GetMemorySizeAndAlignment(firstPrimitiveType)
	filePosition := genContext.toFilePosition(list.FetchPositionLength())
	code.ListLiteral(target.Pos, variables, assembler_sp.StackRange(itemSize), opcode_sp_type.MemoryAlign(itemAlign), filePosition)
	return nil
}

func handleList(code *assembler_sp.Code,
	list *decorated.ListLiteral, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := genContext.context.stackMemory.Allocate(uint(opcode_sp_type.Sizeof64BitPointer), uint32(opcode_sp_type.Alignof64BitPointer), "listLiteral")
	if err := generateList(code, posRange, list, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
