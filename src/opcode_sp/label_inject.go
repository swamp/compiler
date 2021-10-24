/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package opcode_sp

import "github.com/swamp/compiler/src/opcode_sp_type"

type LabelInject struct {
	referencedLabel      *opcode_sp_type.Label
	logicalOrigoPosition opcode_sp_type.ProgramCounter
	positionInStream     opcode_sp_type.ProgramCounter
	offsetLabel          *opcode_sp_type.Label
}

const (
	OctetSizeOfLabel = 2
)

func NewLabelInject(l *opcode_sp_type.Label, positionInStream opcode_sp_type.ProgramCounter) *LabelInject {
	return &LabelInject{
		referencedLabel: l, positionInStream: positionInStream,
		logicalOrigoPosition: positionInStream.Add(OctetSizeOfLabel),
	}
}

func NewLabelInjectWithOffset(l *opcode_sp_type.Label, positionInStream opcode_sp_type.ProgramCounter,
	offsetLabel *opcode_sp_type.Label) *LabelInject {
	return &LabelInject{
		referencedLabel: l, positionInStream: positionInStream,
		logicalOrigoPosition: positionInStream.Add(OctetSizeOfLabel), offsetLabel: offsetLabel,
	}
}

func (l *LabelInject) ReferencedLabel() *opcode_sp_type.Label {
	return l.referencedLabel
}

func (l *LabelInject) ForwardDeltaPC() (opcode_sp_type.DeltaPC, error) {
	targetPc := l.referencedLabel.DefinedProgramCounter()
	if l.offsetLabel != nil {
		beforePc := l.offsetLabel.DefinedProgramCounter()
		newPc, newPcErr := DeltaProgramCounter(targetPc, beforePc)
		return newPc, newPcErr
	}
	return DeltaProgramCounter(targetPc, l.logicalOrigoPosition)
}

func (l *LabelInject) LocatedAtPosition() opcode_sp_type.ProgramCounter {
	return l.positionInStream
}
