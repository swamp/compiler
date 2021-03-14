/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func decorateLet(d DecorateStream, let *ast.Let, context *VariableContext) (*decorated.Let, decshared.DecoratedError) {
	var decoratedAssignments []*decorated.LetAssignment
	letVariableContext := context.MakeVariableContext()

	for _, assignment := range let.Assignments() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, assignment.Expression(), letVariableContext)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}

		identifierCount := len(assignment.Identifiers())
		isMultiple := identifierCount > 1

		var letVariables []*decorated.LetVariable
		if isMultiple {
			unaliased := dectype.UnaliasWithResolveInvoker(decoratedExpression.Type())
			atom, resolveErr := unaliased.Resolve()
			if resolveErr != nil {
				return nil, decorated.NewInternalError(fmt.Errorf("wasn't a tuple"))
			}

			tuple, wasTuple := atom.(*dectype.TupleTypeAtom)
			if !wasTuple {
				return nil, decorated.NewInternalError(fmt.Errorf("wasn't a tuple"))
			}
			if tuple.ParameterCount() != identifierCount {
				return nil, decorated.NewInternalError(fmt.Errorf("wrong number of identifiers for the tuple %v vs %v", tuple.ParameterCount(), identifierCount))
			}

			for index, ident := range assignment.Identifiers() {
				variableType := tuple.ParameterTypes()[index]
				letVar := decorated.NewLetVariable(ident, variableType, nil)
				letVariables = append(letVariables, letVar)
			}
		} else {
			letVar := decorated.NewLetVariable(assignment.Identifiers()[0], decoratedExpression.Type(), nil)
			letVariables = []*decorated.LetVariable{letVar}
		}

		decoratedAssignment := decorated.NewLetAssignment(assignment, letVariables, decoratedExpression)
		decoratedAssignments = append(decoratedAssignments, decoratedAssignment)

		for _, letVariable := range letVariables {
			tempNamedExpression := decorated.NewNamedDecoratedExpression("let", nil, letVariable)
			tempNamedExpression.SetReferenced()
			letVariableContext.Add(letVariable.Name(), tempNamedExpression)
		}
	}

	decoratedConsequence, decoratedConsequenceErr := DecorateExpression(d, let.Consequence(), letVariableContext)
	if decoratedConsequenceErr != nil {
		return nil, decoratedConsequenceErr
	}
	return decorated.NewLet(let, decoratedAssignments, decoratedConsequence), nil
}
