/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type CasePatternMatchingJump struct {
	literal opcode_sp_type.SourceStackPosition
	jump    *opcode_sp_type.Label
}

func NewCasePatternMatchingJump(literal opcode_sp_type.SourceStackPosition, jump *opcode_sp_type.Label) CasePatternMatchingJump {
	return CasePatternMatchingJump{literal: literal, jump: jump}
}

func (e CasePatternMatchingJump) String() string {
	return fmt.Sprintf("[%v %v]", e.literal, e.jump)
}

type CasePatternMatching struct {
	source opcode_sp_type.SourceStackPositionRange
	jumps  []CasePatternMatchingJump
}

func NewCasePatternMatching(source opcode_sp_type.SourceStackPositionRange,
	jumps []CasePatternMatchingJump) *CasePatternMatching {
	return &CasePatternMatching{source: source, jumps: jumps}
}

func (c *CasePatternMatching) Write(writer OpcodeWriter) error {
	writer.Command(CmdCasePatternMatching)
	writer.SourceStackPositionRange(c.source)

	writer.Count(len(c.jumps))

	var lastLabel *opcode_sp_type.Label

	for _, jump := range c.jumps {
		writer.SourceStackPosition(jump.literal)

		if lastLabel != nil {
			writer.LabelWithOffset(jump.jump, lastLabel)
		} else {
			writer.Label(jump.jump)
		}

		lastLabel = jump.jump
	}

	return nil
}

func (c *CasePatternMatching) String() string {
	return fmt.Sprintf("csep %v %v", c.source, c.jumps)
}
