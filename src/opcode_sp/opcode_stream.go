/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package opcode_sp

import (
	"encoding/binary"

	"github.com/swamp/compiler/src/instruction_sp"
	"github.com/swamp/compiler/src/opcode_sp_type"
)

type OpCodeStream struct {
	octets       []byte
	labelInjects []*LabelInject
}

func NewOpCodeStream() *OpCodeStream {
	return &OpCodeStream{}
}

func (s *OpCodeStream) LabelInjects() []*LabelInject {
	return s.labelInjects
}

func (s *OpCodeStream) Octets() []byte {
	return s.octets
}

func (s *OpCodeStream) Write(c uint8) {
	s.octets = append(s.octets, c)
}

func (s *OpCodeStream) WriteUint32(v uint32) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	s.octets = append(s.octets, b[:]...)
}

func (s *OpCodeStream) WriteUint16(v uint16) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	s.octets = append(s.octets, b[:]...)
}

func (s *OpCodeStream) ArgOffsetSize(r opcode_sp_type.ArgOffsetSize) {
	s.WriteUint16(r.Offset)
	s.WriteUint16(r.Size)
}

func (s *OpCodeStream) StackPosition(r opcode_sp_type.SourceStackPosition) {
	s.WriteUint32(uint32(r))
}

func (s *OpCodeStream) SourceStackPosition(r opcode_sp_type.SourceStackPosition) {
	s.WriteUint32(uint32(r))
}

func (s *OpCodeStream) SourceDynamicMemoryPosition(r opcode_sp_type.SourceDynamicMemoryPosition) {
	s.WriteUint32(uint32(r))
}

func (s *OpCodeStream) StackRange(r opcode_sp_type.StackRange) {
	s.WriteUint16(uint16(r))
}

func (s *OpCodeStream) SourceStackRange(r opcode_sp_type.SourceStackRange) {
	if uint(r) == 0 {
		panic("not allowed for it to be zero range")
	}
	s.WriteUint16(uint16(r))
}

func (s *OpCodeStream) SourceStackPositionRange(r opcode_sp_type.SourceStackPositionRange) {
	s.StackPosition(r.Position)
	s.SourceStackRange(r.Range)
}

func (s *OpCodeStream) TargetStackPosition(r opcode_sp_type.TargetStackPosition) {
	s.WriteUint32(uint32(r))
}

func (s *OpCodeStream) TargetFieldOffset(f opcode_sp_type.TargetFieldOffset) {
	s.WriteUint16(uint16(f))
}

func (s *OpCodeStream) Int32(v int32) {
	s.WriteUint32(uint32(v))
}

func (s *OpCodeStream) Rune(v instruction_sp.ShortRune) {
	s.Write(uint8(v))
}

func (s *OpCodeStream) Boolean(v bool) {
	value := 0
	if v {
		value = 1
	}
	s.Write(uint8(value))
}

func (s *OpCodeStream) EnumValue(v uint8) {
	s.Write(v)
}

func (s *OpCodeStream) programCounter() opcode_sp_type.ProgramCounter {
	return opcode_sp_type.NewProgramCounter(uint16(len(s.octets)))
}

func (s *OpCodeStream) DeltaPC(pc opcode_sp_type.DeltaPC) {
	s.Write(uint8(pc))
}

func (s *OpCodeStream) addLabelInject(inject *LabelInject) {
	s.DeltaPC(0xff)
	s.labelInjects = append(s.labelInjects, inject)
}

func (s *OpCodeStream) Label(l *opcode_sp_type.Label) {
	inject := NewLabelInject(l, s.programCounter())
	s.addLabelInject(inject)
}

func (s *OpCodeStream) LabelWithOffset(l *opcode_sp_type.Label, offset *opcode_sp_type.Label) {
	inject := NewLabelInjectWithOffset(l, s.programCounter(), offset)
	s.addLabelInject(inject)
}

func (s *OpCodeStream) Count(c int) {
	s.Write(uint8(c))
}

func (s *OpCodeStream) TypeIDConstant(c uint16) {
	s.Write(uint8(c >> 8))
	s.Write(uint8(c & 0xff))
}

func (s *OpCodeStream) Command(cmd instruction_sp.Commands) {
	s.Write(uint8(cmd))
}
