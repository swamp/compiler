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

type StringLiteral struct {
	str              *ast.StringLiteral
	globalStringType dtype.Type
}

func NewStringLiteral(str *ast.StringLiteral, globalStringType dtype.Type) *StringLiteral {
	return &StringLiteral{str: str, globalStringType: globalStringType}
}

func (i *StringLiteral) Type() dtype.Type {
	return i.globalStringType
}

func (i *StringLiteral) Value() string {
	return i.str.Value()
}

func (i *StringLiteral) String() string {
	return fmt.Sprintf("[String %v]", i.str.Value())
}

func (i *StringLiteral) FetchPositionLength() token.SourceFileReference {
	return i.str.Token.SourceFileReference
}
