/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

// CallExternalWithSizesAlign is an instruction that calls into the embedder (usually C/C++ code).
type CallExternalWithSizesAlign struct {
	newBasePointer opcode_sp_type.TargetStackPosition
	function       opcode_sp_type.SourceStackPosition
	sizes          []opcode_sp_type.ArgOffsetSizeAlign
}

func NewCallExternalWithSizesAlign(newBasePointer opcode_sp_type.TargetStackPosition, function opcode_sp_type.SourceStackPosition, sizes []opcode_sp_type.ArgOffsetSizeAlign) *CallExternalWithSizesAlign {
	return &CallExternalWithSizesAlign{newBasePointer: newBasePointer, function: function, sizes: sizes}
}

func (c *CallExternalWithSizesAlign) Write(writer OpcodeWriter) error {
	writer.Command(CmdCallExternalWithSizesAlign)
	writer.TargetStackPosition(c.newBasePointer)
	writer.SourceStackPosition(c.function)

	writer.Count(len(c.sizes))
	for _, arg := range c.sizes {
		writer.ArgOffsetSizeAlign(arg)
	}

	return nil
}

func (c *CallExternalWithSizesAlign) String() string {
	return fmt.Sprintf("%s %v %v %v", OpcodeToMnemonic(CmdCallExternalWithSizesAlign), c.newBasePointer, c.function, c.sizes)
}
