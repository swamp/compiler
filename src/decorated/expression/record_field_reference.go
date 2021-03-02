/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type RecordFieldReference struct {
	ident           *ast.VariableIdentifier
	recordType      *dectype.RecordAtom
	recordTypeField *dectype.RecordField
}

func (g *RecordFieldReference) Type() dtype.Type {
	return g.recordTypeField.Type()
}

func (g *RecordFieldReference) String() string {
	return fmt.Sprintf("[recordfieldref %v %v]", g.ident, g.recordTypeField)
}

func (g *RecordFieldReference) HumanReadable() string {
	return fmt.Sprintf("%v", g.ident)
}

func (g *RecordFieldReference) RecordTypeField() *dectype.RecordField {
	return g.recordTypeField
}

func (g *RecordFieldReference) AstIdentifier() *ast.VariableIdentifier {
	return g.ident
}

func NewRecordFieldReference(ident *ast.VariableIdentifier, recordType *dectype.RecordAtom, recordTypeField *dectype.RecordField) *RecordFieldReference {
	ref := &RecordFieldReference{ident: ident, recordType: recordType, recordTypeField: recordTypeField}

	return ref
}

func (g *RecordFieldReference) FetchPositionLength() token.SourceFileReference {
	return g.ident.FetchPositionLength()
}
