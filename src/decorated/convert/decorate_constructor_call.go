/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/concretize"
	"github.com/swamp/compiler/src/decorated/debug"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func decorateConstructorCall(d DecorateStream, call *ast.ConstructorCall,
	context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	var decoratedArgumentExpressions []decorated.Expression

	for _, rawExpression := range call.Arguments() {
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, rawExpression, context)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}

		decoratedArgumentExpressions = append(decoratedArgumentExpressions, decoratedExpression)
	}

	var argumentTypes []dtype.Type
	for _, argExpression := range decoratedArgumentExpressions {
		argumentTypes = append(argumentTypes, argExpression.Type())
	}

	variantConstructor, err := d.TypeReferenceMaker().CreateSomeTypeReference(call.TypeReference().SomeTypeIdentifier())
	if err != nil {
		return nil, err
	}
	log.Printf("variantConstructor: %T", variantConstructor)
	// 	concretize.ConcretizeLocalTypeContextUsingArguments()

	unaliasedConstructor := dectype.Unalias(variantConstructor)

	switch t := variantConstructor.(type) {
	case *dectype.AliasReference:
		{
			_, wasPrimitive := unaliasedConstructor.(*dectype.PrimitiveAtom)
			if wasPrimitive {
				return decorated.NewAliasReference(variantConstructor.NameReference(), t), nil
			}

			e, wasRecordType := unaliasedConstructor.(*dectype.RecordAtom)
			if !wasRecordType {
				return nil, decorated.NewInternalError(fmt.Errorf("variantconstructor was not a record atom %T",
					unaliasedConstructor))
			}
			argumentCount := len(decoratedArgumentExpressions)
			if argumentCount == 1 {
				first := decoratedArgumentExpressions[0]
				recordLiteral, wasRecord := first.(*decorated.RecordLiteral)
				if wasRecord {
					compatibleErr := dectype.CompatibleTypes(recordLiteral.Type(), e)
					if compatibleErr == nil {
						return decorated.NewRecordConstructorFromRecord(call, t, e, recordLiteral), nil
					}
				}
			}

			if len(decoratedArgumentExpressions) != len(e.ParseOrderedFields()) {
				return nil, decorated.NewWrongNumberOfFieldsInConstructor(e, call)
			}

			if len(e.ParseOrderedFields()) > 4 {
				return nil, decorated.NewInternalError(fmt.Errorf("maximum of 4 constructor arguments"))
			}

			alphaOrderedAssignments := make([]*decorated.RecordLiteralAssignment, len(decoratedArgumentExpressions))
			parsedOrderedAssignments := make([]*decorated.RecordLiteralAssignment, len(decoratedArgumentExpressions))
			for index, expr := range decoratedArgumentExpressions {
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

			return decorated.NewRecordConstructorFromParameters(call, t, e, alphaOrderedAssignments,
				decoratedArgumentExpressions), nil
		}
	case *dectype.CustomTypeVariantReference:
		localTypeContext, _ := t.Next().(*dectype.LocalTypeNameOnlyContextReference)
		var variantRef *dectype.CustomTypeVariantReference
		log.Printf("variant ref: %T", t.Next())
		if localTypeContext != nil {
			concrete, resolveErr := concretize.ConcretizeLocalTypeContextUsingArguments(localTypeContext, argumentTypes)
			if resolveErr != nil {
				return nil, resolveErr
			}
			log.Printf("resolved to %s", debug.TreeString(concrete))
		} else {
			variantRef = t
		}
		return decorated.NewCustomTypeVariantConstructor(variantRef, decoratedArgumentExpressions), nil
	default:
		panic(fmt.Errorf("not sure what it is now %T", variantConstructor))
	}

	switch unaliasedConstructor.(type) {
	case *dectype.CustomTypeVariantAtom:
	case *dectype.RecordAtom:

	default:
		log.Printf("expected a constructor here %T", unaliasedConstructor)
		return nil, decorated.NewExpectedCustomTypeVariantConstructor(call)
	}

	return nil, decorated.NewInternalError(fmt.Errorf("was not a record atom %T", unaliasedConstructor))
}
