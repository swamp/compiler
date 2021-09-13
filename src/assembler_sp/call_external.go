/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type CallExternal struct {
	target         TargetStackPosRange
	function       SourceStackPos
	newBasePointer SourceStackPos
}

func (o *CallExternal) String() string {
	return fmt.Sprintf("[CallExternal %v (%v)]", o.function, o.newBasePointer)
}
