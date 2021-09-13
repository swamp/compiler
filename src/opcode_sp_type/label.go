/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package opcode_sp_type

import "fmt"

type Label struct {
	name      string
	pc        ProgramCounter
	isDefined bool
}

func NewLabelDefined(name string, pc ProgramCounter) *Label {
	return &Label{name: name, pc: pc, isDefined: true}
}

func NewLabel(name string) *Label {
	return &Label{name: name, isDefined: false}
}

func (l *Label) String() string {
	if l.name != "" {
		return fmt.Sprintf("[label %v %v]", l.name, l.pc)
	}

	return fmt.Sprintf("[label %v]", l.pc)
}

func (l *Label) IsDefined() bool {
	return l.isDefined
}

func (l *Label) DefinedProgramCounter() ProgramCounter {
	if !l.isDefined {
		panic("swamp opcodes: you can not read this label yet. Not defined")
	}

	return l.pc
}

func (l *Label) Define(pc ProgramCounter) {
	l.isDefined = true
	l.pc = pc
}
