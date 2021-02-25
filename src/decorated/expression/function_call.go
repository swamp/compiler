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
	functionType Expression
	assignments  []Expression
	returnType   dtype.Type
}

func NewFunctionCall(functionType Expression, returnType dtype.Type, assignments []Expression) *FunctionCall {
	return &FunctionCall{functionType: functionType, assignments: assignments, returnType: returnType}
}

func (c *FunctionCall) FunctionValue() Expression {
	return c.functionType
}

func (c *FunctionCall) Arguments() []Expression {
	return c.assignments
}

func (c *FunctionCall) Type() dtype.Type {
	return c.returnType
}

func (c *FunctionCall) String() string {
	return fmt.Sprintf("[fcall %v %v]", c.functionType, c.assignments)
}

func (c *FunctionCall) FetchPositionLength() token.SourceFileReference {
	return c.assignments[0].FetchPositionLength()
}
