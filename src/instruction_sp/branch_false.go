/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

// BranchFalse branches if the test register is False.
type BranchFalse struct {
	test opcode_sp_type.SourceStackPosition
	jump *opcode_sp_type.Label
}

// NewBranchFalse creates a new BranchFalse struct.
func NewBranchFalse(test opcode_sp_type.SourceStackPosition, jump *opcode_sp_type.Label) *BranchFalse {
	return &BranchFalse{test: test, jump: jump}
}

func (c *BranchFalse) Write(writer OpcodeWriter) error {
	writer.Command(CmdBranchFalse)
	writer.SourceStackPosition(c.test)
	writer.Label(c.jump)

	return nil
}

func (c *BranchFalse) String() string {
	return fmt.Sprintf("brfa %v %v", c.test, c.jump)
}
