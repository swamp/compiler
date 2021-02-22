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

type IntegerLiteral struct {
	integer           *ast.IntegerLiteral
	globalIntegerType dtype.Type
}

func NewIntegerLiteral(integer *ast.IntegerLiteral, globalIntegerType dtype.Type) *IntegerLiteral {
	return &IntegerLiteral{integer: integer, globalIntegerType: globalIntegerType}
}

func (i *IntegerLiteral) Type() dtype.Type {
	return i.globalIntegerType
}

func (i *IntegerLiteral) Value() int32 {
	return i.integer.Value()
}

func (i *IntegerLiteral) String() string {
	return fmt.Sprintf("[integer %v]", i.integer.Value())
}

func (i *IntegerLiteral) FetchPositionLength() token.SourceFileReference {
	return i.integer.Token.SourceFileReference
}
