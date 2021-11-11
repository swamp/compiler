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

func decorateResourceName(d DecorateStream, resourceName *ast.ResourceNameLiteral) (decorated.Expression, decshared.DecoratedError) {
	resourceNameType := d.TypeReferenceMaker().FindBuiltInType("ResourceName")
	if resourceNameType == nil {
		panic("internal error. ResourceName is an unknown type")
	}
	decoratedInteger := decorated.NewResourceNameLiteral(resourceName, resourceNameType.(*dectype.PrimitiveAtom))
	return decoratedInteger, nil
}
