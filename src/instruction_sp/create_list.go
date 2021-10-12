/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type CreateList struct {
	destination opcode_sp_type.TargetStackPosition
	arguments   []opcode_sp_type.SourceStackPosition
	itemSize    opcode_sp_type.StackRange
	itemAlign   opcode_sp_type.MemoryAlign
}

func NewCreateList(destination opcode_sp_type.TargetStackPosition, itemSize opcode_sp_type.StackRange, itemAlign opcode_sp_type.MemoryAlign, arguments []opcode_sp_type.SourceStackPosition) *CreateList {
	return &CreateList{destination: destination, itemSize: itemSize, itemAlign: itemAlign, arguments: arguments}
}

func (c *CreateList) Write(writer OpcodeWriter) error {
	writer.Command(CmdCreateList)
	writer.TargetStackPosition(c.destination)
	writer.StackRange(c.itemSize)
	writer.MemoryAlign(c.itemAlign)

	writer.Count(len(c.arguments))
	for _, argument := range c.arguments {
		writer.SourceStackPosition(argument)
	}

	return nil
}

func (c *CreateList) String() string {
	return fmt.Sprintf("%s %v %v (%d, %d)", OpcodeToMnemonic(CmdCreateList), c.destination, c.arguments, c.itemSize, c.itemAlign)
}
