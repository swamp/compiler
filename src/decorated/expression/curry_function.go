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

type CurryFunction struct {
	functionValueExpression Expression
	curryFunctionType       *dectype.FunctionAtom
	argumentsToSave         []Expression
	astFunctionCall         *ast.FunctionCall
	originalFunctionType    *dectype.FunctionAtom
}

func NewCurryFunction(astFunctionCall *ast.FunctionCall, curryFunctionType *dectype.FunctionAtom, functionValueExpression Expression, argumentsToSave []Expression) *CurryFunction {
	originalFunctionType, _ := dectype.UnaliasWithResolveInvoker(functionValueExpression.Type()).(*dectype.FunctionAtom)
	
	return &CurryFunction{astFunctionCall: astFunctionCall, curryFunctionType: curryFunctionType, functionValueExpression: functionValueExpression, originalFunctionType: originalFunctionType, argumentsToSave: argumentsToSave}
}

func (c *CurryFunction) ArgumentsToSave() []Expression {
	return c.argumentsToSave
}

func (c *CurryFunction) FunctionAtom() *dectype.FunctionAtom {
	return c.curryFunctionType
}

func (c *CurryFunction) FunctionValue() Expression {
	return c.functionValueExpression
}

func (c *CurryFunction) OriginalFunctionType() *dectype.FunctionAtom {
	return c.originalFunctionType
}

func (c *CurryFunction) Type() dtype.Type {
	return c.curryFunctionType
}

func (c *CurryFunction) AstFunctionCall() *ast.FunctionCall {
	return c.astFunctionCall
}

func (c *CurryFunction) String() string {
	return fmt.Sprintf("[Curry %v %v]", c.functionValueExpression, c.argumentsToSave)
}

func (c *CurryFunction) HumanReadable() string {
	return "Curry Function"
}

func (c *CurryFunction) FetchPositionLength() token.SourceFileReference {
	return c.astFunctionCall.FetchPositionLength()
}
