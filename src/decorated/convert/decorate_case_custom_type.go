/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func decorateCaseCustomType(d DecorateStream, caseExpression *ast.CaseForCustomType, context *VariableContext) (
	*decorated.CaseCustomType, decshared.DecoratedError,
) {
	//log.Printf("before %v", caseExpression)
	decoratedTest, decoratedTestErr := DecorateExpression(d, caseExpression.Test(), context)
	if decoratedTestErr != nil {
		return nil, decoratedTestErr
	}
	//log.Printf("after %v", decoratedTest)

	pureTestType := dectype.ResolveToAtom(decoratedTest.Type())
	testType := pureTestType

	customType, _ := testType.(*dectype.CustomTypeAtom)
	if customType == nil {
		return nil, decorated.NewMustBeCustomType(decoratedTest)
	}

	handledCustomTypeVariants := make([]bool, customType.VariantCount())

	var decoratedConsequences []*decorated.CaseConsequenceForCustomType

	var defaultCase decorated.Expression

	var previousConsequenceType dtype.Type

	for _, consequenceField := range caseExpression.Consequences() {
		var foundVariant *dectype.CustomTypeVariantAtom
		var parameters []*decorated.CaseConsequenceParameterForCustomType

		consequenceVariableContext := context.MakeVariableContext()

		if !consequenceField.Identifier().IsDefaultSymbol() {
			foundVariant = customType.FindVariant(consequenceField.Identifier().Name())
			if foundVariant == nil {
				return nil, decorated.NewCaseCouldNotFindCustomVariantType(caseExpression, consequenceField)
			}

			foundVariantIndex := foundVariant.Index()

			if handledCustomTypeVariants[foundVariantIndex] {
				return nil, decorated.NewAlreadyHandledCustomTypeVariant(caseExpression, consequenceField, foundVariant)
			}

			handledCustomTypeVariants[foundVariantIndex] = true

			numberOfVariantArguments := len(foundVariant.ParameterTypes())
			if numberOfVariantArguments != len(consequenceField.Arguments()) {
				return nil, decorated.NewCaseWrongParameterCountInCustomTypeVariant(
					caseExpression,
					consequenceField, foundVariant,
				)
			}

			log.Printf("type for found variant: %T %v", customType, customType)
			for index, argumentType := range foundVariant.ParameterTypes() {
				ident := consequenceField.Arguments()[index]
				param := decorated.NewCaseConsequenceParameterForCustomType(ident, argumentType)
				log.Printf("param: %T %v", argumentType, argumentType)
				fakeNamedExpression := decorated.NewNamedDecoratedExpression(
					ident.Name(), nil,
					param,
				)
				consequenceVariableContext.Add(ident, fakeNamedExpression)
				parameters = append(parameters, param)
			}
			/*
				for parameterIndex, parameter := range consequenceField.Arguments() {
					parameterType := foundVariant.ParameterTypes()[parameterIndex]
					parameterExpand := decorated.NewCaseConsequenceParameterForCustomType(parameter, parameterType, foundVariant, parameterIndex)


				}

			*/
		}

		decoratedExpression, decoratedExpressionErr := DecorateExpression(
			d, consequenceField.Expression(),
			consequenceVariableContext,
		)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}
		if decoratedExpression == nil {
			panic("decoratedExpression == nil ")
		}

		if previousConsequenceType != nil {
			incompatibleErr := dectype.CompatibleTypesCheckCustomType(
				previousConsequenceType, decoratedExpression.Type(),
			)
			if incompatibleErr != nil {
				return nil, decorated.NewUnMatchingTypes(
					consequenceField.Expression(), previousConsequenceType,
					decoratedExpression.Type(), incompatibleErr,
				)
			}
		}
		previousConsequenceType = decoratedExpression.Type()

		if consequenceField.Identifier().IsDefaultSymbol() {
			defaultCase = decoratedExpression
			break
		} else {
			expectedArgumentCount := len(foundVariant.ParameterTypes())
			actualArgumentCount := len(consequenceField.Arguments())
			if expectedArgumentCount != actualArgumentCount {
				return nil, decorated.NewCaseWrongParameterCountInCustomTypeVariant(
					caseExpression, consequenceField,
					foundVariant,
				)
			}

			// Intentionally without module reference for easier reading
			fieldTypeRef := ast.NewTypeReference(consequenceField.Identifier(), nil)
			named := dectype.NewNamedDefinitionTypeReference(nil, fieldTypeRef)
			variantReference := dectype.NewCustomTypeVariantReference(named, foundVariant)
			decoratedConsequence := decorated.NewCaseConsequenceForCustomType(
				foundVariant.Index(), variantReference,
				parameters, decoratedExpression, consequenceField,
			)
			decoratedConsequences = append(decoratedConsequences, decoratedConsequence)
		}
	}

	if defaultCase == nil {
		var unhandledVariants []*dectype.CustomTypeVariantAtom
		for index, isHandled := range handledCustomTypeVariants {
			if !isHandled {
				unhandledVariants = append(unhandledVariants, customType.Variants()[index])
			}
		}
		if len(unhandledVariants) != 0 {
			return nil, decorated.NewUnhandledCustomTypeVariants(caseExpression, unhandledVariants)
		}
	}

	c, err := decorated.NewCaseCustomType(caseExpression, decoratedTest, decoratedConsequences, defaultCase)
	return c, err
}
