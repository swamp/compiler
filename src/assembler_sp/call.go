/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type Call struct {
	function       SourceStackPos
	newBasePointer TargetStackPos
}

func (o *Call) String() string {
	return fmt.Sprintf("[call %v]", o.function)
}
