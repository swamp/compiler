/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/instruction_sp"
)

type EnumBinaryOperator struct {
	target   TargetStackPos
	a        SourceStackPos
	b        SourceStackPos
	operator instruction_sp.BinaryOperatorType
}

func (o *EnumBinaryOperator) String() string {
	return fmt.Sprintf("[ebinop %v <= %v %v %v]", o.target, o.operator, o.a, o.b)
}
