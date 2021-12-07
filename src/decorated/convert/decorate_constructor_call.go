/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"log"

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

	variantConstructor, err := d.TypeReferenceMaker().CreateSomeTypeReference(call.TypeReference().SomeTypeIdentifier())
	if err != nil {
		return nil, err
	}

	unaliasedConstructor := dectype.Unalias(variantConstructor)

	switch t := variantConstructor.(type) {
	case *dectype.AliasReference:
		{
			_, wasPrimitive := unaliasedConstructor.(*dectype.PrimitiveAtom)
			if wasPrimitive {
				return decorated.NewAliasReference(t), nil
			}

			e, wasRecordType := unaliasedConstructor.(*dectype.RecordAtom)
			if !wasRecordType {
				return nil, decorated.NewInternalError(fmt.Errorf("variantconstructor was not a record atom %T", unaliasedConstructor))
			}
			argumentCount := len(decoratedExpressions)
			if argumentCount == 1 {
				first := decoratedExpressions[0]
				recordLiteral, wasRecord := first.(*decorated.RecordLiteral)
				if wasRecord {
					compatibleErr := dectype.CompatibleTypes(recordLiteral.Type(), e)
					if compatibleErr == nil {
						return decorated.NewRecordConstructorFromRecord(call, t, e, recordLiteral), nil
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
				literalField := decorated.NewRecordLiteralField(field.VariableIdentifier())
				assignment := decorated.NewRecordLiteralAssignment(targetIndex, literalField, expr)
				alphaOrderedAssignments[targetIndex] = assignment
				parsedOrderedAssignments[index] = assignment
				unaliasedFieldType := dectype.Unalias(field.Type())
				unaliasedExprType := dectype.Unalias(expr.Type())
				compatibleErr := dectype.CompatibleTypes(unaliasedFieldType, unaliasedExprType)
				if compatibleErr != nil {
					return nil, decorated.NewWrongTypeForRecordConstructorField(field, expr, call, compatibleErr)
				}
			}

			return decorated.NewRecordConstructorFromParameters(call, t, e, alphaOrderedAssignments, decoratedExpressions), nil
		}
	}

	switch e := unaliasedConstructor.(type) {
	case *dectype.CustomTypeVariantReference:
		return decorated.NewCustomTypeVariantConstructor(e, decoratedExpressions), nil
	case *dectype.RecordAtom:

	default:
		log.Printf("expected a constructor here %T", unaliasedConstructor)
		return nil, decorated.NewExpectedCustomTypeVariantConstructor(call)
	}

	return nil, decorated.NewInternalError(fmt.Errorf("was not a record atom %T", unaliasedConstructor))
}
