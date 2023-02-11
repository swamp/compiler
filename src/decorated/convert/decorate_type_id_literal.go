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
)

func decorateTypeId(d DecorateStream, typeId *ast.TypeId) (decorated.Expression, decshared.DecoratedError) {
	typeRefType := d.TypeReferenceMaker().FindBuiltInType("TypeRef")
	if typeRefType == nil {
		panic("internal error. TypeRef is an unknown type")
	}

	contextForTypeRef := dectype.FindNameOnlyContextWithUnalias(typeRefType)
	if contextForTypeRef == nil {
		panic(fmt.Errorf("internal error, TypeRef must have name only context"))
	}
	
	decoratedType, err := d.TypeReferenceMaker().CreateSomeTypeReference(typeId.TypeIdentifier())
	if err != nil {
		return nil, decorated.NewInternalError(err)
	}

	constructedType, err2 := concretize.ConcreteArguments(contextForTypeRef, []dtype.Type{decoratedType})
	if err2 != nil {
		return nil, decorated.NewInternalError(err2)
	}

	return decorated.NewTypeIdLiteral(typeId, constructedType, decoratedType), nil
}
