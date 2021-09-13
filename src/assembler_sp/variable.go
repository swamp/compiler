/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"

	swampopcodetype "github.com/swamp/opcodes/type"
)

type VariableNode struct {
	context       *Context
	someID        int
	debugString   string
	register      swampopcodetype.Register
	registerIsSet bool
}

type VariableImpl struct {
	VariableNode
	identifier *VariableName
}

func NewVariable(context *Context, someID int, identifier *VariableName) *VariableImpl {
	return &VariableImpl{VariableNode: VariableNode{context: context, someID: someID, register: swampopcodetype.NewRegister(0xff)}, identifier: identifier}
}

func NewKeepVariable(context *Context, someID int, identifier *VariableName) *VariableImpl {
	return &VariableImpl{VariableNode: VariableNode{context: context, someID: someID, register: swampopcodetype.NewRegister(0xff)}, identifier: identifier}
}

func NewTempVariable(context *Context, someID int, debugString string) *VariableImpl {
	return &VariableImpl{VariableNode: VariableNode{context: context, someID: someID, debugString: debugString, register: swampopcodetype.NewRegister(0xff)}}
}

func (v *VariableImpl) SomeID() int {
	return v.someID
}

func (v *VariableImpl) SetRegister(r swampopcodetype.Register) {
	v.VariableNode.register = r
	v.VariableNode.registerIsSet = true
}

func (v *VariableImpl) Register() swampopcodetype.Register {
	if !v.VariableNode.registerIsSet {
		panic(fmt.Sprintf("swamp assembler: can not read variableimpl register %v (%v)", v.identifier, v.debugString))
	}
	return v.VariableNode.register
}

func (v *VariableImpl) RegisterIsSet() bool {
	return v.registerIsSet
}

func (v *VariableImpl) IAmTarget() bool {
	return true
}

func (v *VariableImpl) String() string {
	if v.identifier != nil {
		return fmt.Sprintf("[var%v %v #%v]", v.someID, v.identifier, v.VariableNode.register)
	}
	return fmt.Sprintf("[tmpvar%v '%v' #%v]", v.someID, v.debugString, v.VariableNode.register)
}
