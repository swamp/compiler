/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

type Commands uint8

const (
	CmdCreateStruct Commands = 0x01
	CmdUpdateStruct Commands = 0x02
	CmdListConj     Commands = 0x05

	CmdEnumCase    Commands = 0x06
	CmdBranchFalse Commands = 0x07
	CmdJump        Commands = 0x08

	CmdCall         Commands = 0x09
	CmdReturn       Commands = 0x0a
	CmdCallExternal Commands = 0x0b
	CmdTailCall     Commands = 0x0c

	CmdIntAdd Commands = 0x0d
	CmdIntSub Commands = 0x0e
	CmdIntMul Commands = 0x0f
	CmdIntDiv Commands = 0x10

	// Boolean operators
	CmdIntEqual          Commands = 0x11
	CmdIntNotEqual       Commands = 0x12
	CmdIntLess           Commands = 0x13
	CmdIntLessOrEqual    Commands = 0x14
	CmdIntGreater        Commands = 0x15
	CmdIntGreaterOrEqual Commands = 0x16

	// Bitwise operators
	CmdIntBitwiseAnd Commands = 0x17
	CmdIntBitwiseOr  Commands = 0x18
	CmdIntBitwiseXor Commands = 0x19
	CmdIntBitwiseNot Commands = 0x1a

	CmdBoolLogicalNot      Commands = 0x1b
	CmdBranchTrue          Commands = 0x1c
	CmdCasePatternMatching Commands = 0x1d
	CmdValueEqual          Commands = 0x1e
	CmdValueNotEqual       Commands = 0x1f

	CmdCurry        Commands = 0x20
	CmdCreateList   Commands = 0x21
	CmdListAppend   Commands = 0x22
	CmdCreateEnum   Commands = 0x23
	CmdStringAppend Commands = 0x24

	CmdFixedMul    Commands = 0x25
	CmdFixedDiv    Commands = 0x26
	CmdIntNegate   Commands = 0x27
	CmdCreateArray Commands = 0x29
)

func OpcodeToName(cmd Commands) string {
	names := map[Commands]string{
		CmdCreateStruct:        "crs",
		CmdUpdateStruct:        "upd",
		CmdListConj:            "conj",
		CmdEnumCase:            "case",
		CmdBranchFalse:         "bne",
		CmdJump:                "jmp",
		CmdCall:                "call",
		CmdReturn:              "ret",
		CmdCallExternal:        "ecall",
		CmdTailCall:            "tcl",
		CmdIntAdd:              "add",
		CmdIntSub:              "sub",
		CmdIntMul:              "mul",
		CmdIntDiv:              "div",
		CmdIntEqual:            "cpeq",
		CmdIntNotEqual:         "cpne",
		CmdIntLess:             "cpl",
		CmdIntLessOrEqual:      "cple",
		CmdIntGreater:          "cpg",
		CmdIntGreaterOrEqual:   "cpge",
		CmdIntBitwiseAnd:       "band",
		CmdIntBitwiseOr:        "bor",
		CmdIntBitwiseXor:       "bxor",
		CmdIntBitwiseNot:       "bnot",
		CmdBoolLogicalNot:      "not",
		CmdBranchTrue:          "brt",
		CmdCasePatternMatching: "csep",
		CmdValueEqual:          "cpve",
		CmdValueNotEqual:       "cpvne",
		CmdCurry:               "curry",
		CmdCreateList:          "crl",
		CmdListAppend:          "lap",
		CmdCreateEnum:          "cre",
		CmdStringAppend:        "sap",
		CmdFixedMul:            "fxmul",
		CmdFixedDiv:            "fxdiv",
		CmdIntNegate:           "neg",
		CmdCreateArray:         "carr",
	}

	return names[cmd]
}
