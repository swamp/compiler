/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type EnumCaseJump struct {
	enumValue uint8
	jump      *opcode_sp_type.Label
}

func NewEnumCaseJump(enumValue uint8, jump *opcode_sp_type.Label) EnumCaseJump {
	return EnumCaseJump{enumValue: enumValue, jump: jump}
}

func (e EnumCaseJump) String() string {
	return fmt.Sprintf("[%v %v]", e.enumValue, e.jump)
}

type EnumCase struct {
	source opcode_sp_type.SourceStackPosition
	jumps  []EnumCaseJump
}

func NewEnumCase(source opcode_sp_type.SourceStackPosition, jumps []EnumCaseJump) *EnumCase {
	return &EnumCase{source: source, jumps: jumps}
}

func (c *EnumCase) Write(writer OpcodeWriter) error {
	writer.Command(CmdEnumCase)
	writer.SourceStackPosition(c.source)

	writer.Count(len(c.jumps))

	var lastLabel *opcode_sp_type.Label

	for _, jump := range c.jumps {
		writer.EnumValue(jump.enumValue)

		if lastLabel != nil {
			writer.LabelWithOffset(jump.jump, lastLabel)
		} else {
			writer.Label(jump.jump)
		}

		lastLabel = jump.jump
	}

	return nil
}

func (c *EnumCase) String() string {
	return fmt.Sprintf("cse %v %v", c.source, c.jumps)
}
