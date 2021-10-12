/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type ArrayLiteral struct {
	target    TargetStackPos
	itemSize  StackRange
	itemAlign opcode_sp_type.MemoryAlign
	values    []SourceStackPos
}

func (o *ArrayLiteral) String() string {
	return fmt.Sprintf("[array %v (%d, %d) <= %v]", o.target, o.itemSize, o.itemAlign, o.values)
}
