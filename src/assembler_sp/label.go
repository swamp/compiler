/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package assembler_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/opcode_sp_type"
)

type Label struct {
	identifier  *VariableName
	debugString string
	opLabel     *opcode_sp_type.Label
	offset      *opcode_sp_type.Label
}

func (o *Label) String() string {
	if o.identifier != nil {
		return fmt.Sprintf("%v: # (%v)]", o.identifier, o.debugString)
	}
	return fmt.Sprintf("%v:", o.debugString)
}

func (o *Label) SetOpLabel(opLabel *opcode_sp_type.Label) {
	o.opLabel = opLabel
}

func (o *Label) OpLabel() *opcode_sp_type.Label {
	return o.opLabel
}

func (o *Label) OffsetLabel() *opcode_sp_type.Label {
	return o.offset
}

func (o *Label) Name() string {
	if o.identifier != nil {
		return o.identifier.Name()
	}
	return o.debugString
}
