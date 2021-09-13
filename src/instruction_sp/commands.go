/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

type Commands uint8

const (
	CmdCreateStruct Commands = 0x01
	CmdUpdateStruct          = 0x02
	CmdStructGet             = 0x03
	CmdRegCopy               = 0x04
	CmdListConj              = 0x05

	CmdEnumCase    = 0x06
	CmdBranchFalse = 0x07
	CmdJump        = 0x08

	CmdCall         = 0x09
	CmdReturn       = 0x0a
	CmdCallExternal = 0x0b
	CmdTailCall     = 0x0c

	CmdIntAdd = 0x0d
	CmdIntSub = 0x0e
	CmdIntMul = 0x0f
	CmdIntDiv = 0x10

	// Boolean operators
	CmdIntEqual          = 0x11
	CmdIntNotEqual       = 0x12
	CmdIntLess           = 0x13
	CmdIntLessOrEqual    = 0x14
	CmdIntGreater        = 0x15
	CmdIntGreaterOrEqual = 0x16

	// Bitwise operators
	CmdIntBitwiseAnd = 0x17
	CmdIntBitwiseOr  = 0x18
	CmdIntBitwiseXor = 0x19
	CmdIntBitwiseNot = 0x1a

	CmdBoolLogicalNot      = 0x1b
	CmdBranchTrue          = 0x1c
	CmdCasePatternMatching = 0x1d
	CmdValueEqual          = 0x1e
	CmdValueNotEqual       = 0x1f

	CmdCurry        = 0x20
	CmdCreateList   = 0x21
	CmdListAppend   = 0x22
	CmdCreateEnum   = 0x23
	CmdStringAppend = 0x24

	CmdFixedMul          = 0x25
	CmdFixedDiv          = 0x26
	CmdIntNegate         = 0x27
	CmdStructSplit       = 0x28
	CmdReturnWithMemMove = 0x29
)

func OpcodeToName(cmd Commands) string {
	names := map[Commands]string{
		CmdCreateStruct:        "crs",
		CmdUpdateStruct:        "upd",
		CmdStructGet:           "get",
		CmdRegCopy:             "lr",
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
	}

	return names[cmd]
}
