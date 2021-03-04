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

func decorateCaseCustomType(d DecorateStream, caseExpression *ast.CaseCustomType, context *VariableContext) (*decorated.CaseCustomType, decshared.DecoratedError) {
	decoratedTest, decoratedTestErr := DecorateExpression(d, caseExpression.Test(), context)
	if decoratedTestErr != nil {
		return nil, decoratedTestErr
	}

	pureTestType := dectype.UnaliasWithResolveInvoker(decoratedTest.Type())
	testType := pureTestType

	customType, _ := testType.(*dectype.CustomTypeAtom)
	if customType == nil {
		return nil, decorated.NewMustBeCustomType(decoratedTest)
	}

	handledCustomTypeVariants := make([]bool, customType.VariantCount())

	var decoratedConsequences []*decorated.CaseConsequenceCustomType

	var defaultCase decorated.Expression

	var previousConsequenceType dtype.Type

	for _, consequenceField := range caseExpression.Consequences() {
		var foundVariant *dectype.CustomTypeVariant

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
				return nil, decorated.NewCaseWrongParameterCountInCustomTypeVariant(caseExpression,
					consequenceField, foundVariant)
			}
			for parameterIndex, parameter := range consequenceField.Arguments() {
				parameterType := foundVariant.ParameterTypes()[parameterIndex]
				fakeExpression := NewFakeExpression(parameterType)
				fakeNamedExpresison := decorated.NewNamedDecoratedExpression("__", nil,
					fakeExpression)
				consequenceVariableContext.Add(parameter, fakeNamedExpresison)
			}
		}

		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, consequenceField.Expression(),
			consequenceVariableContext)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}
		if decoratedExpression == nil {
			panic("decoratedExpression == nil ")
		}

		if previousConsequenceType != nil {
			incompatibleErr := dectype.CompatibleTypes(previousConsequenceType, decoratedExpression.Type())
			if incompatibleErr != nil {
				return nil, decorated.NewUnMatchingTypes(consequenceField.Expression(), previousConsequenceType,
					decoratedExpression.Type(), incompatibleErr)
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
				return nil, decorated.NewCaseWrongParameterCountInCustomTypeVariant(caseExpression, consequenceField,
					foundVariant)
			}
			var parameters []*decorated.CaseConsequenceParameter
			for index, argumentType := range foundVariant.ParameterTypes() {
				ident := consequenceField.Arguments()[index]
				param := decorated.NewCaseConsequenceParameter(ident, argumentType)
				parameters = append(parameters, param)
			}

			// Intentionally without module reference for easier reading
			named := decorated.NewNamedDefinitionTypeReference(nil, consequenceField.Identifier())
			variantReference := decorated.NewCustomTypeVariantReference(named, foundVariant)
			decoratedConsequence := decorated.NewCaseConsequenceCustomType(foundVariant.Index(), variantReference,
				parameters, decoratedExpression)
			decoratedConsequences = append(decoratedConsequences, decoratedConsequence)
		}
	}

	if defaultCase == nil {
		var unhandledVariants []*dectype.CustomTypeVariant
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
