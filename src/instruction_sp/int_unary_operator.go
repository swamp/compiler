/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type IntUnaryOperator struct {
	opcode      Commands
	a           opcode_sp_type.SourceStackPosition
	destination opcode_sp_type.TargetStackPosition
}

func NewIntUnaryOperator(opcode Commands, destination opcode_sp_type.TargetStackPosition,
	a opcode_sp_type.SourceStackPosition) *IntUnaryOperator {
	return &IntUnaryOperator{opcode: opcode, destination: destination, a: a}
}

func (c *IntUnaryOperator) String() string {
	return fmt.Sprintf("%s %v,%v", OpcodeToName(c.opcode), c.destination, c.a)
}

func (c *IntUnaryOperator) Write(writer OpcodeWriter) error {
	writer.Command(c.opcode)
	writer.TargetStackPosition(c.destination)
	writer.SourceStackPosition(c.a)

	return nil
}
