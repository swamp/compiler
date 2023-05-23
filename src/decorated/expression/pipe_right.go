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

type PipeRightOperator struct {
	BinaryOperator
}

func NewPipeRightOperator(left Expression, right Expression, resolvedType dtype.Type) *PipeRightOperator {
	return &PipeRightOperator{
		BinaryOperator: BinaryOperator{
			ExpressionNode: ExpressionNode{decoratedType: resolvedType},
			left:           left,
			right:          right,
		},
	}
}

func (b *PipeRightOperator) FetchPositionLength() token.SourceFileReference {
	inclusive := token.MakeInclusiveSourceFileReference(b.left.FetchPositionLength(), b.right.FetchPositionLength())
	return inclusive
}

func (b *PipeRightOperator) String() string {
	return fmt.Sprintf("%v |> %v", b.left, b.right)
}
