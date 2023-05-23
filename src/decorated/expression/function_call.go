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
	functionValueExpression Expression            `debug:"true"`
	arguments               []Expression          `debug:"true"`
	smashedFunctionType     *dectype.FunctionAtom `debug:"true"`
	astFunctionCall         *ast.FunctionCall
	inclusive               token.SourceFileReference
}

func NewFunctionCall(astFunctionCall *ast.FunctionCall, functionValueExpression Expression,
	smashedFunctionType *dectype.FunctionAtom, arguments []Expression) *FunctionCall {
	//log.Printf("handling %v arguments: %v", astFunctionCall.String(), len(arguments))

	inclusive := astFunctionCall.FetchPositionLength()
	if len(arguments) > 0 {
		inclusive = token.MakeInclusiveSourceFileReferenceFlipIfNeeded(
			astFunctionCall.FetchPositionLength(), arguments[len(arguments)-1].FetchPositionLength(),
		)
	}

	return &FunctionCall{
		astFunctionCall: astFunctionCall, functionValueExpression: functionValueExpression, arguments: arguments,
		smashedFunctionType: smashedFunctionType, inclusive: inclusive,
	}
}

func (c *FunctionCall) AstFunctionCall() *ast.FunctionCall {
	return c.astFunctionCall
}

func (c *FunctionCall) FunctionExpression() Expression {
	return c.functionValueExpression
}

func (c *FunctionCall) Arguments() []Expression {
	return c.arguments
}

func (c *FunctionCall) Type() dtype.Type {
	return c.smashedFunctionType.ReturnType()
}

func (c *FunctionCall) SmashedFunctionType() *dectype.FunctionAtom {
	return c.smashedFunctionType
}

func (c *FunctionCall) String() string {
	return fmt.Sprintf(
		"[FnCall %v %v %v]", c.smashedFunctionType, c.functionValueExpression, c.arguments,
	) // c.functionValueExpression, c.arguments)
}

func (c *FunctionCall) HumanReadable() string {
	return "Function Call"
}

func (c *FunctionCall) FetchPositionLength() token.SourceFileReference {
	return c.inclusive
}
