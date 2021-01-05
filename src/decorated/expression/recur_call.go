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

type RecurCall struct {
	assignments  []DecoratedExpression
	returnType dtype.Type
}

func NewRecurCall(returnType dtype.Type, assignments []DecoratedExpression) *RecurCall {
	return &RecurCall{assignments: assignments, returnType: returnType}
}

func (c *RecurCall) Arguments() []DecoratedExpression {
	return c.assignments
}

func (c *RecurCall) Type() dtype.Type {
	return c.returnType
}

func (c *RecurCall) String() string {
	return fmt.Sprintf("[rcall %v]", c.assignments)
}

func (c *RecurCall) FetchPositionAndLength() token.PositionLength {
	return c.assignments[0].FetchPositionAndLength()
}
