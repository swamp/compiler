package decorated

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type Constant struct {
	functionValue *FunctionValue
	identifier    ast.ScopedOrNormalVariableIdentifier
}

func NewConstant(identifier ast.ScopedOrNormalVariableIdentifier, functionValue *FunctionValue) *Constant {
	return &Constant{functionValue: functionValue, identifier: identifier}
}

func (c *Constant) String() string {
	return "constant"
}

func (c *Constant) FunctionValue() *FunctionValue {
	return c.functionValue
}

func (c *Constant) Expression() Expression {
	return c.functionValue.Expression()
}

func (c *Constant) FetchPositionLength() token.SourceFileReference {
	return c.identifier.FetchPositionLength()
}

func (c *Constant) HumanReadable() string {
	return "Constant"
}

func (c *Constant) Type() dtype.Type {
	return c.functionValue.ForcedFunctionType().ReturnType()
}
