/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type Jump struct {
	jump *opcode_sp_type.Label
}

func (c *Jump) Write(writer OpcodeWriter) error {
	writer.Command(CmdJump)
	writer.Label(c.jump)

	return nil
}

func NewJump(jump *opcode_sp_type.Label) *Jump {
	return &Jump{jump: jump}
}

func (c *Jump) String() string {
	return fmt.Sprintf("jmp %v", c.jump)
}
