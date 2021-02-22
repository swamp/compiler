/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type FunctionCall struct {
	functionType DecoratedExpression
	assignments  []DecoratedExpression
	returnType   dtype.Type
}

func NewFunctionCall(functionType DecoratedExpression, returnType dtype.Type, assignments []DecoratedExpression) *FunctionCall {
	return &FunctionCall{functionType: functionType, assignments: assignments, returnType: returnType}
}

func (c *FunctionCall) FunctionValue() DecoratedExpression {
	return c.functionType
}

func (c *FunctionCall) Arguments() []DecoratedExpression {
	return c.assignments
}

func (c *FunctionCall) Type() dtype.Type {
	return c.returnType
}

func (c *FunctionCall) String() string {
	return fmt.Sprintf("[fcall %v %v]", c.functionType, c.assignments)
}

func (c *FunctionCall) FetchPositionLength() token.Range {
	return c.assignments[0].FetchPositionLength()
}
