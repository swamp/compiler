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

func decorateTypeId(d DecorateStream, typeId *ast.TypeId) (decorated.Expression, decshared.DecoratedError) {
	typeRefType := d.TypeRepo().FindTypeFromName("TypeRef")
	if typeRefType == nil {
		panic("internal error. TypeRef is an unknown type")
	}

	decoratedType, err := ConvertFromAstToDecorated(typeId.TypeRef(), d.TypeRepo())
	if err != nil {
		return nil, decorated.NewInternalError(err)
	}

	/*
		constructedType, err2 := dectype.NewInvokerType(typeRefType, []dtype.Type{decoratedType})
		if err2 != nil {
			return nil, decorated.NewInternalError(err2)
		}
	*/

	constructedType, err2 := dectype.CallType(typeRefType, []dtype.Type{decoratedType})
	if err2 != nil {
		return nil, decorated.NewInternalError(err2)
	}

	return decorated.NewTypeIdLiteral(typeId, constructedType, decoratedType), nil
}
