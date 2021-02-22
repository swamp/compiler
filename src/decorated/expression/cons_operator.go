/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/token"
)

type ConsOperator struct {
	BinaryOperator
}

func NewConsOperator(left DecoratedExpression, right DecoratedExpression) (*ConsOperator, decshared.DecoratedError) {
	a := &ConsOperator{}
	a.BinaryOperator.left = left
	a.BinaryOperator.right = right
	a.BinaryOperator.DecoratedExpressionNode.decoratedType = right.Type()
	return a, nil
}

func (a *ConsOperator) Left() DecoratedExpression {
	return a.left
}

func (a *ConsOperator) Right() DecoratedExpression {
	return a.right
}

func (a *ConsOperator) String() string {
	return fmt.Sprintf("[cons left:%v right:%v]", a.left, a.right)
}

func (a *ConsOperator) FetchPositionLength() token.Range {
	return a.Left().FetchPositionLength()
}
