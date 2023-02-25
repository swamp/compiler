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
	"log"
)

func decorateContainerLiteral(d DecorateStream, expressions []ast.Expression, context *VariableContext, containerName string) (*dectype.PrimitiveTypeReference, []decorated.Expression, decshared.DecoratedError) {
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
					return nil, nil, decorated.NewEveryItemInThelistMustHaveTheSameType(nil, expression, detectedType, decoratedExpression.Type(), compatibleErr)
				}
			}
			listExpressions = append(listExpressions, decoratedExpression)
		}
	} else {
		// Empty list
		detectedType = dectype.NewAnyType()
	}

	listType := d.TypeReferenceMaker().FindBuiltInType(containerName)
	if listType == nil {
		panic("container literal must have a container type defined to use literals")
	}
	unaliasListType := dectype.Unalias(listType)
	localNameContext, wasCollectionType := unaliasListType.(*dectype.LocalTypeNameContext)
	if !wasCollectionType {
		panic(fmt.Errorf("must have a List type defined to use [] list literals %T", listType))
	}

	concretizedLiteral, concreteErr := concretize.ConcreteArguments(localNameContext, []dtype.Type{detectedType})
	if concreteErr != nil {
		return nil, nil, concreteErr
	}
	log.Printf("concreteListLiteral %v", concretizedLiteral)

	primitiveAtom, _ := concretizedLiteral.(*dectype.PrimitiveAtom)
	typeIdent := ast.NewTypeIdentifier(token.NewTypeSymbolToken(containerName, listType.FetchPositionLength(), 0))
	typeRef := ast.NewTypeReference(typeIdent, nil)
	decTypeRef := dectype.NewPrimitiveTypeReference(dectype.NewNamedDefinitionTypeReference(nil, typeRef), primitiveAtom)

	return decTypeRef, listExpressions, nil
}

func decorateListLiteral(d DecorateStream, list *ast.ListLiteral, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	wrappedType, listExpressions, err := decorateContainerLiteral(d, list.Expressions(), context, "List")
	if err != nil {
		return nil, err
	}

	return decorated.NewListLiteral(list, wrappedType, listExpressions), nil
}
