/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"
)

type RecordLiteralFieldAssignment struct {
	identifier *VariableIdentifier
	expression Expression
}

func NewRecordLiteralFieldAssignment(identifier *VariableIdentifier, expression Expression) *RecordLiteralFieldAssignment {
	return &RecordLiteralFieldAssignment{identifier: identifier, expression: expression}
}

func (i *RecordLiteralFieldAssignment) Identifier() *VariableIdentifier {
	return i.identifier
}

func (i *RecordLiteralFieldAssignment) Expression() Expression {
	return i.expression
}

func (i *RecordLiteralFieldAssignment) String() string {
	return fmt.Sprintf("[%v = %v]", i.identifier, i.expression)
}

func (i *RecordLiteralFieldAssignment) DebugString() string {
	return fmt.Sprintf("[record-literal-field-assignment]")
}
