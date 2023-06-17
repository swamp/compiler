/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type PipeLeftOperator struct {
	BinaryOperator
}

func NewPipeLeftOperator(incompleteLeft *IncompleteFunctionCall, right Expression) *PipeLeftOperator {
	return &PipeLeftOperator{
		BinaryOperator: BinaryOperator{
			ExpressionNode: ExpressionNode{decoratedType: incompleteLeft.Type()},
			left:           incompleteLeft,
			right:          right,
		},
	}
}

func (b *PipeLeftOperator) FetchPositionLength() token.SourceFileReference {
	inclusive := token.MakeInclusiveSourceFileReference(b.left.FetchPositionLength(), b.right.FetchPositionLength())
	return inclusive
}

func (b *PipeLeftOperator) String() string {
	return fmt.Sprintf("[%v <| %v (%v)]", b.left, b.right, b.BinaryOperator.Type())
}
