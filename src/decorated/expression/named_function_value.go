package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type NamedFunctionValue struct {
	identifier *ast.VariableIdentifier
	value      *FunctionValue
	inclusive  token.SourceFileReference
}

func NewNamedFunctionValue(identifier *ast.VariableIdentifier, value *FunctionValue) *NamedFunctionValue {
	inclusive := token.MakeInclusiveSourceFileReference(identifier.FetchPositionLength(), value.sourceFileReference)
	return &NamedFunctionValue{
		identifier: identifier,
		value:      value,
		inclusive:  inclusive,
	}
}

func (n *NamedFunctionValue) String() string {
	return fmt.Sprintf("named function value %v = %v", n.identifier, n.value)
}

func (n *NamedFunctionValue) Identifier() *ast.VariableIdentifier {
	return n.identifier
}

func (n *NamedFunctionValue) Value() *FunctionValue {
	return n.value
}

func (n *NamedFunctionValue) FetchPositionLength() token.SourceFileReference {
	return n.inclusive
}
