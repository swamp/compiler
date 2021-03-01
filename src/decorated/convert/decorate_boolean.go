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
)

func decorateBoolean(d DecorateStream, boolean *ast.BooleanLiteral) (decorated.Expression, decshared.DecoratedError) {
	boolType := d.TypeRepo().FindTypeFromName("Bool")
	if boolType == nil {
		return nil, decorated.NewTypeNotFound("Bool")
	}
	decoratedBoolean := decorated.NewBooleanLiteral(boolean, boolType.(dtype.Type))
	return decoratedBoolean, nil
}
