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
	_ ConstantType = iota
	ConstantTypeString
	ConstantTypeResourceName
	ConstantTypeFunction
	ConstantTypeFunctionExternal
	ConstantTypeResourceNameChunk
)

type Constant struct {
	constantType   ConstantType
	str            string
	source         SourceDynamicMemoryPosRange
	debugString    string
	resourceNameId uint
}

func (v *Constant) ConstantType() ConstantType {
	return v.constantType
}

func (v *Constant) StringValue() string {
	return v.str
}

func (v *Constant) PosRange() SourceDynamicMemoryPosRange {
	return v.source
}

func (v *Constant) ResourceID() uint {
	return v.resourceNameId
}

func (v *Constant) FunctionReferenceFullyQualifiedName() string {
	return v.str
}

func NewStringConstant(debugString string, str string, source SourceDynamicMemoryPosRange) *Constant {
	return &Constant{constantType: ConstantTypeString, str: str, source: source, debugString: debugString}
}

func NewResourceNameConstant(resourceNameId uint, str string, source SourceDynamicMemoryPosRange) *Constant {
	return &Constant{constantType: ConstantTypeResourceName, str: str, source: source, debugString: "resourceName:" + str, resourceNameId: resourceNameId}
}

func NewFunctionReferenceConstantWithDebug(debugString string, uniqueFullyQualifiedName string, source SourceDynamicMemoryPosRange) *Constant {
	return &Constant{constantType: ConstantTypeFunction, str: uniqueFullyQualifiedName, source: source, debugString: debugString}
}

func NewExternalFunctionReferenceConstantWithDebug(debugString string, uniqueFullyQualifiedName string, source SourceDynamicMemoryPosRange) *Constant {
	return &Constant{constantType: ConstantTypeFunctionExternal, str: uniqueFullyQualifiedName, source: source, debugString: debugString}
}

func (c *Constant) internalString() string {
	switch c.constantType {
	case ConstantTypeString:
		return c.str
	case ConstantTypeResourceName:
		return "@" + c.str
	case ConstantTypeFunction:
		return "func:" + c.str
	case ConstantTypeFunctionExternal:
		return "funcExternal:" + c.str
	}

	panic("swamp assembler: unknown constant")
}

func (c *Constant) String() string {
	return fmt.Sprintf("[constant%v %v %v]", c.debugString, c.source, c.internalString())
}
