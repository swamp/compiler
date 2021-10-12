/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type Return struct {
	stackPointerAdd uint32
}

func (o *Return) String() string {
	return fmt.Sprintf("[ret %d]", o.stackPointerAdd)
}
