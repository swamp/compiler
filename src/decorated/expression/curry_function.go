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

type CurryFunction struct {
	functionType    DecoratedExpression
	argumentsToSave []DecoratedExpression
}

func NewCurryFunction(functionType DecoratedExpression, argumentsToSave []DecoratedExpression) *CurryFunction {
	return &CurryFunction{functionType: functionType, argumentsToSave: argumentsToSave}
}

func (c *CurryFunction) ArgumentsToSave() []DecoratedExpression {
	return c.argumentsToSave
}

func (c *CurryFunction) FunctionValue() DecoratedExpression {
	return c.functionType
}

func (c *CurryFunction) Type() dtype.Type {
	return c.functionType.Type()
}

func (c *CurryFunction) String() string {
	return fmt.Sprintf("[curry %v %v]", c.functionType, c.argumentsToSave)
}

func (c *CurryFunction) FetchPositionLength() token.Range {
	return c.argumentsToSave[0].FetchPositionLength()
}
