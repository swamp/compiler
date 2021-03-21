/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type AliasReference struct {
	ident     *NamedDefinitionTypeReference
	reference *Alias
}

func (g *AliasReference) NamedTypeRef() *NamedDefinitionTypeReference {
	return g.ident
}

func (g *AliasReference) String() string {
	return fmt.Sprintf("[letvarref %v %v]", g.ident, g.reference)
}

func (g *AliasReference) HumanReadable() string {
	return "Alias reference"
}

func (g *AliasReference) Alias() *Alias {
	return g.reference
}

func NewAliasReference(ident *NamedDefinitionTypeReference, reference *Alias) *AliasReference {
	if reference == nil {
		panic("cant be nil")
	}

	ref := &AliasReference{ident: ident, reference: reference}

	reference.AddReferee(ref)

	return ref
}

func (g *AliasReference) Resolve() (dtype.Atom, error) {
	return g.reference.Resolve()
}

func (g *AliasReference) Next() dtype.Type {
	return g.reference.Next()
}

func (g *AliasReference) ParameterCount() int {
	return g.reference.ParameterCount()
}

func (g *AliasReference) FetchPositionLength() token.SourceFileReference {
	return g.ident.FetchPositionLength()
}
