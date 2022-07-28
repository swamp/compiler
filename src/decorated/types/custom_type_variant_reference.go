/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type CustomTypeVariantReference struct {
	named             *NamedDefinitionTypeReference
	customTypeVariant *CustomTypeVariantAtom
}

func (g *CustomTypeVariantReference) Type() dtype.Type {
	return g.customTypeVariant
}

func (g *CustomTypeVariantReference) String() string {
	return fmt.Sprintf("[VariantRef %v %v]", g.named, g.customTypeVariant)
}

func (g *CustomTypeVariantReference) Next() dtype.Type {
	return g.customTypeVariant
}

func (g *CustomTypeVariantReference) HumanReadable() string {
	return g.customTypeVariant.HumanReadable()
}

func (g *CustomTypeVariantReference) NameReference() *NamedDefinitionTypeReference {
	return g.named
}

func (g *CustomTypeVariantReference) CustomTypeVariant() *CustomTypeVariantAtom {
	return g.customTypeVariant
}

func (g *CustomTypeVariantReference) AstIdentifier() ast.TypeReferenceScopedOrNormal {
	return g.named.ident
}

func (g *CustomTypeVariantReference) ParameterCount() int {
	return g.customTypeVariant.ParameterCount()
}

func (g *CustomTypeVariantReference) Resolve() (dtype.Atom, error) {
	return g.customTypeVariant.Resolve()
}

func NewCustomTypeVariantReference(named *NamedDefinitionTypeReference, customTypeVariant *CustomTypeVariantAtom) *CustomTypeVariantReference {
	ref := &CustomTypeVariantReference{named: named, customTypeVariant: customTypeVariant}

	customTypeVariant.AddReferee(ref)

	return ref
}

func (g *CustomTypeVariantReference) FetchPositionLength() token.SourceFileReference {
	return g.named.FetchPositionLength()
}

func (g *CustomTypeVariantReference) WasReferenced() bool {
	return false
}
