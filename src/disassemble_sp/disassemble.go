/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package swampdisasm_sp

import (
	"encoding/binary"
	"fmt"

	"github.com/swamp/compiler/src/instruction_sp"
	"github.com/swamp/compiler/src/opcode_sp"
	"github.com/swamp/compiler/src/opcode_sp_type"
)

type Register struct {
	id uint8
}

type Argument interface {
	String() string
}

type OpcodeInStream struct {
	position int
	octets   []byte
}

func NewOpcodeInStream(octets []byte) *OpcodeInStream {
	return &OpcodeInStream{octets: octets}
}

func (s *OpcodeInStream) IsEOF() bool {
	return s.position >= len(s.octets)
}

func (s *OpcodeInStream) readUint8() uint8 {
	if s.position == len(s.octets) {
		panic("swamp disassembler: read too far")
	}

	a := s.octets[s.position]

	s.position++

	return a
}

func (s *OpcodeInStream) readUint16() uint16 {
	if s.position+2 == len(s.octets) {
		panic("swamp disassembler: read too far uint16")
	}

	pointer := binary.LittleEndian.Uint16(s.octets[s.position : s.position+2])

	s.position += 2

	return pointer
}

func (s *OpcodeInStream) readUint32() uint32 {
	if s.position+4 == len(s.octets) {
		panic("swamp disassembler: read too far uint32")
	}

	pointer := binary.LittleEndian.Uint32(s.octets[s.position : s.position+4])

	s.position += 4

	return pointer
}

func (s *OpcodeInStream) readCommand() instruction_sp.Commands {
	return instruction_sp.Commands(s.readUint8())
}

func (s *OpcodeInStream) programCounter() opcode_sp_type.ProgramCounter {
	return opcode_sp_type.NewProgramCounter(uint16(s.position))
}

func (s *OpcodeInStream) readTypeIDConstant() uint16 {
	return s.readUint16()
}

func (s *OpcodeInStream) readCount() int {
	return int(s.readUint8())
}

func (s *OpcodeInStream) readArgOffsetSize() opcode_sp_type.ArgOffsetSize {
	return opcode_sp_type.ArgOffsetSize{
		Offset: s.readUint16(),
		Size:   s.readUint16(),
	}
}

func (s *OpcodeInStream) readInt32() int32 {
	return int32(s.readUint32())
}

func (s *OpcodeInStream) readBoolean() bool {
	return s.readUint8() != 0
}

func (s *OpcodeInStream) readItemSize() opcode_sp_type.StackRange {
	return opcode_sp_type.StackRange(s.readUint16())
}

func (s *OpcodeInStream) readAlign() opcode_sp_type.MemoryAlign {
	return opcode_sp_type.MemoryAlign(s.readUint8())
}

func (s *OpcodeInStream) readLabel() *opcode_sp_type.Label {
	delta := s.readUint16()
	resultingPosition := s.programCounter().Add(delta)

	return opcode_sp_type.NewLabelDefined("", resultingPosition)
}

func (s *OpcodeInStream) readLabelOffset(offset opcode_sp_type.ProgramCounter) *opcode_sp_type.Label {
	delta := s.readUint16()
	resultingPosition := offset.Add(delta)

	return opcode_sp_type.NewLabelDefined("offset", resultingPosition)
}

func (s *OpcodeInStream) readSourceStackPosition() opcode_sp_type.SourceStackPosition {
	pointer := s.readUint32()
	return opcode_sp_type.SourceStackPosition(pointer)
}

func (s *OpcodeInStream) readSourceStackPositionRange() opcode_sp_type.SourceStackPositionRange {
	pointer := s.readUint32()
	size := s.readUint16()
	if size == 0 {
		panic("disassemble: we can not allow zero size in range")
	}
	return opcode_sp_type.SourceStackPositionRange{
		Position: opcode_sp_type.SourceStackPosition(pointer),
		Range:    opcode_sp_type.SourceStackRange(size),
	}
}

func (s *OpcodeInStream) readSourceStackPositions() []opcode_sp_type.SourceStackPosition {
	count := s.readCount()
	targetArray := make([]opcode_sp_type.SourceStackPosition, count)
	for i := 0; i < count; i++ {
		targetArray[i] = s.readSourceStackPosition()
	}

	return targetArray
}

func (s *OpcodeInStream) readTargetStackPosition() opcode_sp_type.TargetStackPosition {
	pointer := s.readUint32()
	return opcode_sp_type.TargetStackPosition(pointer)
}

func (s *OpcodeInStream) readSourceDynamicMemoryPosition() opcode_sp_type.SourceDynamicMemoryPosition {
	pointer := s.readUint32()
	return opcode_sp_type.SourceDynamicMemoryPosition(pointer)
}

func disassembleListConj(s *OpcodeInStream) *instruction_sp.ListConj {
	destination := s.readTargetStackPosition()
	list := s.readSourceStackPosition()
	item := s.readSourceStackPosition()

	return instruction_sp.NewListConj(destination, item, list)
}

func disassembleListAppend(s *OpcodeInStream) *instruction_sp.ListAppend {
	destination := s.readTargetStackPosition()
	a := s.readSourceStackPosition()
	b := s.readSourceStackPosition()

	return instruction_sp.NewListAppend(destination, a, b)
}

func disassembleStringAppend(s *OpcodeInStream) *instruction_sp.StringAppend {
	destination := s.readTargetStackPosition()
	a := s.readSourceStackPosition()
	b := s.readSourceStackPosition()

	return instruction_sp.NewStringAppend(destination, a, b)
}

func disassembleBinaryOperator(cmd instruction_sp.Commands, s *OpcodeInStream) *instruction_sp.BinaryOperator {
	destination := s.readTargetStackPosition()
	a := s.readSourceStackPosition()
	b := s.readSourceStackPosition()

	return instruction_sp.NewBinaryOperator(cmd, destination, a, b)
}

func disassembleStringBinaryOperator(cmd instruction_sp.Commands, s *OpcodeInStream) *instruction_sp.BinaryOperator {
	destination := s.readTargetStackPosition()
	a := s.readSourceStackPosition()
	b := s.readSourceStackPosition()

	return instruction_sp.NewBinaryOperator(cmd, destination, a, b)
}

func disassembleEnumBinaryOperator(cmd instruction_sp.Commands, s *OpcodeInStream) *instruction_sp.BinaryOperator {
	destination := s.readTargetStackPosition()
	a := s.readSourceStackPosition()
	b := s.readSourceStackPosition()

	return instruction_sp.NewBinaryOperator(cmd, destination, a, b)
}

func disassembleBitwiseOperator(cmd instruction_sp.Commands, s *OpcodeInStream) *instruction_sp.BinaryOperator {
	destination := s.readTargetStackPosition()
	a := s.readSourceStackPosition()
	b := s.readSourceStackPosition()

	return instruction_sp.NewBinaryOperator(cmd, destination, a, b)
}

func disassembleBitwiseUnaryOperator(cmd instruction_sp.Commands, s *OpcodeInStream) *instruction_sp.IntUnaryOperator {
	destination := s.readTargetStackPosition()
	a := s.readSourceStackPosition()

	return instruction_sp.NewIntUnaryOperator(cmd, destination, a)
}

func disassembleLoadInteger(s *OpcodeInStream) *instruction_sp.LoadInteger {
	destination := s.readTargetStackPosition()
	a := s.readInt32()

	return instruction_sp.NewLoadInteger(destination, a)
}

func disassembleLoadRune(s *OpcodeInStream) *instruction_sp.LoadRune {
	destination := s.readTargetStackPosition()
	shortRune := s.readUint8()

	return instruction_sp.NewLoadRune(destination, instruction_sp.ShortRune(shortRune))
}

func disassembleLoadBoolean(s *OpcodeInStream) *instruction_sp.LoadBool {
	destination := s.readTargetStackPosition()
	a := s.readBoolean()

	return instruction_sp.NewLoadBool(destination, a)
}

func disassembleSetEnum(s *OpcodeInStream) *instruction_sp.SetEnum {
	destination := s.readTargetStackPosition()
	a := s.readUint8()

	return instruction_sp.NewSetEnum(destination, a)
}

func disassembleLoadZeroMemoryPointer(s *OpcodeInStream) *instruction_sp.LoadZeroMemoryPointer {
	destination := s.readTargetStackPosition()
	source := s.readSourceDynamicMemoryPosition()

	return instruction_sp.NewLoadZeroMemoryPointer(destination, source)
}

func disassembleCreateList(s *OpcodeInStream) *instruction_sp.CreateList {
	destination := s.readTargetStackPosition()
	itemSize := s.readItemSize()
	memoryAlign := s.readAlign()
	arguments := s.readSourceStackPositions()

	return instruction_sp.NewCreateList(destination, itemSize, memoryAlign, arguments)
}

func disassembleCreateArray(s *OpcodeInStream) *instruction_sp.CreateArray {
	destination := s.readTargetStackPosition()
	itemSize := s.readItemSize()
	memoryAlign := s.readAlign()
	arguments := s.readSourceStackPositions()

	return instruction_sp.NewCreateArray(destination, itemSize, memoryAlign, arguments)
}

func disassembleCall(s *OpcodeInStream) *instruction_sp.Call {
	newStackPointer := s.readTargetStackPosition()
	functionRegister := s.readSourceStackPosition()

	return instruction_sp.NewCall(newStackPointer, functionRegister)
}

func disassembleCallExternal(s *OpcodeInStream) *instruction_sp.CallExternal {
	newStackPointer := s.readTargetStackPosition()
	functionRegister := s.readSourceStackPosition()

	return instruction_sp.NewCallExternal(newStackPointer, functionRegister)
}

func disassembleCallExternalWithSizes(s *OpcodeInStream) *instruction_sp.CallExternalWithSizes {
	newStackPointer := s.readTargetStackPosition()
	functionRegister := s.readSourceStackPosition()
	count := s.readCount()
	targetArgs := make([]opcode_sp_type.ArgOffsetSize, count)
	for i := 0; i < count; i++ {
		targetArgs[i] = s.readArgOffsetSize()
	}

	return instruction_sp.NewCallExternalWithSizes(newStackPointer, functionRegister, targetArgs)
}

func disassembleCurry(s *OpcodeInStream) *instruction_sp.Curry {
	destination := s.readTargetStackPosition()
	typeIDConstant := s.readTypeIDConstant()
	functionRegister := s.readSourceStackPosition()
	arguments := s.readSourceStackPositionRange()

	return instruction_sp.NewCurry(destination, typeIDConstant, functionRegister, arguments)
}

func disassembleEnumCase(s *OpcodeInStream) *instruction_sp.EnumCase {
	source := s.readSourceStackPosition()
	count := s.readCount()

	var jumps []instruction_sp.EnumCaseJump

	var lastLabel *opcode_sp_type.Label

	for i := 0; i < count; i++ {
		enumValue := s.readUint8()

		var label *opcode_sp_type.Label

		if lastLabel != nil {
			label = s.readLabelOffset(lastLabel.DefinedProgramCounter())
		} else {
			label = s.readLabel()
		}

		lastLabel = label
		jump := instruction_sp.NewEnumCaseJump(enumValue, label)
		jumps = append(jumps, jump)
	}

	return instruction_sp.NewEnumCase(source, jumps)
}

/*
	var matchingType instruction_sp.PatternMatchingType
	switch cmd {
	case instruction_sp.CmdPatternMatchingInt:
		matchingType = instruction_sp.PatternMatchingTypeInt
	case instruction_sp.CmdPatternMatchingString:
		matchingType = instruction_sp.PatternMatchingTypeString
	default:
		panic(fmt.Errorf("unknown matching type %v", cmd))
	}
*/
func disassemblePatternMatchingInt(cmd instruction_sp.Commands, s *OpcodeInStream) *instruction_sp.PatternMatchingInt {
	source := s.readSourceStackPosition()
	count := s.readCount()

	var jumps []instruction_sp.EnumCasePatternMatchingIntJump

	var lastLabel *opcode_sp_type.Label

	for i := 0; i < count; i++ {
		matchInteger := s.readInt32()

		var label *opcode_sp_type.Label

		if lastLabel != nil {
			label = s.readLabelOffset(lastLabel.DefinedProgramCounter())
		} else {
			label = s.readLabel()
		}

		lastLabel = label
		jump := instruction_sp.NewEnumCasePatternMatchingIntJump(matchInteger, label)
		jumps = append(jumps, jump)
	}

	defaultLabel := s.readLabelOffset(lastLabel.DefinedProgramCounter())

	return instruction_sp.NewPatternMatchingInt(source, jumps, defaultLabel)
}

func disassembleMemoryCopy(s *OpcodeInStream) *instruction_sp.MemoryCopy {
	destination := s.readTargetStackPosition()
	source := s.readSourceStackPositionRange()

	return instruction_sp.NewMemoryCopy(destination, source)
}

func disassembleTailCall(s *OpcodeInStream) *instruction_sp.TailCall {
	return nil
}

func disassembleReturn(s *OpcodeInStream) *instruction_sp.Return {
	return instruction_sp.NewReturn()
}

func disassembleJump(s *OpcodeInStream) *instruction_sp.Jump {
	label := s.readLabel()

	return instruction_sp.NewJump(label)
}

func disassembleBranchFalse(s *OpcodeInStream) *instruction_sp.BranchFalse {
	test := s.readSourceStackPosition()
	label := s.readLabel()

	return instruction_sp.NewBranchFalse(test, label)
}

func disassembleBranchTrue(s *OpcodeInStream) *instruction_sp.BranchTrue {
	test := s.readSourceStackPosition()
	label := s.readLabel()

	return instruction_sp.NewBranchTrue(test, label)
}

func decodeOpcode(cmd instruction_sp.Commands, s *OpcodeInStream) opcode_sp.Instruction {
	switch cmd {
	case instruction_sp.CmdIntAdd:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdIntSub:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdIntDiv:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdIntMul:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdIntEqual:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdIntNotEqual:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdIntLess:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdIntLessOrEqual:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdIntGreater:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdIntGreaterOrEqual:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdFixedDiv:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdFixedMul:
		return disassembleBinaryOperator(cmd, s)
	case instruction_sp.CmdListConj:
		return disassembleListConj(s)
	case instruction_sp.CmdListAppend:
		return disassembleListAppend(s)
	case instruction_sp.CmdStringAppend:
		return disassembleStringAppend(s)
	case instruction_sp.CmdCreateList:
		return disassembleCreateList(s)
	case instruction_sp.CmdCreateArray:
		return disassembleCreateArray(s)
	case instruction_sp.CmdEnumCase:
		return disassembleEnumCase(s)
	case instruction_sp.CmdPatternMatchingInt:
		return disassemblePatternMatchingInt(cmd, s)
	case instruction_sp.CmdPatternMatchingString:
		panic("not implemented")
	case instruction_sp.CmdCopyMemory:
		return disassembleMemoryCopy(s)
	case instruction_sp.CmdCall:
		return disassembleCall(s)
	case instruction_sp.CmdCallExternal:
		return disassembleCallExternal(s)
	case instruction_sp.CmdCallExternalWithSizes:
		return disassembleCallExternalWithSizes(s)
	case instruction_sp.CmdTailCall:
		return disassembleTailCall(s)
	case instruction_sp.CmdCurry:
		return disassembleCurry(s)
	case instruction_sp.CmdReturn:
		return disassembleReturn(s)
	case instruction_sp.CmdJump:
		return disassembleJump(s)
	case instruction_sp.CmdBranchFalse:
		return disassembleBranchFalse(s)
	case instruction_sp.CmdBranchTrue:
		return disassembleBranchTrue(s)
	case instruction_sp.CmdIntBitwiseAnd:
		return disassembleBitwiseOperator(cmd, s)
	case instruction_sp.CmdIntBitwiseOr:
		return disassembleBitwiseOperator(cmd, s)
	case instruction_sp.CmdIntBitwiseXor:
		return disassembleBitwiseOperator(cmd, s)
	case instruction_sp.CmdIntBitwiseNot:
		return disassembleBitwiseUnaryOperator(cmd, s)
	case instruction_sp.CmdBoolLogicalNot:
		return disassembleBitwiseUnaryOperator(cmd, s)
	case instruction_sp.CmdIntNegate:
		return disassembleBitwiseUnaryOperator(cmd, s)
	case instruction_sp.CmdLoadInteger:
		return disassembleLoadInteger(s)
	case instruction_sp.CmdLoadRune:
		return disassembleLoadRune(s)
	case instruction_sp.CmdLoadBoolean:
		return disassembleLoadBoolean(s)
	case instruction_sp.CmdLoadZeroMemoryPointer:
		return disassembleLoadZeroMemoryPointer(s)
	case instruction_sp.CmdSetEnum:
		return disassembleSetEnum(s)
	case instruction_sp.CmdStringEqual:
		return disassembleStringBinaryOperator(cmd, s)
	case instruction_sp.CmdStringNotEqual:
		return disassembleStringBinaryOperator(cmd, s)
	case instruction_sp.CmdEnumEqual:
		return disassembleEnumBinaryOperator(cmd, s)
	case instruction_sp.CmdEnumNotEqual:
		return disassembleEnumBinaryOperator(cmd, s)
	}

	panic(fmt.Sprintf("swamp disassembler: unknown opcode:%v", cmd))

	// return nil
}

func Disassemble(octets []byte) []string {
	var lines []string

	s := NewOpcodeInStream(octets)

	for !s.IsEOF() {
		startPc := s.programCounter()
		cmd := s.readCommand()

		// log.Printf("disasembling :%s (%02x)\n", instruction_sp.OpcodeToMnemonic(cmd), cmd)
		args := decodeOpcode(cmd, s)
		line := fmt.Sprintf("%04x: %v", startPc.Value(), args)
		lines = append(lines, line)
	}

	return lines
}
