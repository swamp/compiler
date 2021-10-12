/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type Jump struct {
	jump *Label
}

func (o *Jump) String() string {
	return fmt.Sprintf("[jmp jump:%v]", o.jump)
}

func (o *Jump) Jump() *Label {
	return o.jump
}
