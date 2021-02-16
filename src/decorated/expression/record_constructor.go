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

type RecordConstructor struct {
	arguments      []*RecordLiteralAssignment
	recordType     *dectype.RecordAtom
	typeIdentifier *ast.TypeIdentifier
}

func NewRecordConstructor(typeIdentifier *ast.TypeIdentifier, recordType *dectype.RecordAtom, arguments []*RecordLiteralAssignment) *RecordConstructor {
	return &RecordConstructor{typeIdentifier: typeIdentifier, arguments: arguments, recordType: recordType}
}

func (c *RecordConstructor) SortedAssignments() []*RecordLiteralAssignment {
	return c.arguments
}

func (c *RecordConstructor) Type() dtype.Type {
	return c.recordType
}

func (c *RecordConstructor) String() string {
	return fmt.Sprintf("[record-constructor %v %v]", c.typeIdentifier, c.arguments)
}

func (c *RecordConstructor) FetchPositionAndLength() token.PositionLength {
	return c.typeIdentifier.Symbol().FetchPositionLength()
}
