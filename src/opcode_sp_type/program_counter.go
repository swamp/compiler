/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package opcode_sp_type

import "fmt"

type DeltaPC uint16

type ProgramCounter struct {
	position uint16
}

func (p ProgramCounter) String() string {
	return fmt.Sprintf("@%04x", p.position)
}

func NewProgramCounter(position uint16) ProgramCounter {
	return ProgramCounter{position: position}
}

func (p ProgramCounter) Add(delta uint16) ProgramCounter {
	return ProgramCounter{position: p.position + delta}
}

func (p ProgramCounter) Value() uint16 {
	return p.position
}

func (p ProgramCounter) IsAfter(other ProgramCounter) bool {
	return p.position > other.position
}

func (p ProgramCounter) Delta(after ProgramCounter) (DeltaPC, error) {
	delta := int(after.position) - int(p.position)

	if delta < 0 || delta > 0xffff {
		return DeltaPC(0), fmt.Errorf("illegal jump. Jumping forward %d, but maximum is 65535. Please split your code into separate functions.", delta)
	}

	return DeltaPC(delta), nil
}
