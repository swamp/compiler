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
	target              opcode_sp_type.TargetStackPosition
	typeIDConstant      uint16
	firstParameterAlign opcode_sp_type.MemoryAlign
	function            opcode_sp_type.SourceStackPosition
	arguments           opcode_sp_type.SourceStackPositionRange
}

func NewCurry(target opcode_sp_type.TargetStackPosition, typeIDConstant uint16, firstParameterAlign opcode_sp_type.MemoryAlign,
	function opcode_sp_type.SourceStackPosition, arguments opcode_sp_type.SourceStackPositionRange) *Curry {
	return &Curry{
		target: target, typeIDConstant: typeIDConstant,
		firstParameterAlign: firstParameterAlign,
		function:            function, arguments: arguments,
	}
}

func (c *Curry) Write(writer OpcodeWriter) error {
	writer.Command(CmdCurry)
	writer.TargetStackPosition(c.target)
	writer.TypeIDConstant(c.typeIDConstant)
	writer.MemoryAlign(c.firstParameterAlign)
	writer.SourceStackPosition(c.function)
	writer.SourceStackPositionRange(c.arguments)

	return nil
}

func (c *Curry) String() string {
	return fmt.Sprintf("curry %v,%v,%v (typeId:%v, align:%v)", c.target, c.function, c.arguments, c.typeIDConstant, c.firstParameterAlign)
}
