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

func decorateFixed(d DecorateStream, fixed *ast.FixedLiteral) (decorated.Expression, decshared.DecoratedError) {
	fixedType := d.TypeReferenceMaker().FindBuiltInType("Fixed")
	if fixedType == nil {
		panic("internal error. Int is an unknown type")
	}
	decoratedInteger := decorated.NewFixedLiteral(fixed, fixedType.(*dectype.PrimitiveAtom))
	return decoratedInteger, nil
}
