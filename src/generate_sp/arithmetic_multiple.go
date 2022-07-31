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

func generateArithmeticMultiple(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, e *decorated.ArithmeticOperator,
	genContext *generateContext) error {
	leftPrimitive, _ := dectype.UnReference(e.Left().Type()).(*dectype.PrimitiveAtom)
	switch {
	case dectype.IsListLike(e.Left().Type()) && e.OperatorType() == decorated.ArithmeticAppend:
		return generateListAppend(code, target, e, genContext)
	case leftPrimitive != nil && leftPrimitive.AtomName() == "String" && e.OperatorType() == decorated.ArithmeticAppend:
		return generateStringAppend(code, target, e, genContext)
	case dectype.IsIntLike(e.Left().Type()):
		return generateArithmetic(code, target, e, genContext)
	default:
		return fmt.Errorf("cant generate arithmetic for type: %v <-> %v (%v)",
			e.Left().Type(), e.Right().Type(), e.OperatorType())
	}
}

func handleArithmeticMultiple(code *assembler_sp.Code, e *decorated.ArithmeticOperator,
	genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	leftPrimitive, _ := dectype.UnReference(e.Left().Type()).(*dectype.PrimitiveAtom)
	var memorySize dectype.MemorySize
	var memoryAlign dectype.MemoryAlign
	switch {
	case dectype.IsListLike(e.Left().Type()) && e.OperatorType() == decorated.ArithmeticAppend:
		memorySize = dectype.MemorySize(opcode_sp_type.Sizeof64BitPointer)
		memoryAlign = dectype.MemoryAlign(opcode_sp_type.Alignof64BitPointer)
	case leftPrimitive != nil && leftPrimitive.AtomName() == "String" && e.OperatorType() == decorated.ArithmeticAppend:
		memorySize = dectype.MemorySize(opcode_sp_type.Sizeof64BitPointer)
		memoryAlign = dectype.MemoryAlign(opcode_sp_type.Alignof64BitPointer)
	case dectype.IsIntLike(e.Left().Type()):
		memorySize = dectype.MemorySize(opcode_sp_type.SizeofSwampInt)
		memoryAlign = dectype.MemoryAlign(opcode_sp_type.AlignOfSwampInt)
	default:
		panic(fmt.Errorf("cant generate arithmetic for type: %v <-> %v (%v)",
			e.Left().Type(), e.Right().Type(), e.OperatorType()))
	}

	target := genContext.context.stackMemory.Allocate(uint(memorySize), uint32(memoryAlign), "arithmetic multiple")
	if err := generateArithmeticMultiple(code, target, e, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(target), nil
}
