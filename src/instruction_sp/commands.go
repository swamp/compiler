/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"
)

type Commands uint8

const (
	// Branching
	CmdEnumCase    Commands = 0x01
	CmdBranchFalse Commands = 0x02
	CmdBranchTrue  Commands = 0x03
	CmdJump        Commands = 0x04

	// Call
	CmdCall         Commands = 0x05
	CmdReturn       Commands = 0x06
	CmdCallExternal Commands = 0x07
	CmdTailCall     Commands = 0x08
	CmdCurry        Commands = 0x09

	// Arithmetic
	CmdIntAdd    Commands = 0x0a
	CmdIntSub    Commands = 0x0b
	CmdIntMul    Commands = 0x0c
	CmdIntDiv    Commands = 0x0d
	CmdIntNegate Commands = 0x0e

	// Arithmetic fixed point
	CmdFixedMul Commands = 0x0f
	CmdFixedDiv Commands = 0x10

	// Boolean operators
	CmdIntEqual          Commands = 0x11
	CmdIntNotEqual       Commands = 0x12
	CmdIntLess           Commands = 0x13
	CmdIntLessOrEqual    Commands = 0x14
	CmdIntGreater        Commands = 0x15
	CmdIntGreaterOrEqual Commands = 0x16

	// Boolean
	CmdBoolLogicalNot Commands = 0x17

	// Boolean strings
	CmdStringEqual    Commands = 0x18
	CmdStringNotEqual Commands = 0x19

	// Bitwise operators
	CmdIntBitwiseAnd Commands = 0x1a
	CmdIntBitwiseOr  Commands = 0x1b
	CmdIntBitwiseXor Commands = 0x1c
	CmdIntBitwiseNot Commands = 0x1d

	// Creating dynamic structures
	CmdCreateList  Commands = 0x1e
	CmdCreateArray Commands = 0x1f

	// Append
	CmdListConj     Commands = 0x20
	CmdListAppend   Commands = 0x21
	CmdStringAppend Commands = 0x22

	// Load
	CmdLoadInteger           Commands = 0x23
	CmdLoadBoolean           Commands = 0x24
	CmdLoadRune              Commands = 0x25
	CmdLoadZeroMemoryPointer Commands = 0x26
	CmdCopyMemory            Commands = 0x27
	CmdSetEnum               Commands = 0x28
	CmdCallExternalWithSizes Commands = 0x29

	// enum operator
	CmdEnumEqual    Commands = 0x2a
	CmdEnumNotEqual Commands = 0x2b
)

func OpcodeToMnemonic(cmd Commands) string {
	names := map[Commands]string{
		CmdListConj:              "lconj",
		CmdEnumCase:              "jmpe",
		CmdBranchFalse:           "bne",
		CmdJump:                  "jmp",
		CmdCall:                  "call",
		CmdCallExternalWithSizes: "callvar",
		CmdReturn:                "ret",
		CmdCallExternal:          "ecall",
		CmdTailCall:              "tcall",
		CmdIntAdd:                "addi",
		CmdIntSub:                "subi",
		CmdIntMul:                "muli",
		CmdIntDiv:                "divi",
		CmdIntEqual:              "cpeqi",
		CmdIntNotEqual:           "cpnei",
		CmdIntLess:               "cplti",
		CmdIntLessOrEqual:        "cplei",
		CmdIntGreater:            "cpgti",
		CmdIntGreaterOrEqual:     "cpgei",
		CmdStringEqual:           "cpeqs",
		CmdStringNotEqual:        "cpnes",
		CmdEnumEqual:             "cpeqe",
		CmdEnumNotEqual:          "cpnee",
		CmdIntBitwiseAnd:         "andi",
		CmdIntBitwiseOr:          "ori",
		CmdIntBitwiseXor:         "xori",
		CmdIntBitwiseNot:         "noti",
		CmdBoolLogicalNot:        "not",
		CmdBranchTrue:            "brt",
		CmdCurry:                 "curry",
		CmdCreateList:            "crl",
		CmdListAppend:            "concatl",
		CmdCopyMemory:            "cpy",
		CmdStringAppend:          "concats",
		CmdFixedMul:              "fxmul",
		CmdFixedDiv:              "fxdiv",
		CmdIntNegate:             "ineg",
		CmdCreateArray:           "cra",
		CmdLoadInteger:           "ldi",
		CmdLoadBoolean:           "ldb",
		CmdLoadZeroMemoryPointer: "ldz",
		CmdSetEnum:               "lde",
		CmdLoadRune:              "ldr",
	}

	mnemonic, found := names[cmd]
	if !found {
		panic(fmt.Errorf("no lookup for %v", cmd))
	}

	return mnemonic
}
