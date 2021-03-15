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

type RecordTypeFieldReference struct {
	ident                *ast.VariableIdentifier
	recordType           *dectype.RecordAtom
	recordTypeField      *dectype.RecordField
	unresolvedRecordType dtype.Type
}

func (g *RecordTypeFieldReference) Type() dtype.Type {
	return g.recordTypeField.Type()
}

func (g *RecordTypeFieldReference) String() string {
	return fmt.Sprintf("[recordfieldref %v %v]", g.ident, g.recordTypeField)
}

func (g *RecordTypeFieldReference) HumanReadable() string {
	return fmt.Sprintf("Record Field in %v", g.unresolvedRecordType.HumanReadable())
}

func (g *RecordTypeFieldReference) RecordTypeField() *dectype.RecordField {
	return g.recordTypeField
}

func (g *RecordTypeFieldReference) UnresolvedRecordType() dtype.Type {
	return g.unresolvedRecordType
}

func (g *RecordTypeFieldReference) AstIdentifier() *ast.VariableIdentifier {
	return g.ident
}

func NewRecordFieldReference(ident *ast.VariableIdentifier, unresolvedRecordType dtype.Type, recordType *dectype.RecordAtom, recordTypeField *dectype.RecordField) *RecordTypeFieldReference {
	ref := &RecordTypeFieldReference{ident: ident, recordType: recordType, recordTypeField: recordTypeField, unresolvedRecordType: unresolvedRecordType}

	return ref
}

func (g *RecordTypeFieldReference) FetchPositionLength() token.SourceFileReference {
	return g.ident.FetchPositionLength()
}
