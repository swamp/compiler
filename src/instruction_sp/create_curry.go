/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type Curry struct {
	target         opcode_sp_type.TargetStackPosition
	typeIDConstant uint16
	function       opcode_sp_type.SourceStackPosition
	arguments      []opcode_sp_type.SourceStackPosition
}

func NewCurry(target opcode_sp_type.TargetStackPosition, typeIDConstant uint16,
	function opcode_sp_type.SourceStackPosition, arguments []opcode_sp_type.SourceStackPosition) *Curry {
	return &Curry{
		target: target, typeIDConstant: typeIDConstant,
		function: function, arguments: arguments,
	}
}

func (c *Curry) Write(writer OpcodeWriter) error {
	writer.Command(CmdCurry)
	writer.TargetStackPosition(c.target)
	writer.TypeIDConstant(c.typeIDConstant)
	writer.SourceStackPosition(c.function)
	writer.Count(len(c.arguments))

	for _, argument := range c.arguments {
		writer.SourceStackPosition(argument)
	}

	return nil
}

func (c *Curry) String() string {
	return fmt.Sprintf("curry %v (%v) %v (%v)", c.target, c.typeIDConstant, c.function, c.arguments)
}
