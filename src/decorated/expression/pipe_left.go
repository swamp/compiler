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
	calculatedLeft Expression
}

func NewPipeLeftOperator(left Expression, right Expression, calculatedLeftExpression Expression) *PipeLeftOperator {
	return &PipeLeftOperator{
		BinaryOperator: BinaryOperator{
			ExpressionNode: ExpressionNode{decoratedType: calculatedLeftExpression.Type()},
			left:           left,
			right:          right,
		},
		calculatedLeft: calculatedLeftExpression,
	}
}

func (b *PipeLeftOperator) GenerateLeft() Expression {
	return b.calculatedLeft
}

func (b *PipeLeftOperator) FetchPositionLength() token.SourceFileReference {
	inclusive := token.MakeInclusiveSourceFileReference(b.left.FetchPositionLength(), b.right.FetchPositionLength())
	return inclusive
}

func (b *PipeLeftOperator) String() string {
	return fmt.Sprintf("pipe left operator \n  %v\n  %v\n  (%v)", b.left, b.right, b.calculatedLeft)
}
