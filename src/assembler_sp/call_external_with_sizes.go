/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type CallExternalWithSizes struct {
	function       SourceStackPos
	newBasePointer TargetStackPos
	sizes          []VariableArgumentPosSize
}

func (o *CallExternalWithSizes) String() string {
	return fmt.Sprintf("[callExternalWithSizes %v %v %v]", o.newBasePointer, o.function, o.sizes)
}
