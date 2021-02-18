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

type ArrayLiteral struct {
	t           dtype.Type
	expressions []DecoratedExpression
}

func NewArrayLiteral(t dtype.Type, expressions []DecoratedExpression) *ArrayLiteral {
	return &ArrayLiteral{t: t, expressions: expressions}
}

func (c *ArrayLiteral) Type() dtype.Type {
	return c.t
}

func (c *ArrayLiteral) Expressions() []DecoratedExpression {
	return c.expressions
}

func (c *ArrayLiteral) String() string {
	return fmt.Sprintf("[ArrayLiteral %v %v]", c.t.HumanReadable(), c.expressions)
}

func (c *ArrayLiteral) FetchPositionAndLength() token.PositionLength {
	return token.PositionLength{}
}
