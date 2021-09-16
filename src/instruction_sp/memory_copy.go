/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type MemoryCopy struct {
	a           opcode_sp_type.SourceStackPositionRange
	destination opcode_sp_type.TargetStackPosition
}

func (c *MemoryCopy) Write(writer OpcodeWriter) error {
	writer.Command(CmdCopyMemory)
	writer.TargetStackPosition(c.destination)
	writer.SourceStackPositionRange(c.a)

	return nil
}

func NewMemoryCopy(destination opcode_sp_type.TargetStackPosition,
	a opcode_sp_type.SourceStackPositionRange) *MemoryCopy {
	return &MemoryCopy{destination: destination, a: a}
}

func (c *MemoryCopy) String() string {
	return fmt.Sprintf("%s %v,%v,%v", OpcodeToName(CmdCopyMemory), c.destination, c.a)
}
