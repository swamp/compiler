package generate_c

import (
	"fmt"
	"io"

	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func logicalOperatorToCCode(operatorType decorated.LogicalOperatorType) string {
	if operatorType == decorated.LogicalAnd {
		return "&&"
	} else if operatorType == decorated.LogicalOr {
		return "||"
	}

	panic("unknown type")
}

func generateLogical(operator *decorated.LogicalOperator, writer io.Writer, indentation int) error {
	fmt.Fprintf(writer, "(")
	leftErr := generateExpression(operator.Left(), writer, indentation)
	if leftErr != nil {
		return leftErr
	}

	comparisonString := logicalOperatorToCCode(operator.OperatorType())
	fmt.Fprintf(writer, " "+comparisonString+" ")

	rightErr := generateExpression(operator.Right(), writer, indentation)
	if rightErr != nil {
		return rightErr
	}

	fmt.Fprintf(writer, ")")

	return nil
}
