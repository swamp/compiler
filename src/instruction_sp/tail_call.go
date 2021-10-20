/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

import (
	"fmt"
)

type TailCall struct{}

func NewTailCall() *TailCall {
	return &TailCall{}
}

func (c *TailCall) Write(writer OpcodeWriter) error {
	writer.Command(CmdTailCall)

	return nil
}

func (c *TailCall) String() string {
	return fmt.Sprintf("%v", OpcodeToMnemonic(CmdTailCall))
}
