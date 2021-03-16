package decorated

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type Constant struct {
	functionReference *FunctionReference
	identifier        ast.ScopedOrNormalVariableIdentifier
}

func NewConstant(identifier ast.ScopedOrNormalVariableIdentifier, functionReference *FunctionReference) *Constant {
	return &Constant{functionReference: functionReference, identifier: identifier}
}

func (c *Constant) String() string {
	return "constant"
}

func (c *Constant) FunctionReference() *FunctionReference {
	return c.functionReference
}

func (c *Constant) Expression() Expression {
	return c.functionReference.FunctionValue().Expression()
}

func (c *Constant) FetchPositionLength() token.SourceFileReference {
	return c.identifier.FetchPositionLength()
}

func (c *Constant) HumanReadable() string {
	return "Constant"
}

func (c *Constant) Type() dtype.Type {
	return c.functionReference.FunctionValue().ForcedFunctionType().ReturnType()
}
