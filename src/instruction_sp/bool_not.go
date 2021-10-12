/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import "github.com/swamp/compiler/src/opcode_sp_type"

// BoolNot inverts the specified register and puts the result in the destination register.
type BoolNot struct {
	source opcode_sp_type.SourceStackPosition
	target opcode_sp_type.TargetStackPosition
}

func NewBoolNot(target opcode_sp_type.TargetStackPosition, source opcode_sp_type.SourceStackPosition) *BoolNot {
	return &BoolNot{target: target, source: source}
}

func (c *BoolNot) Write(writer OpcodeWriter) error {
	writer.Command(CmdBoolLogicalNot)
	writer.TargetStackPosition(c.target)
	writer.SourceStackPosition(c.source)
	return nil
}

func (c *BoolNot) String() string {
	return "boolnot"
}
