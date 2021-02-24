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
	"github.com/swamp/compiler/src/token"
)

func decorateContainerLiteral(d DecorateStream, expressions []ast.Expression, context *VariableContext, containerName string) (*dectype.PrimitiveAtom, []decorated.DecoratedExpression, decshared.DecoratedError) {
	var listExpressions []decorated.DecoratedExpression
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
					return nil, nil, decorated.NewEveryItemInThelistMustHaveTheSameType(nil, expression, detectedType, decoratedExpression.Type(), compatibleErr)
				}
			}
			listExpressions = append(listExpressions, decoratedExpression)
		}
	} else {
		// Empty list
		detectedType = dectype.NewAnyType(ast.NewTypeIdentifier(token.NewTypeSymbolToken("Any", token.SourceFileReference{}, 0)))
	}

	listType := d.TypeRepo().FindTypeFromAlias(containerName)
	if listType == nil {
		panic("container literal must have a container type defined to use literals")
	}
	unaliasListType := dectype.Unalias(listType)
	collectionType, wasCollectionType := unaliasListType.(*dectype.PrimitiveAtom)
	if !wasCollectionType {
		panic("must have a List type defined to use [] list literals")
	}
	wrapped := dectype.NewPrimitiveType(collectionType.PrimitiveName(), []dtype.Type{detectedType})

	return wrapped, listExpressions, nil
}

func decorateListLiteral(d DecorateStream, list *ast.ListLiteral, context *VariableContext) (decorated.DecoratedExpression, decshared.DecoratedError) {
	wrappedType, listExpressions, err := decorateContainerLiteral(d, list.Expressions(), context, "List")
	if err != nil {
		return nil, err
	}

	return decorated.NewListLiteral(wrappedType, listExpressions), nil
}
