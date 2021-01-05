/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type FunctionValue struct {
	parameters  []*VariableIdentifier
	expression  Expression
	debugAssignedValue token.VariableSymbolToken
}

func NewFunctionValue(debugAssignedValue token.VariableSymbolToken, parameters []*VariableIdentifier, expression Expression) *FunctionValue {
	return &FunctionValue{debugAssignedValue: debugAssignedValue, parameters: parameters, expression: expression}
}

func (i *FunctionValue) Parameters() []*VariableIdentifier {
	return i.parameters
}

func (i *FunctionValue) Expression() Expression {
	return i.expression
}

func (i *FunctionValue) PositionLength() token.PositionLength {
	return i.debugAssignedValue.FetchPositionLength()
}

func (i *FunctionValue) DebugFunctionIdentifier() token.VariableSymbolToken {
	return i.debugAssignedValue
}

func (i *FunctionValue) String() string {
	return fmt.Sprintf("[func (%v) -> %v]", i.parameters, i.expression)
}

func (i *FunctionValue) DebugString() string {
	return fmt.Sprintf("[function]")
}
