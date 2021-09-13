/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

type UnaryOperatorType uint8

const (
	UnaryOperatorBitwiseNot UnaryOperatorType = iota
	UnaryOperatorNot
	UnaryOperatorNegate
)

func UnaryOperatorToOpCode(operator UnaryOperatorType) Commands {
	switch operator {
	case UnaryOperatorBitwiseNot:
		return CmdIntBitwiseNot
	case UnaryOperatorNot:
		return CmdBoolLogicalNot
	case UnaryOperatorNegate:
		return CmdIntNegate
	}

	panic("swamp opcodes: illegal unary operator")
}
