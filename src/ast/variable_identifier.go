/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"github.com/swamp/compiler/src/token"
)

type ScopedOrNormalVariableIdentifier interface {
	Symbol() token.VariableSymbolToken
	FetchPositionLength() token.SourceFileReference
	Name() string
}

type VariableIdentifier struct {
	symbol token.VariableSymbolToken
}

func NewVariableIdentifier(symbol token.VariableSymbolToken) *VariableIdentifier {
	return &VariableIdentifier{symbol: symbol}
}

func (i *VariableIdentifier) Symbol() token.VariableSymbolToken {
	return i.symbol
}

func (i *VariableIdentifier) IsIgnore() bool {
	return i.symbol.IsIgnore()
}

func (i *VariableIdentifier) FetchPositionLength() token.SourceFileReference {
	return i.symbol.SourceFileReference
}

func (i *VariableIdentifier) Name() string {
	return i.symbol.Name()
}

func (i *VariableIdentifier) String() string {
	return i.symbol.String()
}

func (i *VariableIdentifier) DebugString() string {
	return "[identifier]"
}
