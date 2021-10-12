/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type CreateArray struct {
	destination opcode_sp_type.TargetStackPosition
	arguments   []opcode_sp_type.SourceStackPosition
	itemSize    opcode_sp_type.StackRange
	itemAlign   opcode_sp_type.MemoryAlign
}

func NewCreateArray(destination opcode_sp_type.TargetStackPosition, itemSize opcode_sp_type.StackRange, itemAlign opcode_sp_type.MemoryAlign, arguments []opcode_sp_type.SourceStackPosition) *CreateArray {
	return &CreateArray{destination: destination, itemSize: itemSize, itemAlign: itemAlign, arguments: arguments}
}

func (c *CreateArray) Write(writer OpcodeWriter) error {
	writer.Command(CmdCreateArray)
	writer.TargetStackPosition(c.destination)
	writer.StackRange(c.itemSize)
	writer.MemoryAlign(c.itemAlign)

	writer.Count(len(c.arguments))
	for _, argument := range c.arguments {
		writer.SourceStackPosition(argument)
	}

	return nil
}

func (c *CreateArray) String() string {
	return fmt.Sprintf("%v %v %v (%d, %d)", OpcodeToMnemonic(CmdCreateArray), c.destination, c.arguments, c.itemSize, c.itemAlign)
}
