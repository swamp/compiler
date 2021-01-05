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

type ListLiteral struct {
	t           dtype.Type
	expressions []DecoratedExpression
}

func NewListLiteral(t dtype.Type, expressions []DecoratedExpression) *ListLiteral {
	return &ListLiteral{t: t, expressions: expressions}
}

func (c *ListLiteral) Type() dtype.Type {
	return c.t
}

func (c *ListLiteral) Expressions() []DecoratedExpression {
	return c.expressions
}

func (c *ListLiteral) String() string {
	return fmt.Sprintf("[ListLiteral %v %v]", c.t.HumanReadable(), c.expressions)
}

func (c *ListLiteral) FetchPositionAndLength() token.PositionLength {
	return token.PositionLength{}
}
