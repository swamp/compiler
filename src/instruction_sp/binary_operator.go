/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

// Package swampopcodeinst holds all the instructions for the Swamp runtime.
package instruction_sp

// BinaryOperatorType defines the type of binary operator.
type BinaryOperatorType uint8

// The binary operator types.
const (
	BinaryOperatorArithmeticIntPlus BinaryOperatorType = iota
	BinaryOperatorArithmeticIntMinus
	BinaryOperatorArithmeticIntDivide
	BinaryOperatorArithmeticIntMultiply
	BinaryOperatorBooleanIntEqual
	BinaryOperatorBooleanIntNotEqual
	BinaryOperatorBooleanIntLess
	BinaryOperatorBooleanIntLessOrEqual
	BinaryOperatorBooleanIntGreater
	BinaryOperatorBooleanIntGreaterOrEqual
	BinaryOperatorBitwiseIntAnd
	BinaryOperatorBitwiseIntOr
	BinaryOperatorBitwiseIntXor
	BinaryOperatorArithmeticListAppend
	BinaryOperatorArithmeticFixedDivide
	BinaryOperatorArithmeticFixedMultiply
	BinaryOperatorBooleanValueEqual
	BinaryOperatorBooleanValueNotEqual
)

// BinaryOperatorToOpCode converts from the type of binary operator to the actual opcode instruction.
func BinaryOperatorToOpCode(operator BinaryOperatorType) Commands {
	switch operator {
	case BinaryOperatorArithmeticIntPlus:
		return CmdIntAdd
	case BinaryOperatorArithmeticIntMinus:
		return CmdIntSub
	case BinaryOperatorArithmeticIntDivide:
		return CmdIntDiv
	case BinaryOperatorArithmeticIntMultiply:
		return CmdIntMul
	case BinaryOperatorArithmeticFixedDivide:
		return CmdFixedDiv
	case BinaryOperatorArithmeticFixedMultiply:
		return CmdFixedMul
	case BinaryOperatorArithmeticListAppend:
		return CmdListAppend
	case BinaryOperatorBooleanIntEqual:
		return CmdIntEqual
	case BinaryOperatorBooleanIntNotEqual:
		return CmdIntNotEqual
	case BinaryOperatorBooleanIntLess:
		return CmdIntLess
	case BinaryOperatorBooleanIntLessOrEqual:
		return CmdIntLessOrEqual
	case BinaryOperatorBooleanIntGreater:
		return CmdIntGreater
	case BinaryOperatorBooleanIntGreaterOrEqual:
		return CmdIntGreaterOrEqual
	case BinaryOperatorBitwiseIntAnd:
		return CmdIntBitwiseAnd
	case BinaryOperatorBitwiseIntOr:
		return CmdIntBitwiseOr
	case BinaryOperatorBitwiseIntXor:
		return CmdIntBitwiseXor
	case BinaryOperatorBooleanValueEqual:
		return CmdValueEqual
	case BinaryOperatorBooleanValueNotEqual:
		return CmdValueNotEqual
	}

	panic("swamp opcodes: unknown binary operator")
}
