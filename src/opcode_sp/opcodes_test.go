/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package opcode_sp_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	swampopcodeinst "github.com/swamp/opcodes/instruction"
	swampopcode "github.com/swamp/opcodes/opcode"
	swampopcodetype "github.com/swamp/opcodes/type"
)

func checkEqual(t *testing.T, expected []byte, actual []byte) {
	diff := bytes.Compare(expected, actual)
	if diff != 0 {
		t.Errorf("Octets not equal. Expected %v but got %v", hex.EncodeToString(expected), hex.EncodeToString(actual))
	}
}

func TestSerialize(t *testing.T) {
	s := swampopcode.NewStream()

	destination := swampopcodetype.NewRegister(0x44)
	source := swampopcodetype.NewRegister(0x11)
	source2 := swampopcodetype.NewRegister(0x13)
	arguments := []swampopcodetype.Register{source, source2}
	s.CreateStruct(destination, arguments)
	octets, _ := s.Serialize()
	checkEqual(t, []byte{uint8(swampopcodeinst.CmdCreateStruct), 0x44, 0x02, 0x11, 0x13}, octets)
}

func TestLabel(t *testing.T) {
	s := swampopcode.NewStream()

	label := s.CreateLabel("jumphere")
	s.Jump(label)
	s.Return()
	label.Define(swampopcodetype.NewProgramCounter(uint16(0x03)))

	octets, _ := s.Serialize()
	checkEqual(t, []byte{swampopcodeinst.CmdJump, 0x01, swampopcodeinst.CmdReturn}, octets)
}
