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
	path                []*TypeIdentifier
	optionalAlias       *TypeIdentifier
	typesToExpose       []*TypeIdentifier
	definitionsToExpose []*VariableIdentifier
	exposeAll           bool
	keyword             token.VariableSymbolToken
	precedingComments   token.CommentBlock
}

func NewImport(keyword token.VariableSymbolToken, path []*TypeIdentifier,
	optionalAlias *TypeIdentifier, typesToExpose []*TypeIdentifier,
	definitionsToExpose []*VariableIdentifier,
	exposeAll bool, precedingComments token.CommentBlock) *Import {
	return &Import{
		keyword: keyword, path: path, optionalAlias: optionalAlias,
		exposeAll:           exposeAll,
		typesToExpose:       typesToExpose,
		definitionsToExpose: definitionsToExpose,
		precedingComments:   precedingComments,
	}
}

func (i *Import) Path() []*TypeIdentifier {
	return i.path
}

func (i *Import) ExposeAll() bool {
	return i.exposeAll
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
	s := fmt.Sprintf("[import %v", i.path)
	if i.optionalAlias != nil {
		s += fmt.Sprintf(" as %v", i.optionalAlias)
	}
	if len(i.typesToExpose) > 0 || len(i.definitionsToExpose) > 0 {
		s += fmt.Sprintf(" exposing (%v %v)", i.typesToExpose, i.definitionsToExpose)
	} else if i.exposeAll {
		s += fmt.Sprintf(" exposing (..)")
	}
	s += "]"

	return s
}

func (i *Import) DebugString() string {
	return fmt.Sprintf("[Import]")
}
