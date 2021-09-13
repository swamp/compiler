/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type StringAppend struct {
	a           opcode_sp_type.SourceStackPosition
	b           opcode_sp_type.SourceStackPosition
	destination opcode_sp_type.TargetStackPosition
}

func (c *StringAppend) Write(writer OpcodeWriter) error {
	writer.Command(CmdStringAppend)
	writer.TargetStackPosition(c.destination)
	writer.SourceStackPosition(c.a)
	writer.SourceStackPosition(c.b)

	return nil
}

func NewStringAppend(destination opcode_sp_type.TargetStackPosition, a opcode_sp_type.SourceStackPosition,
	b opcode_sp_type.SourceStackPosition) *StringAppend {
	return &StringAppend{destination: destination, a: a, b: b}
}

func (c *StringAppend) String() string {
	return fmt.Sprintf("stringappend %v,%v,%v", c.destination, c.a, c.b)
}
