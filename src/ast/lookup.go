/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type Lookups struct {
	contextIdentifier *VariableIdentifier
	fieldNames        []*VariableIdentifier
	inclusive         token.SourceFileReference
}

func NewLookups(contextIdentifier *VariableIdentifier, identifiers []*VariableIdentifier) *Lookups {
	inclusive := token.MakeInclusiveSourceFileReference(contextIdentifier.FetchPositionLength(), identifiers[len(identifiers)-1].FetchPositionLength())
	return &Lookups{contextIdentifier: contextIdentifier, fieldNames: identifiers, inclusive: inclusive}
}

func (i *Lookups) ContextIdentifier() *VariableIdentifier {
	return i.contextIdentifier
}

func (i *Lookups) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *Lookups) FieldNames() []*VariableIdentifier {
	return i.fieldNames
}

func (i *Lookups) String() string {
	return fmt.Sprintf("[lookups %v %v]", i.contextIdentifier, i.fieldNames)
}

func (i *Lookups) DebugString() string {
	return fmt.Sprintf("[Lookups]")
}
