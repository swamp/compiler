/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type TupleLiteral struct {
	tupleType   *dectype.TupleTypeAtom
	expressions []Expression `debug:"true"`
	astLiteral  *ast.TupleLiteral
}

func NewTupleLiteral(astLiteral *ast.TupleLiteral, tupleType *dectype.TupleTypeAtom,
	expressions []Expression) *TupleLiteral {
	return &TupleLiteral{expressions: expressions, tupleType: tupleType, astLiteral: astLiteral}
}

func (c *TupleLiteral) Type() dtype.Type {
	return c.tupleType
}

func (c *TupleLiteral) TupleType() *dectype.TupleTypeAtom {
	return c.tupleType
}

func (c *TupleLiteral) AstTupleLiteral() *ast.TupleLiteral {
	return c.astLiteral
}

func (c *TupleLiteral) Expressions() []Expression {
	return c.expressions
}

func (c *TupleLiteral) String() string {
	return fmt.Sprintf("[TupleLiteral %v %v %v]", c.tupleType, c.astLiteral, c.expressions)
}

func (c *TupleLiteral) FetchPositionLength() token.SourceFileReference {
	return c.astLiteral.FetchPositionLength()
}

func (c *TupleLiteral) HumanReadable() string {
	return "Tuple Literal"
}
