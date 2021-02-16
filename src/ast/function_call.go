/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type FunctionCaller interface {
	Arguments() []Expression
	OverwriteArguments(args []Expression)
	String() string
	DebugString() string
	PositionLength() token.PositionLength
}

type FunctionCall struct {
	arguments          []Expression
	functionExpression Expression
}

func NewFunctionCall(functionExpression Expression, arguments []Expression) *FunctionCall {
	return &FunctionCall{functionExpression: functionExpression, arguments: arguments}
}

func (i *FunctionCall) Arguments() []Expression {
	return i.arguments
}

func (i *FunctionCall) OverwriteArguments(args []Expression) {
	i.arguments = args
}

func (i *FunctionCall) PositionLength() token.PositionLength {
	return i.functionExpression.PositionLength()
}

func (i *FunctionCall) FunctionExpression() Expression {
	return i.functionExpression
}

func (i *FunctionCall) String() string {
	return fmt.Sprintf("[call %v %v]", i.functionExpression, i.arguments)
}

func (i *FunctionCall) DebugString() string {
	return "[FunctionCall]"
}
