/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/instruction_sp"
)

type UnaryOperator struct {
	target   TargetStackPos
	a        SourceStackPos
	operator instruction_sp.UnaryOperatorType
}

func (o *UnaryOperator) String() string {
	return fmt.Sprintf("[unary %v <= %v %v]", o.target, o.operator, o.a)
}
