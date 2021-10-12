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
