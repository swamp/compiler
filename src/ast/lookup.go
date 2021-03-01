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
}

func NewLookups(contextIdentifier *VariableIdentifier, identifiers []*VariableIdentifier) *Lookups {
	return &Lookups{contextIdentifier: contextIdentifier, fieldNames: identifiers}
}

func (i *Lookups) ContextIdentifier() *VariableIdentifier {
	return i.contextIdentifier
}

func (i *Lookups) FetchPositionLength() token.SourceFileReference {
	return i.contextIdentifier.FetchPositionLength()
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
