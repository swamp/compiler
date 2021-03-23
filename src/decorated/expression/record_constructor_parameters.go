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

type RecordConstructorFromParameters struct {
	arguments            []*RecordLiteralAssignment
	parseOrderArguments  []Expression
	recordAliasReference *dectype.AliasReference
	recordType           *dectype.RecordAtom
	astConstructorCall   *ast.ConstructorCall
}

// typeIdentifier ast.TypeReferenceScopedOrNormal
func NewRecordConstructorFromParameters(astConstructorCall *ast.ConstructorCall, recordAliasReference *dectype.AliasReference, recordType *dectype.RecordAtom, arguments []*RecordLiteralAssignment, parseOrderArguments []Expression) *RecordConstructorFromParameters {
	if recordAliasReference == nil {
		panic("can not be nil")
	}
	return &RecordConstructorFromParameters{astConstructorCall: astConstructorCall, recordAliasReference: recordAliasReference, arguments: arguments, parseOrderArguments: parseOrderArguments, recordType: recordType}
}

func (c *RecordConstructorFromParameters) SortedAssignments() []*RecordLiteralAssignment {
	return c.arguments
}

func (c *RecordConstructorFromParameters) NamedTypeReference() *dectype.NamedDefinitionTypeReference {
	return c.recordAliasReference.NameReference()
}

func (c *RecordConstructorFromParameters) ParseOrderArguments() []Expression {
	return c.parseOrderArguments
}

func (c *RecordConstructorFromParameters) Type() dtype.Type {
	return c.recordType
}

func (c *RecordConstructorFromParameters) RecordType() *dectype.RecordAtom {
	return c.recordType
}

func (c *RecordConstructorFromParameters) String() string {
	return fmt.Sprintf("[record-constructor %v %v]", c.recordAliasReference, c.arguments)
}

func (c *RecordConstructorFromParameters) HumanReadable() string {
	return "Record Constructor"
}

func (c *RecordConstructorFromParameters) AstConstructorCall() *ast.ConstructorCall {
	return c.astConstructorCall
}

func (c *RecordConstructorFromParameters) FetchPositionLength() token.SourceFileReference {
	return c.astConstructorCall.FetchPositionLength()
}
