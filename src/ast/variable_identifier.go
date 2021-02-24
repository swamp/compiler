/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type VariableIdentifier struct {
	symbol          token.VariableSymbolToken
	moduleReference *ModuleReference
}

func NewVariableIdentifier(symbol token.VariableSymbolToken) *VariableIdentifier {
	return &VariableIdentifier{symbol: symbol}
}

func NewQualifiedVariableIdentifier(variable *VariableIdentifier, moduleReference *ModuleReference) *VariableIdentifier {
	return &VariableIdentifier{symbol: variable.symbol, moduleReference: moduleReference}
}

func (i *VariableIdentifier) Symbol() token.VariableSymbolToken {
	return i.symbol
}

func (i *VariableIdentifier) ModuleReference() *ModuleReference {
	return i.moduleReference
}

func (i *VariableIdentifier) FetchPositionLength() token.SourceFileReference {
	return i.symbol.SourceFileReference
}

func (i *VariableIdentifier) Name() string {
	if i.moduleReference != nil {
		return i.moduleReference.ModuleName() + "." + i.symbol.Name()
	}
	return i.symbol.Name()
}

func (i *VariableIdentifier) String() string {
	if i.moduleReference != nil {
		return i.moduleReference.ModuleName() + "." + i.symbol.String()
	}
	return i.symbol.String()
}

func (i *VariableIdentifier) DebugString() string {
	return fmt.Sprintf("[identifier]")
}
