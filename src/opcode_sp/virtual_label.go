/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package opcode_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/instruction_sp"
	"github.com/swamp/compiler/src/opcode_sp_type"
)

type VirtualLabel struct {
	label *opcode_sp_type.Label
}

func NewVirtualLabel(label *opcode_sp_type.Label) *VirtualLabel {
	return &VirtualLabel{label: label}
}

func (c *VirtualLabel) Label() *opcode_sp_type.Label {
	return c.label
}

func (c *VirtualLabel) Write(writer instruction_sp.OpcodeWriter) error {
	return nil
}

func (c *VirtualLabel) String() string {
	return fmt.Sprintf("label: %v", c.label)
}
