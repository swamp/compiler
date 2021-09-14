/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"
)

type VariableName string

func (o VariableName) String() string {
	return fmt.Sprintf("[var %s]", string(o))
}

type VariableNode struct {
	debugString string
	source      SourceStackPosRange
}

type VariableImpl struct {
	VariableNode
	identifier VariableName
}

func NewVariable(identifier VariableName, source SourceStackPosRange) *VariableImpl {
	return &VariableImpl{identifier: identifier, VariableNode: VariableNode{source: source}}
}

func (v *VariableImpl) String() string {
	return fmt.Sprintf("[var %v #%v]", v.identifier, v.VariableNode.source)
}
