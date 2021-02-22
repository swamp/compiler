/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type FixedLiteral struct {
	integer         *ast.FixedLiteral
	globalFixedType dtype.Type
}

func NewFixedLiteral(integer *ast.FixedLiteral, globalFixedType dtype.Type) *FixedLiteral {
	return &FixedLiteral{integer: integer, globalFixedType: globalFixedType}
}

func (i *FixedLiteral) Type() dtype.Type {
	return i.globalFixedType
}

func (i *FixedLiteral) Value() int32 {
	return i.integer.Value()
}

func (i *FixedLiteral) String() string {
	return fmt.Sprintf("[integer %v]", i.integer.Value())
}

func (i *FixedLiteral) FetchPositionLength() token.Range {
	return i.integer.Token.FetchPositionLength()
}
