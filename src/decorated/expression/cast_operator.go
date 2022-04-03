package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type CastOperator struct {
	expression Expression
	castToType *AliasReference
	infix      *ast.BinaryOperator
	inclusive  token.SourceFileReference
}

func NewCastOperator(expression Expression, castToType *AliasReference, infix *ast.BinaryOperator) *CastOperator {
	inclusive := token.MakeInclusiveSourceFileReference(expression.FetchPositionLength(), castToType.FetchPositionLength())
	return &CastOperator{expression: expression, castToType: castToType, infix: infix, inclusive: inclusive}
}

func (c *CastOperator) FetchPositionLength() token.SourceFileReference {
	return c.inclusive
}

func (c *CastOperator) Expression() Expression {
	return c.expression
}

func (c *CastOperator) String() string {
	return fmt.Sprintf("[Cast %v %v]", c.expression, c.castToType)
}

func (c *CastOperator) HumanReadable() string {
	return "Cast"
}

func (c *CastOperator) Type() dtype.Type {
	return c.castToType.Type()
}
