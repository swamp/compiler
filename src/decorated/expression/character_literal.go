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

type CharacterLiteral struct {
	str                 *ast.CharacterLiteral `debug:"true"`
	globalCharacterType dtype.Type
}

func NewCharacterLiteral(str *ast.CharacterLiteral, globalCharacterType dtype.Type) *CharacterLiteral {
	return &CharacterLiteral{str: str, globalCharacterType: globalCharacterType}
}

func (i *CharacterLiteral) Type() dtype.Type {
	return i.globalCharacterType
}

func (i *CharacterLiteral) Value() rune {
	return i.str.Value()
}

func (i *CharacterLiteral) String() string {
	return fmt.Sprintf("[Char %v]", i.str.Value())
}

func (i *CharacterLiteral) FetchPositionLength() token.SourceFileReference {
	return i.str.Token.SourceFileReference
}
