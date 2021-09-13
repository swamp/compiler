/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"
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
	constantType ConstantType
	str          string
	b            bool
	integer      int32
	source       SourceZeroMemoryPosRange
	debugString  string
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

func (v *Constant) PosRange() SourceZeroMemoryPosRange {
	return v.source
}

func (v *Constant) FunctionReferenceFullyQualifiedName() string {
	return v.str
}

func NewStringConstant(debugString string, str string, source SourceZeroMemoryPosRange) *Constant {
	return &Constant{constantType: ConstantTypeString, str: str, source: source, debugString: debugString}
}

func NewIntegerConstant(debugString string, i int32, source SourceZeroMemoryPosRange) *Constant {
	return &Constant{constantType: ConstantTypeInteger, integer: i, source: source, debugString: debugString}
}

func NewResourceNameConstant(debugString string, str string, source SourceZeroMemoryPosRange) *Constant {
	return &Constant{constantType: ConstantTypeResourceName, str: str, source: source, debugString: debugString}
}

func NewFunctionReferenceConstantWithDebug(debugString string, uniqueFullyQualifiedName string, source SourceZeroMemoryPosRange) *Constant {
	return &Constant{constantType: ConstantTypeFunction, str: uniqueFullyQualifiedName, source: source, debugString: debugString}
}

func NewExternalFunctionReferenceConstantWithDebug(debugString string, uniqueFullyQualifiedName string, source SourceZeroMemoryPosRange) *Constant {
	return &Constant{constantType: ConstantTypeFunctionExternal, str: uniqueFullyQualifiedName, source: source, debugString: debugString}
}

func NewBooleanConstant(debugString string, b bool, source SourceZeroMemoryPosRange) *Constant {
	return &Constant{constantType: ConstantTypeBoolean, b: b, source: source, debugString: debugString}
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
	return fmt.Sprintf("[constant%v %v %v]", c.debugString, c.source, c.internalString())
}
