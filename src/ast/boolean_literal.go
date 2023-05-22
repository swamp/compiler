/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type BooleanLiteral struct {
	booleanToken token.BooleanToken `debug:"true"`
}

func NewBooleanLiteral(booleanToken token.BooleanToken) *BooleanLiteral {
	return &BooleanLiteral{booleanToken: booleanToken}
}

func (i *BooleanLiteral) Value() bool {
	return i.booleanToken.Value()
}

func (i *BooleanLiteral) Token() token.BooleanToken {
	return i.booleanToken
}

func (i *BooleanLiteral) String() string {
	return fmt.Sprintf("€%v", i.booleanToken.Value())
}

func (i *BooleanLiteral) DebugString() string {
	return i.String()
}

func (i *BooleanLiteral) FetchPositionLength() token.SourceFileReference {
	return i.booleanToken.SourceFileReference
}
