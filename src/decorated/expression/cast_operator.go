package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type CastOperator struct {
	expression Expression
	castToType dtype.Type
}

func NewCastOperator(expression Expression, castToType dtype.Type) *CastOperator {
	return &CastOperator{expression: expression, castToType: castToType}
}

func (c *CastOperator) FetchPositionLength() token.SourceFileReference {
	return c.expression.FetchPositionLength()
}

func (c *CastOperator) Expression() Expression {
	return c.expression
}

func (c *CastOperator) String() string {
	return fmt.Sprintf("cast %v %v", c.expression, c.castToType)
}

func (c *CastOperator) Type() dtype.Type {
	return c.castToType
}
