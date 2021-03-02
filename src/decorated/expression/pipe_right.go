/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type PipeRightOperator struct {
	BinaryOperator
	calculatedRight Expression
}

func NewPipeRightOperator(left Expression, right Expression, calculatedRight Expression) *PipeRightOperator {
	return &PipeRightOperator{
		BinaryOperator: BinaryOperator{
			ExpressionNode: ExpressionNode{decoratedType: calculatedRight.Type()},
			left:           left,
			right:          right,
		},
		calculatedRight: calculatedRight,
	}
}

func (b *PipeRightOperator) GenerateRight() Expression {
	return b.calculatedRight
}

func (b *PipeRightOperator) FetchPositionLength() token.SourceFileReference {
	inclusive := token.MakeInclusiveSourceFileReference(b.left.FetchPositionLength(), b.right.FetchPositionLength())
	return inclusive
}

func (b *PipeRightOperator) String() string {
	return fmt.Sprintf("pipe right operator \n  %v\n  %v\n  (%v)", b.left, b.right, b.calculatedRight)
}
