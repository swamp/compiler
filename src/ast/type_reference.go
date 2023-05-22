/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type TypeReferenceScopedOrNormal interface {
	Arguments() []Type
	FetchPositionLength() token.SourceFileReference
	SomeTypeIdentifier() TypeIdentifierNormalOrScoped
}

type TypeReference struct {
	ident     *TypeIdentifier `debug:"true"`
	arguments []Type          `debug:"true"`
}

func (i *TypeReference) String() string {
	if len(i.arguments) == 0 {
		return fmt.Sprintf("[TypeReference %v]", i.ident)
	}
	return fmt.Sprintf("[TypeReference %v %v]", i.ident, i.arguments)
}

func (i *TypeReference) DebugString() string {
	return ""
}

func (i *TypeReference) TypeIdentifier() *TypeIdentifier {
	return i.ident
}

func (i *TypeReference) SomeTypeIdentifier() TypeIdentifierNormalOrScoped {
	return i.ident
}

func (i *TypeReference) Arguments() []Type {
	return i.arguments
}

func (i *TypeReference) DecoratedName() string {
	return ""
}

func (i *TypeReference) Name() string {
	s := ""
	if len(i.arguments) == 0 {
		return fmt.Sprintf("%v", i.ident.Name())
	}

	for index, argument := range i.arguments {
		if index > 0 {
			s += " "
		}
		s += argument.Name()
	}
	return fmt.Sprintf("%v<%v>", i.ident.Name(), s)
}

func (i *TypeReference) FetchPositionLength() token.SourceFileReference {
	return i.ident.FetchPositionLength()
}

func NewTypeReference(ident *TypeIdentifier, arguments []Type) *TypeReference {
	return &TypeReference{ident: ident, arguments: arguments}
}
