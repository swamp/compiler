/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

type TailCall struct{}

func NewTailCall() *TailCall {
	return &TailCall{}
}

func (c *TailCall) Write(writer OpcodeWriter) error {
	writer.Command(CmdTailCall)

	return nil
}
