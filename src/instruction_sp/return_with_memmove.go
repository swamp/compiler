/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package instruction_sp

type ReturnWithMemMove struct{}

func NewReturnWithMemMove() *ReturnWithMemMove {
	return &ReturnWithMemMove{}
}

func (c *ReturnWithMemMove) Write(writer OpcodeWriter) error {
	writer.Command(CmdReturnWithMemMove)

	return nil
}

func (c *ReturnWithMemMove) String() string {
	return "ret"
}
