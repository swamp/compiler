/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type FunctionCall struct {
	functionValueExpression Expression
	assignments             []Expression
	returnType              dtype.Type
	astFunctionCall         *ast.FunctionCall
}

func NewFunctionCall(astFunctionCall *ast.FunctionCall, functionType Expression, returnType dtype.Type, assignments []Expression) *FunctionCall {
	return &FunctionCall{astFunctionCall: astFunctionCall, functionValueExpression: functionType, assignments: assignments, returnType: returnType}
}

func (c *FunctionCall) AstFunctionCall() *ast.FunctionCall {
	return c.astFunctionCall
}

func (c *FunctionCall) FunctionValue() Expression {
	return c.functionValueExpression
}

func (c *FunctionCall) Arguments() []Expression {
	return c.assignments
}

func (c *FunctionCall) Type() dtype.Type {
	return c.returnType
}

func (c *FunctionCall) String() string {
	return fmt.Sprintf("[fcall %v %v]", c.functionValueExpression, c.assignments)
}

func (c *FunctionCall) FetchPositionLength() token.SourceFileReference {
	return c.astFunctionCall.FetchPositionLength()
}
