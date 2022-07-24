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
	inclusive     token.SourceFileReference
}

func NewConstructorCall(functionValue TypeReferenceScopedOrNormal, arguments []Expression) *ConstructorCall {
	inclusive := functionValue.FetchPositionLength()
	if len(arguments) > 0 {
		inclusive = token.MakeInclusiveSourceFileReference(functionValue.FetchPositionLength(), arguments[len(arguments)-1].FetchPositionLength())
	}
	return &ConstructorCall{functionValue: functionValue, arguments: arguments, inclusive: inclusive}
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
	return i.inclusive
}

func (i *ConstructorCall) String() string {
	if i.arguments != nil {
		return fmt.Sprintf("[CCall %v %v]", i.functionValue, i.arguments)
	}
	return fmt.Sprintf("[CCall %v]", i.functionValue)
}

func (i *ConstructorCall) DebugString() string {
	return "[ConstructorCall]"
}
