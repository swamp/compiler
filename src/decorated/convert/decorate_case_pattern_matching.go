/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type FakeExpression struct {
	storedType dtype.Type
}

func NewFakeExpression(storedType dtype.Type) *FakeExpression {
	return &FakeExpression{storedType: storedType}
}

func (f *FakeExpression) Type() dtype.Type {
	return f.storedType
}

func (f *FakeExpression) String() string {
	return "totally fake"
}

func (f *FakeExpression) FetchPositionLength() token.Range {
	return token.Range{}
}

func decorateCasePatternMatching(d DecorateStream, caseExpression *ast.CasePatternMatching, context *VariableContext) (*decorated.CasePatternMatching, decshared.DecoratedError) {
	decoratedTest, decoratedTestErr := DecorateExpression(d, caseExpression.Test(), context)
	if decoratedTestErr != nil {
		return nil, decoratedTestErr
	}

	pureTestType := dectype.UnaliasWithResolveInvoker(decoratedTest.Type())
	testType := pureTestType

	var decoratedConsequences []*decorated.CaseConsequencePatternMatching

	var defaultCase decorated.DecoratedExpression

	var previousConsequenceType dtype.Type

	for _, consequence := range caseExpression.Consequences() {
		var decoratedLiteralExpression decorated.DecoratedExpression
		if consequence.Literal() != nil {
			consequenceVariableContext := context.MakeVariableContext()
			var decoratedLiteralExpressionErr decshared.DecoratedError
			decoratedLiteralExpression, decoratedLiteralExpressionErr = DecorateExpression(d, consequence.Literal(),
				consequenceVariableContext)
			if decoratedLiteralExpressionErr != nil {
				return nil, decoratedLiteralExpressionErr
			}

			incompatibleErr := dectype.CompatibleTypes(testType, decoratedLiteralExpression.Type())
			if incompatibleErr != nil {
				fmt.Printf("test type and literal must be compatible %v %v\n", testType, decoratedLiteralExpression.Type())
				return nil, decorated.NewUnMatchingTypes(consequence.Expression(), testType,
					decoratedLiteralExpression.Type(), incompatibleErr)
			}
		}

		consequenceExpressionContext := context.MakeVariableContext()
		decoratedExpression, decoratedExpressionErr := DecorateExpression(d, consequence.Expression(),
			consequenceExpressionContext)
		if decoratedExpressionErr != nil {
			return nil, decoratedExpressionErr
		}
		if decoratedExpression == nil {
			panic("decoratedExpression == nil ")
		}

		if previousConsequenceType != nil {
			incompatibleErr := dectype.CompatibleTypes(previousConsequenceType, decoratedExpression.Type())
			if incompatibleErr != nil {
				return nil, decorated.NewUnMatchingTypes(consequence.Expression(), previousConsequenceType,
					decoratedExpression.Type(), incompatibleErr)
			}
		}
		previousConsequenceType = decoratedExpression.Type()

		if consequence.Literal() == nil {
			defaultCase = decoratedExpression
			break
		} else {
			decoratedConsequence := decorated.NewCaseConsequencePatternMatching(consequence.Index(), decoratedLiteralExpression,
				decoratedExpression)
			decoratedConsequences = append(decoratedConsequences, decoratedConsequence)
		}
	}

	if defaultCase == nil {
		return nil, decorated.NewInternalError(fmt.Errorf("must have a default case"))
	}

	c, err := decorated.NewCasePatternMatching(decoratedTest, decoratedConsequences, defaultCase)
	return c, err
}
