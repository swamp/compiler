/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package swampspopcode

import (
	"fmt"

	swampopcodeinst "github.com/swamp/opcodes/instruction"
	swampopcodetype "github.com/swamp/opcodes/type"
)

type Stream struct {
	instructions []Instruction
	allLabels    []*swampopcodetype.Label
}

func NewStream() *Stream {
	return &Stream{}
}

func (s *Stream) addInstruction(c Instruction) {
	s.instructions = append(s.instructions, c)
}

func newLabel(name string, index int) *swampopcodetype.Label {
	return swampopcodetype.NewLabel(fmt.Sprintf("%v%v", name, index))
}

func (s *Stream) CreateLabel(name string) *swampopcodetype.Label {
	l := newLabel(name, len(s.allLabels))

	s.allLabels = append(s.allLabels, l)
	return l
}

func (s *Stream) CreateStruct(destination swampopcodetype.Register,
	arguments []swampopcodetype.Register) *swampopcodeinst.CreateStruct {
	c := swampopcodeinst.NewCreateStruct(destination, arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) StructSplit(source swampopcodetype.Register,
	destinations []swampopcodetype.Register) *swampopcodeinst.StructSplit {
	c := swampopcodeinst.NewStructSplit(source, destinations)
	s.addInstruction(c)
	return c
}

func (s *Stream) CreateList(destination swampopcodetype.Register,
	arguments []swampopcodetype.Register) *swampopcodeinst.CreateList {
	c := swampopcodeinst.NewCreateList(destination, arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) UpdateStruct(target swampopcodetype.Register, source swampopcodetype.Register,
	fieldDefinitions []swampopcodeinst.CopyToFieldInfo) *swampopcodeinst.UpdateStruct {
	c := swampopcodeinst.NewUpdateStruct(target, source, fieldDefinitions)
	s.addInstruction(c)
	return c
}

func (s *Stream) GetStruct(target swampopcodetype.Register, source swampopcodetype.Register,
	fieldLookup []swampopcodetype.Field) *swampopcodeinst.GetStruct {
	c := swampopcodeinst.NewGetStruct(target, source, fieldLookup)
	s.addInstruction(c)
	return c
}

func (s *Stream) EnumCase(source swampopcodetype.Register,
	jumps []swampopcodeinst.EnumCaseJump) *swampopcodeinst.EnumCase {
	c := swampopcodeinst.NewEnumCase(source, jumps)
	s.addInstruction(c)

	return c
}

func (s *Stream) CasePatternMatching(source swampopcodetype.Register,
	jumps []swampopcodeinst.CasePatternMatchingJump) *swampopcodeinst.CasePatternMatching {
	c := swampopcodeinst.NewCasePatternMatching(source, jumps)
	s.addInstruction(c)

	return c
}

func (s *Stream) RegCopy(target swampopcodetype.Register, source swampopcodetype.Register) *swampopcodeinst.RegCopy {
	c := swampopcodeinst.NewRegCopy(target, source)
	s.addInstruction(c)

	return c
}

func (s *Stream) TailCall(arguments []swampopcodetype.Register) *swampopcodeinst.TailCall {
	c := swampopcodeinst.NewTailCall(arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) Call(target swampopcodetype.Register, function swampopcodetype.Register,
	arguments []swampopcodetype.Register) *swampopcodeinst.Call {
	c := swampopcodeinst.NewCall(target, function, arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) Curry(target swampopcodetype.Register, typeIDConstant uint16, function swampopcodetype.Register,
	arguments []swampopcodetype.Register) *swampopcodeinst.Curry {
	c := swampopcodeinst.NewCurry(target, typeIDConstant, function, arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) Enum(target swampopcodetype.Register, enumFieldIndex int,
	arguments []swampopcodetype.Register) *swampopcodeinst.Enum {
	c := swampopcodeinst.NewEnum(target, enumFieldIndex, arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) Return() *swampopcodeinst.Return {
	c := swampopcodeinst.NewReturn()
	s.addInstruction(c)
	return c
}

func (s *Stream) CallExternal(destination swampopcodetype.Register, function swampopcodetype.Register,
	arguments []swampopcodetype.Register) *swampopcodeinst.CallExternal {
	c := swampopcodeinst.NewCallExternal(destination, function, arguments)
	s.addInstruction(c)
	return c
}

func (s *Stream) ListConj(destination swampopcodetype.Register, list swampopcodetype.Register,
	item swampopcodetype.Register) *swampopcodeinst.ListConj {
	c := swampopcodeinst.NewListConj(destination, item, list)
	s.addInstruction(c)
	return c
}

func (s *Stream) Jump(delta *swampopcodetype.Label) *swampopcodeinst.Jump {
	c := swampopcodeinst.NewJump(delta)
	s.addInstruction(c)
	return c
}

func (s *Stream) Label(label *swampopcodetype.Label) error {
	c := NewVirtualLabel(label)
	s.addInstruction(c)
	return nil
}

func (s *Stream) BranchFalse(test swampopcodetype.Register, jump *swampopcodetype.Label) *swampopcodeinst.BranchFalse {
	c := swampopcodeinst.NewBranchFalse(test, jump)
	s.addInstruction(c)
	return c
}

func (s *Stream) BranchTrue(test swampopcodetype.Register, jump *swampopcodetype.Label) *swampopcodeinst.BranchTrue {
	c := swampopcodeinst.NewBranchTrue(test, jump)
	s.addInstruction(c)
	return c
}

func (s *Stream) BinaryOperator(destination swampopcodetype.Register,
	operatorType swampopcodeinst.BinaryOperatorType, a swampopcodetype.Register,
	b swampopcodetype.Register) *swampopcodeinst.BinaryOperator {
	opcode := swampopcodeinst.BinaryOperatorToOpCode(operatorType)
	c := swampopcodeinst.NewBinaryOperator(opcode, destination, a, b)
	s.addInstruction(c)
	return c
}

func (s *Stream) ListAppend(destination swampopcodetype.Register, a swampopcodetype.Register,
	b swampopcodetype.Register) *swampopcodeinst.ListAppend {
	c := swampopcodeinst.NewListAppend(destination, a, b)
	s.addInstruction(c)
	return c
}

func (s *Stream) StringAppend(destination swampopcodetype.Register, a swampopcodetype.Register,
	b swampopcodetype.Register) *swampopcodeinst.StringAppend {
	c := swampopcodeinst.NewStringAppend(destination, a, b)
	s.addInstruction(c)
	return c
}

func (s *Stream) IntUnaryOperator(destination swampopcodetype.Register, operatorType swampopcodeinst.UnaryOperatorType,
	a swampopcodetype.Register) *swampopcodeinst.IntUnaryOperator {
	opcode := swampopcodeinst.UnaryOperatorToOpCode(operatorType)
	c := swampopcodeinst.NewIntUnaryOperator(opcode, destination, a)
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
