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

func decorateIdentifier(d DecorateStream, ident *ast.VariableIdentifier, context *VariableContext) (decorated.DecoratedExpression, decshared.DecoratedError) {
	def := context.ResolveVariable(ident)
	if def == nil {
		return nil, decorated.NewUnknownVariable(ident)
	}
	return decorated.NewGetVariable(ident, def), nil
}
