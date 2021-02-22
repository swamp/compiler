/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type TypeIdentifier struct {
	symbolToken     token.TypeSymbolToken
	moduleReference *ModuleReference
}

func NewTypeIdentifier(symbolToken token.TypeSymbolToken) *TypeIdentifier {
	return &TypeIdentifier{symbolToken: symbolToken}
}

func NewQualifiedTypeIdentifier(symbolToken token.TypeSymbolToken, moduleReference *ModuleReference) *TypeIdentifier {
	return &TypeIdentifier{symbolToken: symbolToken, moduleReference: moduleReference}
}

func (i *TypeIdentifier) ModuleReference() *ModuleReference {
	return i.moduleReference
}

func (i *TypeIdentifier) Name() string {
	if i.moduleReference != nil {
		return i.moduleReference.ModuleName() + "." + i.symbolToken.Name()
	}
	return i.symbolToken.Name()
}

func (i *TypeIdentifier) Symbol() token.TypeSymbolToken {
	return i.symbolToken
}

func (i *TypeIdentifier) String() string {
	if i.moduleReference != nil {
		return i.moduleReference.ModuleName() + "." + i.symbolToken.String()
	}
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
