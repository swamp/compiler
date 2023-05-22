/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type ResourceNameLiteral struct {
	Token token.ResourceName
	value string `debug:"true"`
}

func (i *ResourceNameLiteral) String() string {
	return fmt.Sprintf("@%v", i.value)
}

func (i *ResourceNameLiteral) Value() string {
	return i.value
}

func NewResourceNameLiteral(t token.ResourceName, v string) *ResourceNameLiteral {
	return &ResourceNameLiteral{value: v, Token: t}
}

func (i *ResourceNameLiteral) FetchPositionLength() token.SourceFileReference {
	return i.Token.SourceFileReference
}

func (i *ResourceNameLiteral) StringToken() token.ResourceName {
	return i.Token
}

func (i *ResourceNameLiteral) DebugString() string {
	return fmt.Sprintf("[ResourceName %v]", i.value)
}
