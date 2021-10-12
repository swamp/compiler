/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type Call struct {
	newBasePointer opcode_sp_type.TargetStackPosition
	function       opcode_sp_type.SourceStackPosition
}

func NewCall(newBasePointer opcode_sp_type.TargetStackPosition, function opcode_sp_type.SourceStackPosition) *Call {
	return &Call{newBasePointer: newBasePointer, function: function}
}

func (c *Call) Write(writer OpcodeWriter) error {
	writer.Command(CmdCall)
	writer.TargetStackPosition(c.newBasePointer)
	writer.SourceStackPosition(c.function)

	return nil
}

func (c *Call) String() string {
	return fmt.Sprintf("call %v %v", c.newBasePointer, c.function)
}
