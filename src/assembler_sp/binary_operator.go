/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"

	swampopcodeinst "github.com/swamp/opcodes/instruction"
)

type BinaryOperator struct {
	target   TargetStackPos
	a        SourceStackPos
	b        SourceStackPos
	r        SourceStackRange
	operator swampopcodeinst.BinaryOperatorType
}

func (o *BinaryOperator) String() string {
	return fmt.Sprintf("[binop %v <= %v %v %v]", o.target, o.operator, o.a, o.b)
}
