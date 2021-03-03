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

func decorateConstructorCall(d DecorateStream, call *ast.ConstructorCall, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	var decoratedExpressions []decorated.Expression

	for _, rawExpression := range call.Arguments() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, rawExpression, context)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}

		decoratedExpressions = append(decoratedExpressions, decoratedExpression)
	}

	variantConstructor := d.TypeRepo().FindTypeFromName(call.TypeIdentifier().Name())
	unaliasedConstructor := dectype.Unalias(variantConstructor)

	switch e := unaliasedConstructor.(type) {
	case *dectype.CustomTypeVariantConstructorType:
		ref := decorated.NewCustomTypeVariantReference(call.TypeIdentifier(), e.Variant())
		return decorated.NewCustomTypeVariantConstructor(ref, decoratedExpressions), nil
	case *dectype.RecordAtom:
		{
			argumentCount := len(decoratedExpressions)
			if argumentCount == 1 {
				first := decoratedExpressions[0]
				recordLiteral, wasRecord := first.(*decorated.RecordLiteral)
				if wasRecord {
					compatibleErr := dectype.CompatibleTypes(recordLiteral.Type(), e)
					if compatibleErr == nil {
						return decorated.NewRecordConstructorRecord(call.TypeIdentifier(), e, recordLiteral), nil
					}
				}
			}

			if len(decoratedExpressions) != len(e.ParseOrderedFields()) {
				return nil, decorated.NewWrongNumberOfFieldsInConstructor(e, call)
			}

			if len(e.ParseOrderedFields()) > 4 {
				return nil, decorated.NewInternalError(fmt.Errorf("maximum of 4 constructor arguments"))
			}

			alphaOrderedAssignments := make([]*decorated.RecordLiteralAssignment, len(decoratedExpressions))
			parsedOrderedAssignments := make([]*decorated.RecordLiteralAssignment, len(decoratedExpressions))
			for index, expr := range decoratedExpressions {
				field := e.ParseOrderedFields()[index]
				targetIndex := field.Index()
				assignment := decorated.NewRecordLiteralAssignment(targetIndex, decorated.NewRecordLiteralField(field.VariableIdentifier()), expr)
				alphaOrderedAssignments[targetIndex] = assignment
				parsedOrderedAssignments[index] = assignment
				unaliasedFieldType := dectype.Unalias(field.Type())
				unaliasedExprType := dectype.Unalias(expr.Type())
				compatibleErr := dectype.CompatibleTypes(unaliasedFieldType, unaliasedExprType)
				if compatibleErr != nil {
					return nil, decorated.NewWrongTypeForRecordConstructorField(field, expr, call, compatibleErr)
				}
			}

			return decorated.NewRecordConstructor(call.TypeIdentifier(), e, alphaOrderedAssignments, decoratedExpressions), nil
		}
	default:
		return nil, decorated.NewExpectedCustomTypeVariantConstructor(call)
	}
}
