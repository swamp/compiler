/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func decorateTupleLiteral(d DecorateStream, astTuple *ast.TupleLiteral, context *VariableContext) (*decorated.TupleLiteral, decshared.DecoratedError) {
	var tupleExpressions []decorated.Expression
	var foundTypes []dtype.Type
	for _, expression := range astTuple.Expressions() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, expression, context)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}
		tupleExpressions = append(tupleExpressions, decoratedExpression)
		foundTypes = append(foundTypes, decoratedExpression.Type())
	}

	// astTupleType := ast.NewTupleType(token.ParenToken{}, token.ParenToken{}, )
	tupleType := dectype.NewTupleTypeAtom(nil, foundTypes)

	return decorated.NewTupleLiteral(astTuple, tupleType, tupleExpressions), nil
}
