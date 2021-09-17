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
	CmdEnumCase    Commands = 0x06
	CmdBranchFalse Commands = 0x07
	CmdBranchTrue  Commands = 0x1c
	CmdJump        Commands = 0x08

	// Call
	CmdCall         Commands = 0x09
	CmdReturn       Commands = 0x0a
	CmdCallExternal Commands = 0x0b
	CmdTailCall     Commands = 0x0c
	CmdCurry        Commands = 0x20

	// Arithmetic
	CmdIntAdd    Commands = 0x0d
	CmdIntSub    Commands = 0x0e
	CmdIntMul    Commands = 0x0f
	CmdIntDiv    Commands = 0x10
	CmdIntNegate Commands = 0x27

	// Arithmetic fixed point
	CmdFixedMul Commands = 0x25
	CmdFixedDiv Commands = 0x26

	// Boolean operators
	CmdIntEqual          Commands = 0x11
	CmdIntNotEqual       Commands = 0x12
	CmdIntLess           Commands = 0x13
	CmdIntLessOrEqual    Commands = 0x14
	CmdIntGreater        Commands = 0x15
	CmdIntGreaterOrEqual Commands = 0x16

	// Boolean
	CmdBoolLogicalNot Commands = 0x1b

	// Boolean strings
	CmdStringEqual    Commands = 0x35
	CmdStringNotEqual Commands = 0x36

	// Bitwise operators
	CmdIntBitwiseAnd Commands = 0x17
	CmdIntBitwiseOr  Commands = 0x18
	CmdIntBitwiseXor Commands = 0x19
	CmdIntBitwiseNot Commands = 0x1a

	// Creating dynamic structures
	CmdCreateList  Commands = 0x21
	CmdCreateArray Commands = 0x29

	// Append
	CmdListConj     Commands = 0x05
	CmdListAppend   Commands = 0x22
	CmdStringAppend Commands = 0x24

	// Load
	CmdLoadInteger           Commands = 0x31
	CmdLoadBoolean           Commands = 0x32
	CmdLoadRune              Commands = 0x37
	CmdLoadZeroMemoryPointer Commands = 0x33
	CmdCopyMemory            Commands = 0x30
	CmdSetEnum               Commands = 0x34
)

func OpcodeToMnemonic(cmd Commands) string {
	names := map[Commands]string{
		CmdListConj:              "conj",
		CmdEnumCase:              "case",
		CmdBranchFalse:           "bne",
		CmdJump:                  "jmp",
		CmdCall:                  "call",
		CmdReturn:                "ret",
		CmdCallExternal:          "ecall",
		CmdTailCall:              "tcl",
		CmdIntAdd:                "add",
		CmdIntSub:                "sub",
		CmdIntMul:                "mul",
		CmdIntDiv:                "div",
		CmdIntEqual:              "cpieq",
		CmdIntNotEqual:           "cpine",
		CmdIntLess:               "cpl",
		CmdIntLessOrEqual:        "cple",
		CmdIntGreater:            "cpg",
		CmdIntGreaterOrEqual:     "cpge",
		CmdIntBitwiseAnd:         "band",
		CmdIntBitwiseOr:          "bor",
		CmdIntBitwiseXor:         "bxor",
		CmdIntBitwiseNot:         "bnot",
		CmdBoolLogicalNot:        "not",
		CmdBranchTrue:            "brt",
		CmdCurry:                 "curry",
		CmdCreateList:            "crl",
		CmdListAppend:            "lap",
		CmdCopyMemory:            "mcpy",
		CmdStringAppend:          "sap",
		CmdFixedMul:              "fxmul",
		CmdFixedDiv:              "fxdiv",
		CmdIntNegate:             "neg",
		CmdCreateArray:           "carr",
		CmdLoadInteger:           "ldi",
		CmdLoadBoolean:           "ldb",
		CmdLoadZeroMemoryPointer: "lpzm",
		CmdSetEnum:               "ldenum",
		CmdStringEqual:           "cpseq",
		CmdStringNotEqual:        "cpsneq",
		CmdLoadRune:              "ldch",
	}

	mnemonic, found := names[cmd]
	if !found {
		panic(fmt.Errorf("no lookup for %v", cmd))
	}

	return mnemonic
}
