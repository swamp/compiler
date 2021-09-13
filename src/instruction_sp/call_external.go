/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

// CallExternal is an instruction that calls into the embedder (usually C/C++ code).
type CallExternal struct {
	newBasePointer opcode_sp_type.SourceStackPosition
	function       opcode_sp_type.SourceStackPosition
	arguments      []opcode_sp_type.SourceStackPosition
}

func NewCallExternal(newBasePointer opcode_sp_type.SourceStackPosition, function opcode_sp_type.SourceStackPosition, arguments []opcode_sp_type.SourceStackPosition) *CallExternal {
	return &CallExternal{newBasePointer: newBasePointer, function: function, arguments: arguments}
}

func (c *CallExternal) Write(writer OpcodeWriter) error {
	writer.Command(CmdCallExternal)
	writer.SourceStackPosition(c.newBasePointer)
	writer.SourceStackPosition(c.function)
	writer.Count(len(c.arguments))

	for _, argument := range c.arguments {
		writer.SourceStackPosition(argument)
	}

	return nil
}

func (c *CallExternal) String() string {
	return fmt.Sprintf("callexternal %v %v (%v)", c.newBasePointer, c.function, c.arguments)
}
