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

func decorateInteger(d DecorateStream, integer *ast.IntegerLiteral) (decorated.DecoratedExpression, decshared.DecoratedError) {
	integerType := d.TypeRepo().FindTypeFromName("Int")
	if integerType == nil {
		panic("internal error. Int is an unknown type")
	}
	decoratedInteger := decorated.NewIntegerLiteral(integer, integerType.(*dectype.PrimitiveAtom))
	return decoratedInteger, nil
}
