/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package opcode_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/instruction_sp"
	"github.com/swamp/compiler/src/opcode_sp_type"
)

type Stream struct {
	instructions []Instruction
	allLabels    []*opcode_sp_type.Label
}

func NewStream() *Stream {
	return &Stream{}
}

func (s *Stream) addInstruction(c Instruction) {
	s.instructions = append(s.instructions, c)
}

func newLabel(name string, index int) *opcode_sp_type.Label {
	return opcode_sp_type.NewLabel(fmt.Sprintf("%v%v", name, index))
}

func (s *Stream) CreateLabel(name string) *opcode_sp_type.Label {
	l := newLabel(name, len(s.allLabels))

	s.allLabels = append(s.allLabels, l)
	return l
}

func (s *Stream) CreateList(destination opcode_sp_type.TargetStackPosition,
	itemSize opcode_sp_type.StackRange, arguments []opcode_sp_type.SourceStackPosition) *instruction_sp.CreateList {
	c := instruction_sp.NewCreateList(destination, itemSize, arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) CreateArray(destination opcode_sp_type.TargetStackPosition,
	itemSize opcode_sp_type.StackRange, arguments []opcode_sp_type.SourceStackPosition) *instruction_sp.CreateArray {
	c := instruction_sp.NewCreateArray(destination, itemSize, arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) UpdateStruct(target opcode_sp_type.TargetStackPosition, source opcode_sp_type.SourceStackPositionRange,
	fieldDefinitions []instruction_sp.CopyToFieldInfo) *instruction_sp.UpdateStruct {
	c := instruction_sp.NewUpdateStruct(target, source, fieldDefinitions)
	s.addInstruction(c)
	return c
}

func (s *Stream) EnumCase(source opcode_sp_type.SourceStackPosition,
	jumps []instruction_sp.EnumCaseJump) *instruction_sp.EnumCase {
	c := instruction_sp.NewEnumCase(source, jumps)
	s.addInstruction(c)

	return c
}

func (s *Stream) CasePatternMatching(source opcode_sp_type.SourceStackPositionRange,
	jumps []instruction_sp.CasePatternMatchingJump) *instruction_sp.CasePatternMatching {
	c := instruction_sp.NewCasePatternMatching(source, jumps)
	s.addInstruction(c)

	return c
}

func (s *Stream) TailCall() *instruction_sp.TailCall {
	c := instruction_sp.NewTailCall()
	s.addInstruction(c)
	return c
}

func (s *Stream) Call(newBasePointer opcode_sp_type.TargetStackPosition, function opcode_sp_type.SourceStackPosition,
) *instruction_sp.Call {
	c := instruction_sp.NewCall(newBasePointer, function)
	s.addInstruction(c)
	return c
}

func (s *Stream) Curry(target opcode_sp_type.TargetStackPosition, typeIDConstant uint16, function opcode_sp_type.SourceStackPosition,
	arguments []opcode_sp_type.SourceStackPosition) *instruction_sp.Curry {
	c := instruction_sp.NewCurry(target, typeIDConstant, function, arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) Return() *instruction_sp.Return {
	c := instruction_sp.NewReturn()
	s.addInstruction(c)
	return c
}

func (s *Stream) CallExternal(basePointer opcode_sp_type.SourceStackPosition, function opcode_sp_type.SourceStackPosition,
	arguments []opcode_sp_type.SourceStackPosition) *instruction_sp.CallExternal {
	c := instruction_sp.NewCallExternal(basePointer, function, arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) ListConj(destination opcode_sp_type.TargetStackPosition, list opcode_sp_type.SourceStackPosition,
	item opcode_sp_type.SourceStackPosition) *instruction_sp.ListConj {
	c := instruction_sp.NewListConj(destination, item, list)
	s.addInstruction(c)
	return c
}

func (s *Stream) Jump(delta *opcode_sp_type.Label) *instruction_sp.Jump {
	c := instruction_sp.NewJump(delta)
	s.addInstruction(c)
	return c
}

func (s *Stream) Label(label *opcode_sp_type.Label) error {
	c := NewVirtualLabel(label)
	s.addInstruction(c)
	return nil
}

func (s *Stream) BranchFalse(test opcode_sp_type.SourceStackPosition, jump *opcode_sp_type.Label) *instruction_sp.BranchFalse {
	c := instruction_sp.NewBranchFalse(test, jump)
	s.addInstruction(c)
	return c
}

func (s *Stream) BranchTrue(test opcode_sp_type.SourceStackPosition, jump *opcode_sp_type.Label) *instruction_sp.BranchTrue {
	c := instruction_sp.NewBranchTrue(test, jump)
	s.addInstruction(c)
	return c
}

func (s *Stream) BinaryOperator(destination opcode_sp_type.TargetStackPosition,
	operatorType instruction_sp.BinaryOperatorType, a opcode_sp_type.SourceStackPosition,
	b opcode_sp_type.SourceStackPosition) *instruction_sp.BinaryOperator {
	opcode := instruction_sp.BinaryOperatorToOpCode(operatorType)
	c := instruction_sp.NewBinaryOperator(opcode, destination, a, b)
	s.addInstruction(c)
	return c
}

func (s *Stream) ListAppend(destination opcode_sp_type.TargetStackPosition, a opcode_sp_type.SourceStackPosition,
	b opcode_sp_type.SourceStackPosition) *instruction_sp.ListAppend {
	c := instruction_sp.NewListAppend(destination, a, b)
	s.addInstruction(c)
	return c
}

func (s *Stream) StringAppend(destination opcode_sp_type.TargetStackPosition, a opcode_sp_type.SourceStackPosition,
	b opcode_sp_type.SourceStackPosition) *instruction_sp.StringAppend {
	c := instruction_sp.NewStringAppend(destination, a, b)
	s.addInstruction(c)
	return c
}

func (s *Stream) IntUnaryOperator(destination opcode_sp_type.TargetStackPosition, operatorType instruction_sp.UnaryOperatorType,
	a opcode_sp_type.SourceStackPosition) *instruction_sp.IntUnaryOperator {
	opcode := instruction_sp.UnaryOperatorToOpCode(operatorType)
	c := instruction_sp.NewIntUnaryOperator(opcode, destination, a)
	s.addInstruction(c)
	return c
}

func (s *Stream) Serialize() ([]byte, error) {
	writer := NewOpCodeStream()

	for _, instruction := range s.instructions {
		lbl, _ := instruction.(*VirtualLabel)
		if lbl != nil {
			lbl.Label().Define(writer.programCounter())
		} else {
			instruction.Write(writer)
		}
	}

	for _, label := range s.allLabels {
		if !label.IsDefined() {
			return nil, fmt.Errorf("Label %v not defined", label)
		}
	}

	beforeOctets := writer.Octets()
	block := NewOctetBlock(beforeOctets)

	fixupErr := block.FixUpLabelInjects(writer.LabelInjects())
	if fixupErr != nil {
		return nil, fixupErr
	}

	fixedUpOctets := block.Octets()

	return fixedUpOctets, nil
}
