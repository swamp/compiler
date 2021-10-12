package generate_c

import (
	"fmt"
	"io"

	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateIf(ifExpr *decorated.If, writer io.Writer, returnString string, indentation int) error {
	fmt.Fprintf(writer, "if (")
	testErr := generateExpression(ifExpr.Condition(), writer, returnString, indentation)
	if testErr != nil {
		return testErr
	}
	fmt.Fprintf(writer, ") {\n%s", indentationString(indentation))
	consErr := generateExpression(ifExpr.Consequence(), writer, returnString, indentation)
	if consErr != nil {
		return consErr
	}
	fmt.Fprintf(writer, "} else {\n%s", indentationString(indentation))

	altErr := generateExpression(ifExpr.Alternative(), writer, returnString, indentation)
	if altErr != nil {
		return altErr
	}

	fmt.Fprintf(writer, "}\n")
	return nil
}
