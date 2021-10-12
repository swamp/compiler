/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type CallExternalWithSizesAlign struct {
	function       SourceStackPos
	newBasePointer TargetStackPos
	sizes          []VariableArgumentPosSizeAlign
}

func (o *CallExternalWithSizesAlign) String() string {
	return fmt.Sprintf("[callExternalWithSizesAlign %v %v %v]", o.newBasePointer, o.function, o.sizes)
}
