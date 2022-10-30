/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func CallIsExternal(fn Expression) (*ast.FunctionDeclarationExpression, bool) {
	fnRef, isFunctionReference := fn.(*FunctionReference)
	if isFunctionReference {
		expression := fnRef.FunctionValue().AstFunctionValue().Expression()
		declarationExpression, isDeclarationExpression := expression.(*ast.FunctionDeclarationExpression)
		if isDeclarationExpression && declarationExpression.IsSomeKindOfExternal() {
			return declarationExpression, true
		} else {
			return nil, false
		}
	}

	return nil, false
}

type FunctionCall struct {
	functionValueExpression Expression
	assignments             []Expression
	smashedFunctionType     *dectype.FunctionAtom
	astFunctionCall         *ast.FunctionCall
}

func NewFunctionCall(astFunctionCall *ast.FunctionCall, functionValueExpression Expression, smashedFunctionType *dectype.FunctionAtom, assignments []Expression) *FunctionCall {
	return &FunctionCall{astFunctionCall: astFunctionCall, functionValueExpression: functionValueExpression, assignments: assignments, smashedFunctionType: smashedFunctionType}
}

func (c *FunctionCall) AstFunctionCall() *ast.FunctionCall {
	return c.astFunctionCall
}

func (c *FunctionCall) FunctionExpression() Expression {
	return c.functionValueExpression
}

func (c *FunctionCall) Arguments() []Expression {
	return c.assignments
}

func (c *FunctionCall) Type() dtype.Type {
	return c.smashedFunctionType.ReturnType()
}

func (c *FunctionCall) SmashedFunctionType() *dectype.FunctionAtom {
	return c.smashedFunctionType
}

func (c *FunctionCall) String() string {
	return fmt.Sprintf("[FnCall %v %v]", c.functionValueExpression, c.assignments) // c.functionValueExpression, c.assignments)
}

func (c *FunctionCall) HumanReadable() string {
	return "Function Call"
}

func (c *FunctionCall) FetchPositionLength() token.SourceFileReference {
	return c.astFunctionCall.FetchPositionLength()
}
