/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type ArrayLiteral struct {
	t               dtype.Type
	expressions     []Expression
	astArrayLiteral *ast.ArrayLiteral
}

func NewArrayLiteral(astArrayLiteral *ast.ArrayLiteral, t dtype.Type, expressions []Expression) *ArrayLiteral {
	return &ArrayLiteral{t: t, expressions: expressions, astArrayLiteral: astArrayLiteral}
}

func (c *ArrayLiteral) Type() dtype.Type {
	return c.t
}

func (c *ArrayLiteral) AstArrayLiteral() *ast.ArrayLiteral {
	return c.astArrayLiteral
}

func (c *ArrayLiteral) Expressions() []Expression {
	return c.expressions
}

func (c *ArrayLiteral) String() string {
	return fmt.Sprintf("[ArrayLiteral %v %v]", c.t.HumanReadable(), c.expressions)
}

func (c *ArrayLiteral) FetchPositionLength() token.SourceFileReference {
	return c.astArrayLiteral.FetchPositionLength()
}
