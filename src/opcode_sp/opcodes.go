/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package opcode_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/instruction_sp"
	"github.com/swamp/compiler/src/opcode_sp_type"
)

type OpCode uint8

type Count uint8

func DeltaProgramCounter(after opcode_sp_type.ProgramCounter, before opcode_sp_type.ProgramCounter) (opcode_sp_type.DeltaPC, error) {
	if before.IsAfter(after) {
		panic(fmt.Sprintf("swamp opcodes: illegal delta program counter %v %v", before, after))
	}

	delta, err := before.Delta(after)
	return delta, err
}

type Instruction interface {
	Write(writer instruction_sp.OpcodeWriter) error
}
