/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type ListConj struct {
	list        opcode_sp_type.SourceStackPosition
	destination opcode_sp_type.TargetStackPosition
	item        opcode_sp_type.SourceStackPosition
}

func (c *ListConj) Write(writer OpcodeWriter) error {
	writer.Command(CmdListConj)
	writer.TargetStackPosition(c.destination)
	writer.SourceStackPosition(c.list)
	writer.SourceStackPosition(c.item)

	return nil
}

func NewListConj(destination opcode_sp_type.TargetStackPosition, item opcode_sp_type.SourceStackPosition,
	list opcode_sp_type.SourceStackPosition) *ListConj {
	return &ListConj{destination: destination, item: item, list: list}
}

func (c *ListConj) String() string {
	return fmt.Sprintf("conj %v %v %v", c.destination, c.item, c.list)
}
