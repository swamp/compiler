/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import "fmt"

type ListLiteral struct {
	target   TargetStackPos
	itemSize StackRange
	values   []SourceStackPos
}

func (o *ListLiteral) String() string {
	return fmt.Sprintf("[list %v (%d) <= %v]", o.target, o.itemSize, o.values)
}
