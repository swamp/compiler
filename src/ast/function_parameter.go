/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"
	"github.com/swamp/compiler/src/token"
)

type FunctionParameter struct {
	identifier    *VariableIdentifier // allowed to be nil
	parameterType Type
	inclusive     token.SourceFileReference
}

func NewFunctionParameter(identifier *VariableIdentifier, parameterType Type) *FunctionParameter {
	inclusive := parameterType.FetchPositionLength()
	if identifier != nil {
		inclusive = token.MakeInclusiveSourceFileReference(identifier.FetchPositionLength(), parameterType.FetchPositionLength())
	}

	return &FunctionParameter{inclusive: inclusive, identifier: identifier, parameterType: parameterType}
}

func (i *FunctionParameter) IsIgnore() bool {
	if i.identifier == nil {
		return true
	}
	return i.identifier.IsIgnore()
}

func (i *FunctionParameter) Identifier() *VariableIdentifier {
	return i.identifier
}

func (i *FunctionParameter) Name() string {
	if i.identifier == nil {
		return "_"
	}
	return i.identifier.Name()
}

func (i *FunctionParameter) Type() Type {
	return i.parameterType
}

func (i *FunctionParameter) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *FunctionParameter) String() string {
	if i.identifier != nil {
		return fmt.Sprintf("[Arg %s: %s]", i.identifier.Symbol(), i.parameterType)
	}
	return fmt.Sprintf("[Arg %s]", i.parameterType)
}

func (i *FunctionParameter) DebugString() string {
	return "[function parameter]"
}
