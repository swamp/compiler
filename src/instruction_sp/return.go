/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

type Return struct{}

func NewReturn() *Return {
	return &Return{}
}

func (c *Return) Write(writer OpcodeWriter) error {
	writer.Command(CmdReturn)

	return nil
}

func (c *Return) String() string {
	return "ret"
}
