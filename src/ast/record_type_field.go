/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"
)

type RecordTypeField struct {
	symbol            *VariableIdentifier
	userType          Type
	index             int
	precedingComments *MultilineComment
}

func NewRecordTypeField(index int, variable *VariableIdentifier, userType Type, precedingComments *MultilineComment) *RecordTypeField {
	return &RecordTypeField{index: index, symbol: variable, userType: userType, precedingComments: precedingComments}
}

func (i *RecordTypeField) Name() string {
	return i.symbol.Name()
}

func (i *RecordTypeField) VariableIdentifier() *VariableIdentifier {
	return i.symbol
}

func (i *RecordTypeField) Type() Type {
	return i.userType
}

func (i *RecordTypeField) FieldIndex() int {
	return i.index
}

func (i *RecordTypeField) Comment() *MultilineComment {
	return i.precedingComments
}

func (i *RecordTypeField) String() string {
	return fmt.Sprintf("[Field: %v %v]", i.symbol, i.userType)
}
