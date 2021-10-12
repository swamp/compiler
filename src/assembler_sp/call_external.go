/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type CallExternal struct {
	function       SourceStackPos
	newBasePointer TargetStackPos
}

func (o *CallExternal) String() string {
	return fmt.Sprintf("[callExternal %v %v]", o.newBasePointer, o.function)
}
