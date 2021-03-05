/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type VariableIdentifierScoped struct {
	symbol          *VariableIdentifier
	moduleReference *ModuleReference
}

func NewQualifiedVariableIdentifierScoped(moduleReference *ModuleReference, variable *VariableIdentifier) *VariableIdentifierScoped {
	return &VariableIdentifierScoped{symbol: variable, moduleReference: moduleReference}
}

func (i *VariableIdentifierScoped) AstVariableReference() *VariableIdentifier {
	return i.symbol
}

func (i *VariableIdentifierScoped) ModuleReference() *ModuleReference {
	return i.moduleReference
}

func (i *VariableIdentifierScoped) FetchPositionLength() token.SourceFileReference {
	return i.symbol.FetchPositionLength()
}

func (i *VariableIdentifierScoped) Name() string {
	return i.moduleReference.ModuleName() + "." + i.symbol.Name()
}

func (i *VariableIdentifierScoped) String() string {
	return i.moduleReference.ModuleName() + "." + i.symbol.String()
}

func (i *VariableIdentifierScoped) DebugString() string {
	return fmt.Sprintf("[scopedidentifier]")
}
