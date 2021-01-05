/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func decorateLet(d DecorateStream, let *ast.Let, context *VariableContext) (*decorated.Let, decshared.DecoratedError) {
	var decoratedAssignments []*decorated.LetAssignment
	letVariableContext := context.MakeVariableContext()

	for _, assignment := range let.Assignments() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, assignment.Expression(), letVariableContext)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}
		decoratedAssignment := decorated.NewLetAssignment(assignment.Identifier(), decoratedExpression)
		decoratedAssignments = append(decoratedAssignments, decoratedAssignment)
		tempNamedExpression := decorated.NewNamedDecoratedExpression("let", nil, decoratedExpression)
		letVariableContext.Add(assignment.Identifier(), tempNamedExpression)
	}

	decoratedConsequence, decoratedConsequenceErr := DecorateExpression(d, let.Consequence(), letVariableContext)
	if decoratedConsequenceErr != nil {
		return nil, decoratedConsequenceErr
	}
	return decorated.NewLet(decoratedAssignments, decoratedConsequence), nil
}
