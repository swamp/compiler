/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type CurryFunction struct {
	functionValueExpression Expression
	curryFunctionType       *dectype.FunctionAtom
	argumentsToSave         []Expression
}

func NewCurryFunction(curryFunctionType *dectype.FunctionAtom, functionValueExpression Expression, argumentsToSave []Expression) *CurryFunction {
	return &CurryFunction{curryFunctionType: curryFunctionType, functionValueExpression: functionValueExpression, argumentsToSave: argumentsToSave}
}

func (c *CurryFunction) ArgumentsToSave() []Expression {
	return c.argumentsToSave
}

func (c *CurryFunction) FunctionValue() Expression {
	return c.functionValueExpression
}

func (c *CurryFunction) Type() dtype.Type {
	return c.curryFunctionType
}

func (c *CurryFunction) String() string {
	return fmt.Sprintf("[curry %v %v]", c.functionValueExpression, c.argumentsToSave)
}

func (c *CurryFunction) FetchPositionLength() token.SourceFileReference {
	return c.argumentsToSave[0].FetchPositionLength()
}
