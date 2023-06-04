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

	var allLetVariables []*decorated.LetVariable
	for _, assignment := range let.Assignments() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, assignment.Expression(),
			letVariableContext)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}

		identifierCount := len(assignment.Identifiers())

		var letVariables []*decorated.LetVariable
		if assignment.WasRecordDestructuring() {
			atom := dectype.ResolveToAtom(decoratedExpression.Type())
			record, wasRecord := atom.(*dectype.RecordAtom)
			if !wasRecord {
				return nil, decorated.NewRecordDestructuringWasNotRecordExpression(decoratedExpression, record)
			}
			for _, ident := range assignment.Identifiers() {
				recordField := record.FindField(ident.Symbol().Name())
				if recordField == nil {
					return nil, decorated.NewRecordDestructuringFieldNotFound(decoratedExpression, record, ident)
				}
				letVar := decorated.NewLetVariable(ident, recordField.Type(), assignment.CommentBlock())
				letVariables = append(letVariables, letVar)
			}
		} else {
			isMultiple := identifierCount > 1
			if isMultiple {
				atom := dectype.ResolveToAtom(decoratedExpression.Type())

				tuple, wasTuple := atom.(*dectype.TupleTypeAtom)
				if !wasTuple {
					return nil, decorated.NewInternalError(fmt.Errorf("wasn't a tuple"))
				}
				if tuple.ParameterCount() != identifierCount {
					return nil, decorated.NewTupleDestructuringWrongNumberOfIdentifiers(decoratedExpression, tuple,
						assignment.Identifiers())
				}

				for index, ident := range assignment.Identifiers() {
					variableType := tuple.ParameterTypes()[index]
					letVar := decorated.NewLetVariable(ident, variableType, assignment.CommentBlock())
					letVariables = append(letVariables, letVar)
				}
			} else {
				letVar := decorated.NewLetVariable(assignment.Identifiers()[0], decoratedExpression.Type(),
					assignment.CommentBlock())
				letVariables = []*decorated.LetVariable{letVar}
			}
		}
		decoratedAssignment := decorated.NewLetAssignment(assignment, letVariables, decoratedExpression)
		decoratedAssignments = append(decoratedAssignments, decoratedAssignment)

		allLetVariables = append(allLetVariables, letVariables...)

		for _, letVariable := range letVariables {
			if letVariable.IsIgnore() {
				continue
			}
			tempNamedExpression := decorated.NewNamedDecoratedExpression("let", nil, letVariable)
			tempNamedExpression.SetReferenced()
			letVariableContext.Add(letVariable.Name(), tempNamedExpression)
		}
	}

	decoratedConsequence, decoratedConsequenceErr := DecorateExpression(d, let.Consequence(), letVariableContext)
	if decoratedConsequenceErr != nil {
		return nil, decoratedConsequenceErr
	}

	for _, letVariable := range allLetVariables {
		if letVariable.IsIgnore() {
			continue
		}
		if !letVariable.WasReferenced() {
			unusedErr := decorated.NewUnusedLetVariable(letVariable)
			d.AddDecoratedError(unusedErr)
		}
	}

	return decorated.NewLet(let, decoratedAssignments, decoratedConsequence), nil
}
