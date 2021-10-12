/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

// BranchTrue branches if the test register is True.
type BranchTrue struct {
	test opcode_sp_type.SourceStackPosition
	jump *opcode_sp_type.Label
}

// NewBranchTrue creates a new BranchTrue struct.
func NewBranchTrue(test opcode_sp_type.SourceStackPosition, jump *opcode_sp_type.Label) *BranchTrue {
	return &BranchTrue{test: test, jump: jump}
}

func (c *BranchTrue) Write(writer OpcodeWriter) error {
	writer.Command(CmdBranchTrue)
	writer.SourceStackPosition(c.test)
	writer.Label(c.jump)

	return nil
}

func (c *BranchTrue) String() string {
	return fmt.Sprintf("brtr %v %v", c.test, c.jump)
}
