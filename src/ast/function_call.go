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
	FetchPositionLength() token.SourceFileReference
}

type FunctionCall struct {
	arguments          []Expression
	functionExpression Expression
	inclusive          token.SourceFileReference
}

func NewFunctionCall(functionExpression Expression, arguments []Expression) *FunctionCall {
	lastSourceFileRef := functionExpression.FetchPositionLength()
	if len(arguments) > 0 {
		lastSourceFileRef = arguments[len(arguments)-1].FetchPositionLength()
	}
	inclusive := token.MakeInclusiveSourceFileReference(functionExpression.FetchPositionLength(), lastSourceFileRef)
	return &FunctionCall{functionExpression: functionExpression, arguments: arguments, inclusive: inclusive}
}

func (i *FunctionCall) Arguments() []Expression {
	return i.arguments
}

func (i *FunctionCall) OverwriteArguments(args []Expression) {
	i.arguments = args
}

func (i *FunctionCall) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
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
