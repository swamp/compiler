/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type Import struct {
	path []*TypeIdentifier
	optionalAlias *TypeIdentifier
	keyword token.VariableSymbolToken
	precedingComments token.CommentBlock
}

func NewImport(keyword token.VariableSymbolToken, path []*TypeIdentifier,
	optionalAlias *TypeIdentifier, precedingComments token.CommentBlock) *Import {
	return &Import{keyword:keyword, path: path, optionalAlias: optionalAlias, precedingComments: precedingComments}
}

func (i *Import) Path() []*TypeIdentifier {
	return i.path
}

func (i *Import) Alias() *TypeIdentifier {
	return i.optionalAlias
}

func (i *Import) ModuleName() []*TypeIdentifier {
	return i.path
}

func (i *Import) PositionLength() token.PositionLength {
	return i.keyword.PositionLength
}

func (i *Import) String() string {
	return fmt.Sprintf("[import %v]", i.path)
}

func (i *Import) DebugString() string {
	return fmt.Sprintf("[Import]")
}
