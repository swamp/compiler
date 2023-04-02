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

type StringLiteral struct {
	str              *ast.StringLiteral
	globalStringType *dectype.PrimitiveTypeReference
}

func NewStringLiteral(str *ast.StringLiteral, globalStringType *dectype.PrimitiveTypeReference) *StringLiteral {
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

func (i *StringLiteral) AstString() *ast.StringLiteral {
	return i.str
}

func (i *StringLiteral) FetchPositionLength() token.SourceFileReference {
	return i.str.Token.SourceFileReference
}
