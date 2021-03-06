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
	functionValue TypeReferenceScopedOrNormal
}

func NewConstructorCall(functionValue TypeReferenceScopedOrNormal, arguments []Expression) *ConstructorCall {
	return &ConstructorCall{functionValue: functionValue, arguments: arguments}
}

func (i *ConstructorCall) TypeReference() TypeReferenceScopedOrNormal {
	return i.functionValue
}

func (i *ConstructorCall) Arguments() []Expression {
	return i.arguments
}

func (i *ConstructorCall) OverwriteArguments(args []Expression) {
	i.arguments = args
}

func (i *ConstructorCall) FetchPositionLength() token.SourceFileReference {
	return i.functionValue.FetchPositionLength()
}

func (i *ConstructorCall) String() string {
	if i.arguments != nil {
		return fmt.Sprintf("[ccall %v %v]", i.functionValue, i.arguments)
	}
	return fmt.Sprintf("[ccall %v]", i.functionValue)
}

func (i *ConstructorCall) DebugString() string {
	return "[ConstructorCall]"
}
