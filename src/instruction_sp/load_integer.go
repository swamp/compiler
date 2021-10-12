/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type LoadInteger struct {
	destination opcode_sp_type.TargetStackPosition
	a           int32
}

func (c *LoadInteger) Write(writer OpcodeWriter) error {
	writer.Command(CmdLoadInteger)
	writer.TargetStackPosition(c.destination)
	writer.Int32(c.a)

	return nil
}

func NewLoadInteger(destination opcode_sp_type.TargetStackPosition,
	a int32) *LoadInteger {
	return &LoadInteger{destination: destination, a: a}
}

func (c *LoadInteger) String() string {
	return fmt.Sprintf("%s %v,%v", OpcodeToMnemonic(CmdLoadInteger), c.destination, c.a)
}
