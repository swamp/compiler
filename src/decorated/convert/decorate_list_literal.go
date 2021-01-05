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

func decorateListLiteral(d DecorateStream, list *ast.ListLiteral, context *VariableContext) (decorated.DecoratedExpression, decshared.DecoratedError) {
	var listExpressions []decorated.DecoratedExpression
	var detectedType dtype.Type

	if len(list.Expressions()) > 0 {
		for _, expression := range list.Expressions() {
			decoratedExpression, decoratedExpressionErr := DecorateExpression(d, expression, context)
			if decoratedExpressionErr != nil {
				return nil, decoratedExpressionErr
			}
			if detectedType == nil {
				detectedType = decoratedExpression.Type()
			} else {
				compatibleErr := dectype.CompatibleTypes(detectedType, decoratedExpression.Type())
				if compatibleErr != nil {
					return nil, decorated.NewEveryItemInThelistMustHaveTheSameType(list, expression, detectedType, decoratedExpression.Type(), compatibleErr)
				}
			}
			listExpressions = append(listExpressions, decoratedExpression)
		}
	} else {
		// Empty list
		detectedType = dectype.NewAnyType()
	}

	listType := d.TypeRepo().FindTypeFromAlias("List")
	if listType == nil {
		panic("list literal must have a List type defined to use [] list literals")
	}
	unaliasListType := dectype.Unalias(listType)
	collectionType, wasCollectionType := unaliasListType.(*dectype.PrimitiveAtom)
	if !wasCollectionType {
		panic("must have a List type defined to use [] list literals")
	}
	wrapped := dectype.NewPrimitiveType(collectionType.PrimitiveName(), []dtype.Type{detectedType})

	return decorated.NewListLiteral(wrapped, listExpressions), nil
}
