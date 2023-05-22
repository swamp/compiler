/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type StringLiteral struct {
	Token token.StringToken `debug:"true"`
}

func (i *StringLiteral) String() string {
	return fmt.Sprintf("'%v'", i.Token.Text())
}

func (i *StringLiteral) Value() string {
	return i.Token.Text()
}

func NewStringLiteral(t token.StringToken) *StringLiteral {
	return &StringLiteral{Token: t}
}

func (i *StringLiteral) FetchPositionLength() token.SourceFileReference {
	return i.Token.SourceFileReference
}

func (i *StringLiteral) StringToken() token.StringToken {
	return i.Token
}

func (i *StringLiteral) DebugString() string {
	return fmt.Sprintf("[String %v]", i.Token.Text())
}
