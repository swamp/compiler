package ast

import (
	"fmt"
	"github.com/swamp/compiler/src/token"
)

type FunctionParameter struct {
	identifier    *VariableIdentifier
	parameterType Type
}

func NewFunctionParameter(identifier *VariableIdentifier, parameterType Type) *FunctionParameter {
	return &FunctionParameter{identifier: identifier, parameterType: parameterType}
}

func (i *FunctionParameter) IsIgnore() bool {
	return i.identifier.IsIgnore()
}

func (i *FunctionParameter) Identifier() *VariableIdentifier {
	return i.identifier
}

func (i *FunctionParameter) Name() string {
	return i.identifier.Name()
}

func (i *FunctionParameter) Type() Type {
	return i.parameterType
}

func (i *FunctionParameter) FetchPositionLength() token.SourceFileReference {
	return i.identifier.FetchPositionLength()
}

func (i *FunctionParameter) String() string {
	return fmt.Sprintf("[Arg %s: %s]", i.identifier.Symbol(), i.parameterType)
}

func (i *FunctionParameter) DebugString() string {
	return "[function parameter]"
}
