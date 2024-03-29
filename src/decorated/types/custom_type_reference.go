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

type CustomTypeReference struct {
	named      *NamedDefinitionTypeReference
	customType *CustomTypeAtom
}

func (g *CustomTypeReference) Type() dtype.Type {
	return g.customType
}

func (g *CustomTypeReference) String() string {
	return fmt.Sprintf("[VariantRef %v]", g.named)
}

func (g *CustomTypeReference) Next() dtype.Type {
	return g.customType
}

func (g *CustomTypeReference) HumanReadable() string {
	return g.customType.HumanReadable()
}

func (g *CustomTypeReference) CustomTypeAtom() *CustomTypeAtom {
	return g.customType
}

func (g *CustomTypeReference) AstIdentifier() ast.TypeReferenceScopedOrNormal {
	return g.named.ident
}

func (g *CustomTypeReference) NameReference() *NamedDefinitionTypeReference {
	return g.named
}

func (g *CustomTypeReference) ParameterCount() int {
	return g.customType.ParameterCount()
}

func (g *CustomTypeReference) Resolve() (dtype.Atom, error) {
	return g.customType.Resolve()
}

func NewCustomTypeReference(named *NamedDefinitionTypeReference, customType *CustomTypeAtom) *CustomTypeReference {
	ref := &CustomTypeReference{named: named, customType: customType}

	customType.AddReferee(ref)

	return ref
}

func (g *CustomTypeReference) FetchPositionLength() token.SourceFileReference {
	return g.named.FetchPositionLength()
}

func (g *CustomTypeReference) WasReferenced() bool {
	return false // You can not reference references
}
