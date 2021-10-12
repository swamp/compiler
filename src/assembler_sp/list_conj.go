/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"
)

type ListConj struct {
	target        TargetStackPos
	item          SourceStackPos
	list          SourceStackPos
	debugItemSize StackItemSize
}

func (o *ListConj) String() string {
	return fmt.Sprintf("[ListConj %v <= item:%v (%d) list:%v]", o.target, o.item, o.debugItemSize, o.list)
}
