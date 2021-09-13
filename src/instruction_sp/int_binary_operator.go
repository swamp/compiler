/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type BinaryOperator struct {
	opcode      Commands
	a           opcode_sp_type.SourceStackPosition
	b           opcode_sp_type.SourceStackPosition
	destination opcode_sp_type.TargetStackPosition
}

func (c *BinaryOperator) Write(writer OpcodeWriter) error {
	writer.Command(c.opcode)
	writer.TargetStackPosition(c.destination)
	writer.SourceStackPosition(c.a)
	writer.SourceStackPosition(c.b)

	return nil
}

func NewBinaryOperator(opcode Commands, destination opcode_sp_type.TargetStackPosition,
	a opcode_sp_type.SourceStackPosition, b opcode_sp_type.SourceStackPosition) *BinaryOperator {
	return &BinaryOperator{opcode: opcode, destination: destination, a: a, b: b}
}

func (c *BinaryOperator) String() string {
	return fmt.Sprintf("%s %v,%v,%v", OpcodeToName(c.opcode), c.destination, c.a, c.b)
}
