/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	dectype "github.com/swamp/compiler/src/decorated/types"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type FixedLiteral struct {
	integer         *ast.FixedLiteral `debug:"true"`
	globalFixedType *dectype.PrimitiveTypeReference
}

func NewFixedLiteral(integer *ast.FixedLiteral, globalFixedType *dectype.PrimitiveTypeReference) *FixedLiteral {
	return &FixedLiteral{integer: integer, globalFixedType: globalFixedType}
}

func (i *FixedLiteral) Type() dtype.Type {
	return i.globalFixedType
}

func (i *FixedLiteral) Value() int32 {
	return i.integer.Value()
}

func (i *FixedLiteral) String() string {
	return fmt.Sprintf("[Fixed %v]", i.integer.Value())
}

func (i *FixedLiteral) FetchPositionLength() token.SourceFileReference {
	return i.integer.Token.SourceFileReference
}
