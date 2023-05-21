/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

func NewAnyType() *PrimitiveAtom {
	return NewPrimitiveType(ast.NewTypeIdentifier(token.NewTypeSymbolToken("Any", token.NewInternalSourceFileReferenceRow(1), 0)), nil)
}
