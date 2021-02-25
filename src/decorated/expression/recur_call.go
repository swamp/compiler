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
	assignments []Expression
	returnType  dtype.Type
}

func NewRecurCall(returnType dtype.Type, assignments []Expression) *RecurCall {
	return &RecurCall{assignments: assignments, returnType: returnType}
}

func (c *RecurCall) Arguments() []Expression {
	return c.assignments
}

func (c *RecurCall) Type() dtype.Type {
	return c.returnType
}

func (c *RecurCall) String() string {
	return fmt.Sprintf("[rcall %v]", c.assignments)
}

func (c *RecurCall) FetchPositionLength() token.SourceFileReference {
	return c.assignments[0].FetchPositionLength()
}
