/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"
)

type CopyConstant struct {
	target TargetStackPos
	source SourceDynamicMemoryPos
}

func (o *CopyConstant) String() string {
	return fmt.Sprintf("[CopyConstant %v <= %v]", o.target, o.source)
}
