/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type TypeIdentifierNormalOrScoped interface {
	IsDefaultSymbol() bool
	FetchPositionLength() token.SourceFileReference
}

type TypeIdentifier struct {
	symbolToken token.TypeSymbolToken
}

func NewTypeIdentifier(symbolToken token.TypeSymbolToken) *TypeIdentifier {
	return &TypeIdentifier{symbolToken: symbolToken}
}

func (i *TypeIdentifier) Name() string {
	return i.symbolToken.Name()
}

func (i *TypeIdentifier) Symbol() token.TypeSymbolToken {
	return i.symbolToken
}

func (i *TypeIdentifier) String() string {
	return i.symbolToken.String()
}

func (i *TypeIdentifier) DebugString() string {
	return fmt.Sprintf("[TypeIdentifier]")
}

func (i *TypeIdentifier) IsDefaultSymbol() bool {
	return i.symbolToken.Name() == "_"
}

func (i *TypeIdentifier) FetchPositionLength() token.SourceFileReference {
	return i.symbolToken.SourceFileReference
}
