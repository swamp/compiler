package generate_c

import (
	"fmt"
	"io"

	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func typeNameString(t dtype.Type) string {
	return t.HumanReadable()
}

func generateLet(let *decorated.Let, writer io.Writer, indentation int) error {
	for _, assignment := range let.Assignments() {
		if len(assignment.LetVariables()) == 1 {
			firstVar := assignment.LetVariables()[0]
			fmt.Fprintf(writer, "%s%v %v = ", indentationString(indentation+1), typeNameString(firstVar.Type()), firstVar.Name().Name())
		} else {
			//tupleType := assignment.Expression().Type().(*dectype.TupleTypeAtom)
			//for index, tupleField := range tupleType.Fields() {
			// variable := assignment.LetVariables()[index]
			//}
		}

		sourceErr := generateExpression(assignment.Expression(), writer, indentation)
		if sourceErr != nil {
			return sourceErr
		}
		fmt.Fprintf(writer, ";\n")

	}

	if !expressionHasOwnReturn(let.Consequence()) {
		fmt.Fprintf(writer, "return ")
	}

	codeErr := generateExpression(let.Consequence(), writer, indentation+1)
	if codeErr != nil {
		return codeErr
	}
	if !expressionHasOwnReturn(let.Consequence()) {
		fmt.Fprintf(writer, ";")
	}

	fmt.Fprintf(writer, "\n%s\n%s", indentationString(indentation), indentationString(indentation))

	return nil
}
