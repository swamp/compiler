/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type CopyToFieldInfo struct {
	Target opcode_sp_type.TargetFieldOffset
	Source opcode_sp_type.SourceStackPositionRange
}

func (c CopyToFieldInfo) String() string {
	return fmt.Sprintf("%v<-%v", c.Target, c.Source)
}

type UpdateStruct struct {
	target           opcode_sp_type.TargetStackPosition
	source           opcode_sp_type.SourceStackPositionRange
	fieldDefinitions []CopyToFieldInfo
}

func NewUpdateStruct(target opcode_sp_type.TargetStackPosition, source opcode_sp_type.SourceStackPositionRange,
	fieldDefinitions []CopyToFieldInfo) *UpdateStruct {
	return &UpdateStruct{target: target, source: source, fieldDefinitions: fieldDefinitions}
}

func (c *UpdateStruct) Write(writer OpcodeWriter) error {
	writer.Command(CmdUpdateStruct)
	writer.TargetStackPosition(c.target)
	writer.SourceStackPositionRange(c.source)

	writer.Count(len(c.fieldDefinitions))

	for _, fieldDefinition := range c.fieldDefinitions {
		writer.TargetFieldOffset(fieldDefinition.Target)
		writer.SourceStackPositionRange(fieldDefinition.Source)
	}

	return nil
}

func (c *UpdateStruct) String() string {
	return fmt.Sprintf("update %v %v %v", c.target, c.source, c.fieldDefinitions)
}
