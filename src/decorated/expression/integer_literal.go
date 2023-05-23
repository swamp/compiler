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

type IntegerLiteral struct {
	integer           *ast.IntegerLiteral             `debug:"true"`
	globalIntegerType *dectype.PrimitiveTypeReference `debug:"true"`
}

func NewIntegerLiteral(integer *ast.IntegerLiteral, globalIntegerType *dectype.PrimitiveTypeReference) *IntegerLiteral {
	return &IntegerLiteral{integer: integer, globalIntegerType: globalIntegerType}
}

func (i *IntegerLiteral) Type() dtype.Type {
	return i.globalIntegerType
}

func (i *IntegerLiteral) Value() int32 {
	return i.integer.Value()
}

func (i *IntegerLiteral) String() string {
	return fmt.Sprintf("[Integer %v]", i.integer.Value())
}

func (i *IntegerLiteral) FetchPositionLength() token.SourceFileReference {
	return i.integer.Token.SourceFileReference
}
