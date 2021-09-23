package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type EnumCasePatternMatchingIntJump struct {
	matchValue int32
	jump       *opcode_sp_type.Label
}

func NewEnumCasePatternMatchingIntJump(matchValue int32, jump *opcode_sp_type.Label) EnumCasePatternMatchingIntJump {
	return EnumCasePatternMatchingIntJump{matchValue: matchValue, jump: jump}
}

func (e EnumCasePatternMatchingIntJump) String() string {
	return fmt.Sprintf("[%v %v]", e.matchValue, e.jump)
}

type PatternMatchingInt struct {
	source      opcode_sp_type.SourceStackPosition
	jumps       []EnumCasePatternMatchingIntJump
	defaultJump *opcode_sp_type.Label
}

func NewPatternMatchingInt(source opcode_sp_type.SourceStackPosition, jumps []EnumCasePatternMatchingIntJump, defaultJump *opcode_sp_type.Label) *PatternMatchingInt {
	return &PatternMatchingInt{source: source, jumps: jumps, defaultJump: defaultJump}
}

func (c *PatternMatchingInt) Write(writer OpcodeWriter) error {
	writer.Command(CmdPatternMatchingInt)
	writer.SourceStackPosition(c.source)

	writer.Count(len(c.jumps))

	var lastLabel *opcode_sp_type.Label

	for _, jump := range c.jumps {
		writer.Int32(jump.matchValue)

		if lastLabel != nil {
			writer.LabelWithOffset(jump.jump, lastLabel)
		} else {
			writer.Label(jump.jump)
		}

		lastLabel = jump.jump
	}

	writer.LabelWithOffset(c.defaultJump, lastLabel)

	return nil
}

func (c *PatternMatchingInt) String() string {
	return fmt.Sprintf("%s %v %v %v", OpcodeToMnemonic(CmdPatternMatchingInt), c.source, c.jumps, c.defaultJump)
}
