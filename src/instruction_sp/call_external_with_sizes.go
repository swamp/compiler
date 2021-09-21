/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

// CallExternalWithSizes is an instruction that calls into the embedder (usually C/C++ code).
type CallExternalWithSizes struct {
	newBasePointer opcode_sp_type.TargetStackPosition
	function       opcode_sp_type.SourceStackPosition
	sizes          []opcode_sp_type.ArgOffsetSize
}

func NewCallExternalWithSizes(newBasePointer opcode_sp_type.TargetStackPosition, function opcode_sp_type.SourceStackPosition, sizes []opcode_sp_type.ArgOffsetSize) *CallExternalWithSizes {
	return &CallExternalWithSizes{newBasePointer: newBasePointer, function: function, sizes: sizes}
}

func (c *CallExternalWithSizes) Write(writer OpcodeWriter) error {
	writer.Command(CmdCallExternalWithSizes)
	writer.TargetStackPosition(c.newBasePointer)
	writer.SourceStackPosition(c.function)

	writer.Count(len(c.sizes))
	for _, arg := range c.sizes {
		writer.ArgOffsetSize(arg)
	}

	return nil
}

func (c *CallExternalWithSizes) String() string {
	return fmt.Sprintf("%s %v %v %v", OpcodeToMnemonic(CmdCallExternalWithSizes), c.newBasePointer, c.function, c.sizes)
}
