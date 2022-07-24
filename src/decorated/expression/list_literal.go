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

type ListLiteral struct {
	t           dtype.Type
	expressions []Expression
	astLiteral  *ast.ListLiteral
}

func NewListLiteral(astLiteral *ast.ListLiteral, t dtype.Type, expressions []Expression) *ListLiteral {
	return &ListLiteral{t: t, expressions: expressions, astLiteral: astLiteral}
}

func (c *ListLiteral) Type() dtype.Type {
	return c.t
}

func (c *ListLiteral) AstListLiteral() *ast.ListLiteral {
	return c.astLiteral
}

func (c *ListLiteral) Expressions() []Expression {
	return c.expressions
}

func (c *ListLiteral) String() string {
	return fmt.Sprintf("[ListLiteral %v]", c.expressions)
}

func (c *ListLiteral) FetchPositionLength() token.SourceFileReference {
	return c.astLiteral.FetchPositionLength()
}

func (c *ListLiteral) HumanReadable() string {
	return "List Literal"
}
