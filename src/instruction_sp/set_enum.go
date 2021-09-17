/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type SetEnum struct {
	destination opcode_sp_type.TargetStackPosition
	enumIndex   uint8
}

func (c *SetEnum) Write(writer OpcodeWriter) error {
	writer.Command(CmdSetEnum)
	writer.TargetStackPosition(c.destination)
	writer.EnumValue(c.enumIndex)

	return nil
}

func NewSetEnum(destination opcode_sp_type.TargetStackPosition,
	enumIndex uint8) *SetEnum {
	return &SetEnum{destination: destination, enumIndex: enumIndex}
}

func (c *SetEnum) String() string {
	return fmt.Sprintf("%s %v,%v", OpcodeToMnemonic(CmdSetEnum), c.destination, c.enumIndex)
}
