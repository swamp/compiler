/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type ConstructorCall struct {
	arguments     []Expression
	functionValue *TypeIdentifier
}

func NewConstructorCall(functionValue *TypeIdentifier, arguments []Expression) *ConstructorCall {
	return &ConstructorCall{functionValue: functionValue, arguments: arguments}
}

func (i *ConstructorCall) TypeIdentifier() *TypeIdentifier {
	return i.functionValue
}

func (i *ConstructorCall) Arguments() []Expression {
	return i.arguments
}

func (i *ConstructorCall) PositionLength() token.PositionLength {
	return i.functionValue.symbolToken.FetchPositionLength()
}

func (i *ConstructorCall) String() string {
	if i.arguments != nil {
		return fmt.Sprintf("[ccall %v %v]", i.functionValue, i.arguments)
	}
	return fmt.Sprintf("[ccall %v]", i.functionValue)
}

func (i *ConstructorCall) DebugString() string {
	return fmt.Sprintf("[ConstructorCall]")
}
