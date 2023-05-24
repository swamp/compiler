/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/concretize"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

func decorateContainerLiteral(d DecorateStream, expressions []ast.Expression, context *VariableContext,
	containerName string, reference token.SourceFileReference) (
	*dectype.ResolvedLocalTypeContext, []decorated.Expression, decshared.DecoratedError,
) {
	var listExpressions []decorated.Expression
	var detectedType dtype.Type

	if len(expressions) > 0 {
		for _, expression := range expressions {
			decoratedExpression, decoratedExpressionErr := DecorateExpression(d, expression, context)
			if decoratedExpressionErr != nil {
				return nil, nil, decoratedExpressionErr
			}
			if detectedType == nil {
				detectedType = decoratedExpression.Type()
			} else {
				compatibleErr := dectype.CompatibleTypes(detectedType, decoratedExpression.Type())
				if compatibleErr != nil {
					return nil, nil, decorated.NewEveryItemInThelistMustHaveTheSameType(
						nil, expression, detectedType, decoratedExpression.Type(), compatibleErr,
					)
				}
			}
			listExpressions = append(listExpressions, decoratedExpression)
		}
	} else {
		// Empty list
		detectedType = dectype.NewAnyType()
	}

	listType := d.TypeReferenceMaker().FindBuiltInType(containerName, reference)
	if listType == nil {
		panic("container literal must have a container type defined to use literals")
	}
	unaliasListType := dectype.Unalias(listType)
	localNameContext, wasCollectionType := unaliasListType.(*dectype.LocalTypeNameOnlyContextReference)
	if !wasCollectionType {
		panic(fmt.Errorf("must have a List type defined to use [] list literals %T", listType))
	}

	concretizedLiteralResolvedContext, concreteErr := concretize.ConcretizeLocalTypeContextUsingArguments(
		localNameContext, []dtype.Type{detectedType},
	)
	if concreteErr != nil {
		return nil, nil, concreteErr
	}

	return concretizedLiteralResolvedContext, listExpressions, nil
}

func decorateListLiteral(d DecorateStream, list *ast.ListLiteral, context *VariableContext) (
	decorated.Expression, decshared.DecoratedError,
) {
	wrappedType, listExpressions, err := decorateContainerLiteral(
		d, list.Expressions(), context, "List", list.FetchPositionLength(),
	)
	if err != nil {
		return nil, err
	}

	return decorated.NewListLiteral(list, wrappedType, listExpressions), nil
}
