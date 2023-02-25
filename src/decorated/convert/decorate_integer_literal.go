/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func decorateInteger(d DecorateStream, integer *ast.IntegerLiteral) (decorated.Expression, decshared.DecoratedError) {
	integerType := d.TypeReferenceMaker().FindBuiltInType("Int")
	if integerType == nil {
		panic("internal error. Int is an unknown type")
	}

	namedTypeRef := dectype.MakeFakeNamedDefinitionTypeReference(integerType.FetchPositionLength(), "Int")

	primitiveReference := dectype.NewPrimitiveTypeReference(namedTypeRef, integerType.(*dectype.PrimitiveAtom))

	decoratedInteger := decorated.NewIntegerLiteral(integer, primitiveReference)
	return decoratedInteger, nil
}
