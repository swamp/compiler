/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"github.com/swamp/compiler/src/token"
)

type TypeIdentifierScoped struct {
	identifier      *TypeIdentifier
	moduleReference *ModuleReference
}

func NewQualifiedTypeIdentifierScoped(moduleReference *ModuleReference, identifier *TypeIdentifier) *TypeIdentifierScoped {
	return &TypeIdentifierScoped{identifier: identifier, moduleReference: moduleReference}
}

func (i *TypeIdentifierScoped) ModuleReference() *ModuleReference {
	return i.moduleReference
}

func (i *TypeIdentifierScoped) Name() string {
	return i.moduleReference.ModuleName() + "." + i.identifier.Name()
}

func (i *TypeIdentifierScoped) Symbol() *TypeIdentifier {
	return i.identifier
}

func (i *TypeIdentifierScoped) IsDefaultSymbol() bool {
	return false
}

func (i *TypeIdentifierScoped) String() string {
	return i.moduleReference.ModuleName() + "." + i.identifier.String()
}

func (i *TypeIdentifierScoped) DebugString() string {
	return "[TypeIdentifierScoped]"
}

func (i *TypeIdentifierScoped) FetchPositionLength() token.SourceFileReference {
	return i.identifier.FetchPositionLength()
}
