package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func decorateArrayLiteral(d DecorateStream, list *ast.ArrayLiteral, context *VariableContext) (decorated.DecoratedExpression, decshared.DecoratedError) {
	wrappedType, listExpressions, err := decorateContainerLiteral(d, list.Expressions(), context, "Array")
	if err != nil {
		return nil, err
	}

	return decorated.NewArrayLiteral(wrappedType, listExpressions), nil
}
