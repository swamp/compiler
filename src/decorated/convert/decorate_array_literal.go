/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func decorateArrayLiteral(d DecorateStream, list *ast.ArrayLiteral, context *VariableContext) (decorated.Expression,
	decshared.DecoratedError) {
	wrappedType, listExpressions, err := decorateContainerLiteral(d, list.Expressions(), context, "Array", list.FetchPositionLength())
	if err != nil {
		return nil, err
	}

	return decorated.NewArrayLiteral(list, wrappedType, listExpressions), nil
}
