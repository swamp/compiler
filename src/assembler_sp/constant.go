/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"

	swampopcodetype "github.com/swamp/opcodes/type"
)

type ConstantType uint

const (
	ConstantTypeString ConstantType = iota
	ConstantTypeBoolean
	ConstantTypeInteger
	ConstantTypeResourceName
	ConstantTypeFunction
	ConstantTypeFunctionExternal
)

type Constant struct {
	VariableNode
	constantType ConstantType
	str          string
	b            bool
	integer      int32
}

func (v *Constant) ConstantType() ConstantType {
	return v.constantType
}

func (v *Constant) IntegerValue() int32 {
	return v.integer
}

func (v *Constant) StringValue() string {
	return v.str
}

func (v *Constant) BooleanValue() bool {
	return v.b
}

func (v *Constant) FunctionReferenceFullyQualifiedName() string {
	return v.str
}

func NewStringConstant(debugID int, str string) *Constant {
	return &Constant{VariableNode: VariableNode{someID: debugID}, constantType: ConstantTypeString, str: str}
}

func NewIntegerConstant(debugID int, i int32) *Constant {
	return &Constant{VariableNode: VariableNode{someID: debugID}, constantType: ConstantTypeInteger, integer: i}
}

func NewResourceNameConstant(debugID int, str string) *Constant {
	return &Constant{VariableNode: VariableNode{someID: debugID}, constantType: ConstantTypeResourceName, str: str}
}

func NewFunctionReferenceConstantWithDebug(debugID int, uniqueFullyQualifiedName string) *Constant {
	return &Constant{VariableNode: VariableNode{someID: debugID}, constantType: ConstantTypeFunction, str: uniqueFullyQualifiedName}
}

func NewExternalFunctionReferenceConstantWithDebug(debugID int, uniqueFullyQualifiedName string) *Constant {
	return &Constant{VariableNode: VariableNode{someID: debugID}, constantType: ConstantTypeFunctionExternal, str: uniqueFullyQualifiedName}
}

func NewBooleanConstant(debugID int, b bool) *Constant {
	return &Constant{VariableNode: VariableNode{someID: debugID}, constantType: ConstantTypeBoolean, b: b}
}

func (v *Constant) SetRegister(r swampopcodetype.Register) {
	if r.Value() == 0 {
		panic("swamp assembler: register value is zero")
	}
	v.VariableNode.register = r
	v.VariableNode.registerIsSet = true
}

func (v *Constant) Register() swampopcodetype.Register {
	if !v.VariableNode.registerIsSet {
		panic("swamp assembler: you can't read unset register")
	}
	return v.VariableNode.register
}

func (v *Constant) RegisterIsSet() bool {
	return v.VariableNode.registerIsSet
}

func (c *Constant) internalString() string {
	switch c.constantType {
	case ConstantTypeString:
		return c.str
	case ConstantTypeResourceName:
		return "@" + c.str
	case ConstantTypeBoolean:
		if c.b {
			return "True"
		}

		return "False"
	case ConstantTypeFunction:
		return "func:" + c.str
	case ConstantTypeFunctionExternal:
		return "funcExternal:" + c.str
	case ConstantTypeInteger:
		return fmt.Sprintf("int:%v", c.integer)
	}

	panic("swamp assembler: unknown constant")
}

func (c *Constant) String() string {
	if !c.VariableNode.registerIsSet {
		return fmt.Sprintf("[constant%v %v #NotSet]", c.VariableNode.someID, c.internalString())
	}

	return fmt.Sprintf("[constant%v %v #%v]", c.VariableNode.someID, c.internalString(), c.VariableNode.register)
}

func (v *Constant) SomeID() int {
	return v.someID
}
