/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type RecordField struct {
	symbol   *VariableIdentifier
	userType Type
	index    int
	precedingComments token.CommentBlock
}

func NewRecordTypeField(index int, variable *VariableIdentifier, userType Type, precedingComments token.CommentBlock) *RecordField {
	return &RecordField{index: index, symbol: variable, userType: userType, precedingComments: precedingComments}
}

func (i *RecordField) Name() string {
	return i.symbol.Name()
}

func (i *RecordField) VariableIdentifier() *VariableIdentifier {
	return i.symbol
}

func (i *RecordField) Type() Type {
	return i.userType
}

func (i *RecordField) FieldIndex() int {
	return i.index
}

func (i *RecordField) String() string {
	return fmt.Sprintf("[field: %v %v]", i.symbol, i.userType)
}
