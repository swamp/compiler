/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"github.com/swamp/compiler/src/token"
)

type VariableIdentifierScoped struct {
	symbol          *VariableIdentifier `debug:"true"`
	moduleReference *ModuleReference
	inclusive       token.SourceFileReference
}

func NewQualifiedVariableIdentifierScoped(moduleReference *ModuleReference, variable *VariableIdentifier) *VariableIdentifierScoped {
	inclusive := token.MakeInclusiveSourceFileReference(moduleReference.FetchPositionLength(), variable.FetchPositionLength())
	return &VariableIdentifierScoped{symbol: variable, moduleReference: moduleReference, inclusive: inclusive}
}

func (i *VariableIdentifierScoped) AstVariableReference() *VariableIdentifier {
	return i.symbol
}

func (i *VariableIdentifierScoped) ModuleReference() *ModuleReference {
	return i.moduleReference
}

func (i *VariableIdentifierScoped) Symbol() token.VariableSymbolToken {
	return i.symbol.Symbol()
}

func (i *VariableIdentifierScoped) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *VariableIdentifierScoped) Name() string {
	return i.moduleReference.ModuleName() + "." + i.symbol.Name()
}

func (i *VariableIdentifierScoped) String() string {
	return i.moduleReference.ModuleName() + "." + i.symbol.String()
}

func (i *VariableIdentifierScoped) DebugString() string {
	return "[scopedidentifier]"
}
