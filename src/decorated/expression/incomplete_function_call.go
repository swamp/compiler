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

type IncompleteFunctionCall struct {
	functionValueExpression Expression   `debug:"true"`
	arguments               []Expression `debug:"true"`
	inclusive               token.SourceFileReference
	assignedType            dtype.Type
}

func NewIncompleteFunctionCall(functionValueExpression Expression,
	arguments []Expression, assignedType dtype.Type) *IncompleteFunctionCall {
	//log.Printf("handling %v arguments: %v", astIncompleteFunctionCall.String(), len(arguments))

	inclusive := functionValueExpression.FetchPositionLength()
	if len(arguments) > 0 {
		inclusive = token.MakeInclusiveSourceFileReferenceFlipIfNeeded(
			functionValueExpression.FetchPositionLength(), arguments[len(arguments)-1].FetchPositionLength(),
		)
	}

	return &IncompleteFunctionCall{
		functionValueExpression: functionValueExpression,
		arguments:               arguments,
		inclusive:               inclusive,
		assignedType:            assignedType,
	}
}

func (c *IncompleteFunctionCall) FunctionExpression() Expression {
	return c.functionValueExpression
}

func (c *IncompleteFunctionCall) Arguments() []Expression {
	return c.arguments
}

func (c *IncompleteFunctionCall) Type() dtype.Type {
	return c.assignedType
}

func (c *IncompleteFunctionCall) String() string {
	return fmt.Sprintf(
		"[IncompleteFnCall %v %v %v]", c.functionValueExpression, c.arguments,
	) // c.functionValueExpression, c.arguments)
}

func (c *IncompleteFunctionCall) HumanReadable() string {
	return "Incomplete Function Call"
}

func (c *IncompleteFunctionCall) FetchPositionLength() token.SourceFileReference {
	return c.inclusive
}
