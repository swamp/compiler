/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"
)

type VariableNode struct {
	context     *Context
	someID      int
	debugString string
	source      SourceStackPosRange
}

type VariableImpl struct {
	VariableNode
	identifier *VariableName
}

func NewVariable(context *Context, someID int, identifier *VariableName, source SourceStackPosRange) *VariableImpl {
	return &VariableImpl{VariableNode: VariableNode{context: context, someID: someID, source: source}}
}

func NewTempVariable(context *Context, someID int, debugString string, source SourceStackPosRange) *VariableImpl {
	return &VariableImpl{VariableNode: VariableNode{context: context, someID: someID, debugString: debugString, source: source}}
}

func (v *VariableImpl) SomeID() int {
	return v.someID
}

func (v *VariableImpl) IAmTarget() bool {
	return true
}

func (v *VariableImpl) String() string {
	if v.identifier != nil {
		return fmt.Sprintf("[var%v %v #%v]", v.someID, v.identifier, v.VariableNode.source)
	}
	return fmt.Sprintf("[tmpvar%v '%v' #%v]", v.someID, v.debugString, v.VariableNode.source)
}
