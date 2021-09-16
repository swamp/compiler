/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type LoadZeroMemoryPointer struct {
	destination opcode_sp_type.TargetStackPosition
	source      opcode_sp_type.SourceDynamicMemoryPosition
}

func (c *LoadZeroMemoryPointer) Write(writer OpcodeWriter) error {
	writer.Command(CmdLoadZeroMemoryPointer)
	writer.TargetStackPosition(c.destination)
	writer.SourceDynamicMemoryPosition(c.source)

	return nil
}

func NewLoadZeroMemoryPointer(destination opcode_sp_type.TargetStackPosition,
	source opcode_sp_type.SourceDynamicMemoryPosition) *LoadZeroMemoryPointer {
	return &LoadZeroMemoryPointer{destination: destination, source: source}
}

func (c *LoadZeroMemoryPointer) String() string {
	return fmt.Sprintf("%s %v,%v", OpcodeToName(CmdLoadZeroMemoryPointer), c.destination, c.source)
}
