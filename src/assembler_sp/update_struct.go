/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type UpdateStruct struct {
	targetClone   TargetStackPos
	structToClone SourceStackPosRange
	updates       []SourceStackPosAndRangeToLocalOffset
}

func (o *UpdateStruct) String() string {
	return fmt.Sprintf("[UpdateStruct %v <= (%v) %v]", o.targetClone, o.structToClone, o.updates)
}
