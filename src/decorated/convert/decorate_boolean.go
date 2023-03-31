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

func decorateBoolean(d DecorateStream, boolean *ast.BooleanLiteral) (decorated.Expression, decshared.DecoratedError) {
	boolType := d.TypeReferenceMaker().FindBuiltInType("Bool")
	if boolType == nil {
		return nil, decorated.NewTypeNotFound("Bool")
	}
	namedTypeRef := dectype.MakeFakeNamedDefinitionTypeReference(boolType.FetchPositionLength(), "Bool")

	boolType = dectype.NewPrimitiveTypeReference(namedTypeRef, boolType.(*dectype.PrimitiveAtom))
	decoratedBoolean := decorated.NewBooleanLiteral(boolean, boolType)
	return decoratedBoolean, nil
}
