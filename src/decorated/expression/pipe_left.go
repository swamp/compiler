/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"

	"github.com/swamp/compiler/src/token"
)

type PipeLeftOperator struct {
	BinaryOperator
}

func NewPipeLeftOperator(left Expression, right Expression, resolvedType dtype.Type) *PipeLeftOperator {
	return &PipeLeftOperator{
		BinaryOperator: BinaryOperator{
			ExpressionNode: ExpressionNode{decoratedType: resolvedType},
			left:           left,
			right:          right,
		},
	}
}

func (b *PipeLeftOperator) FetchPositionLength() token.SourceFileReference {
	inclusive := token.MakeInclusiveSourceFileReference(b.left.FetchPositionLength(), b.right.FetchPositionLength())
	return inclusive
}

func (b *PipeLeftOperator) String() string {
	return fmt.Sprintf("%v <| %v", b.left, b.right)
}
