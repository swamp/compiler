/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type TupleType struct {
	types      []Type
	inclusive  token.SourceFileReference
	startParen token.ParenToken
	endParen   token.ParenToken
}

func (i *TupleType) String() string {
	return fmt.Sprintf("[tuple-type %v]", i.types)
}

func (i *TupleType) DebugString() string {
	return fmt.Sprintf("[tuple-type %v]", i.types)
}

func (i *TupleType) DecoratedName() string {
	return ""
}

func (i *TupleType) Name() string {
	return "tupletype"
}

func (i *TupleType) Types() []Type {
	return i.types
}

func (i *TupleType) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *TupleType) StartParen() token.ParenToken {
	return i.startParen
}

func (i *TupleType) EndParen() token.ParenToken {
	return i.endParen
}

func NewTupleType(startParen token.ParenToken, endParen token.ParenToken, types []Type) *TupleType {
	inclusive := token.MakeInclusiveSourceFileReference(startParen.FetchPositionLength(), endParen.FetchPositionLength())
	return &TupleType{
		types:      types,
		startParen: startParen,
		endParen:   endParen,
		inclusive:  inclusive,
	}
}
