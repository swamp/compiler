/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_c

import (
	"fmt"
	"io"

	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateLocalFunctionParameterReference(getVar *decorated.FunctionParameterReference, writer io.Writer, indentation int) error {
	fmt.Fprintf(writer, "%s", getVar.Identifier().Name())
	return nil
}

func generateLetVariableReference(getVar *decorated.LetVariableReference, writer io.Writer, indentation int) error {
	fmt.Fprintf(writer, "%s", getVar.LetVariable().Name().Name())
	return nil
}
