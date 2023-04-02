/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func decorateIf(d DecorateStream, ifExpression *ast.IfExpression,
	context *VariableContext) (*decorated.If, decshared.DecoratedError) {
	condition, conditionErr := DecorateExpression(d, ifExpression.Condition(), context)
	if conditionErr != nil {
		return nil, conditionErr
	}
	if condition == nil {
		panic("condition is nil")
	}
	boolType := d.TypeReferenceMaker().FindBuiltInType("Bool", condition.FetchPositionLength())
	if boolType == nil {
		panic("internal error. Bool type doesn't exist")
	}
	boolCompatibleErr := dectype.CompatibleTypes(boolType, condition.Type())
	if boolCompatibleErr != nil {
		return nil, decorated.NewIfTestMustHaveBooleanType(ifExpression, condition)
	}

	consequence, consequenceErr := DecorateExpression(d, ifExpression.Consequence(), context)
	if consequenceErr != nil {
		return nil, consequenceErr
	}
	alternative, alternativeErr := DecorateExpression(d, ifExpression.Alternative(), context)
	if alternativeErr != nil {
		return nil, alternativeErr
	}

	compatibleErr := dectype.CompatibleTypesCheckCustomType(consequence.Type(), alternative.Type())

	if compatibleErr != nil {
		consPureType := dectype.UnaliasWithResolveInvoker(consequence.Type())
		altPureType := dectype.UnaliasWithResolveInvoker(alternative.Type())
		log.Printf("ERROR %v", ifExpression.FetchPositionLength().ToStandardReferenceString())
		log.Printf("TYPES:\n%v\n%v\n", consPureType, altPureType)
		return nil, decorated.NewIfConsequenceAndAlternativeMustHaveSameType(ifExpression, consequence,
			alternative, compatibleErr)
	}

	return decorated.NewIf(ifExpression, condition, consequence, alternative), nil
}
