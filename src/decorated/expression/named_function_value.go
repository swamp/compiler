package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type FunctionName struct {
	identifier *ast.VariableIdentifier
}

func (f *FunctionName) Ident() *ast.VariableIdentifier {
	return f.identifier
}

func (f *FunctionName) FetchPositionLength() token.SourceFileReference {
	return f.identifier.FetchPositionLength()
}

func (f *FunctionName) String() string {
	return fmt.Sprintf("function name %v", f.identifier)
}

func (f *FunctionName) HumanReadable() string {
	return "This is the name of the function"
}

type NamedFunctionValue struct {
	identifier *FunctionName
	value      *FunctionValue
	inclusive  token.SourceFileReference
}

func NewNamedFunctionValue(identifier *ast.VariableIdentifier, value *FunctionValue) *NamedFunctionValue {
	inclusive := token.MakeInclusiveSourceFileReference(identifier.FetchPositionLength(), value.sourceFileReference)
	return &NamedFunctionValue{
		identifier: &FunctionName{identifier: identifier},
		value:      value,
		inclusive:  inclusive,
	}
}

func (n *NamedFunctionValue) String() string {
	return fmt.Sprintf("named function value %v = %v", n.identifier, n.value)
}

func (n *NamedFunctionValue) StatementString() string {
	return fmt.Sprintf("named function value %v = %v", n.identifier, n.value)
}

func (n *NamedFunctionValue) FunctionName() *FunctionName {
	return n.identifier
}

func (n *NamedFunctionValue) Value() *FunctionValue {
	return n.value
}

func (n *NamedFunctionValue) FetchPositionLength() token.SourceFileReference {
	return n.inclusive
}
