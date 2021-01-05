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

func decorateAsm(d DecorateStream, asm *ast.Asm) (decorated.DecoratedExpression, decshared.DecoratedError) {
	decoratedAsm := decorated.NewAsmConstant(asm)
	return decoratedAsm, nil
}
