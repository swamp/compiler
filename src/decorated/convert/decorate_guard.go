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

func decorateGuard(d DecorateStream, guardExpression *ast.GuardExpression,
	context *VariableContext) (*decorated.Guard, decshared.DecoratedError) {
	var items []*decorated.GuardItem

	var detectedType dtype.Type
	for index, item := range guardExpression.Items() {
		condition, conditionErr := DecorateExpression(d, item.Condition, context)
		if conditionErr != nil {
			return nil, conditionErr
		}
		if condition == nil {
			panic("condition is nil")
		}
		boolType := d.TypeRepo().FindTypeFromName("Bool")
		if boolType == nil {
			panic("internal error. Bool type doesn't exist")
		}
		boolCompatibleErr := dectype.CompatibleTypes(boolType, condition.Type())
		if boolCompatibleErr != nil {
			return nil, decorated.NewIfTestMustHaveBooleanType(nil, condition)
		}

		consequence, consequenceErr := DecorateExpression(d, item.Consequence, context)
		if consequenceErr != nil {
			return nil, consequenceErr
		}

		item := decorated.NewGuardItem(index, condition, consequence)
		items = append(items, item)
		if index == 0 {
			detectedType = consequence.Type()
		} else {
			allSameErr := dectype.CompatibleTypes(detectedType, consequence.Type())
			if allSameErr != nil {
				return nil, decorated.NewEveryItemInThelistMustHaveTheSameType(nil, nil, nil, nil, allSameErr)
			}
		}
	}

	defaultDecoratedExpression, defaultExpressionErr := DecorateExpression(d, guardExpression.DefaultExpression(), context)
	if defaultExpressionErr != nil {
		return nil, defaultExpressionErr
	}

	compatibleErr := dectype.CompatibleTypes(detectedType, defaultDecoratedExpression.Type())
	if compatibleErr != nil {
		return nil, decorated.NewIfConsequenceAndAlternativeMustHaveSameType(nil, nil,
			nil, compatibleErr)
	}

	return decorated.NewGuard(items, defaultDecoratedExpression)
}
