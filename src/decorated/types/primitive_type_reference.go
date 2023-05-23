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

type PrimitiveTypeReference struct {
	named         *NamedDefinitionTypeReference
	primitiveType *PrimitiveAtom `debug:"true"`
}

func (g *PrimitiveTypeReference) Type() dtype.Type {
	return g.primitiveType
}

func (g *PrimitiveTypeReference) NameReference() *NamedDefinitionTypeReference {
	return g.named
}

func (g *PrimitiveTypeReference) String() string {
	return fmt.Sprintf("[PrimitiveTypeRef %v]", g.named)
}

func (g *PrimitiveTypeReference) Next() dtype.Type {
	return g.primitiveType
}

func (g *PrimitiveTypeReference) HumanReadable() string {
	return g.primitiveType.HumanReadable()
}

func (g *PrimitiveTypeReference) PrimitiveAtom() *PrimitiveAtom {
	return g.primitiveType
}

func (g *PrimitiveTypeReference) AstIdentifier() ast.TypeReferenceScopedOrNormal {
	return g.named.ident
}

func (g *PrimitiveTypeReference) ParameterCount() int {
	return g.primitiveType.ParameterCount()
}

func (g *PrimitiveTypeReference) Resolve() (dtype.Atom, error) {
	return g.primitiveType.Resolve()
}

func NewPrimitiveTypeReference(named *NamedDefinitionTypeReference,
	primitiveType *PrimitiveAtom) *PrimitiveTypeReference {
	if named == nil {
		panic(fmt.Errorf("must have a named type reference"))
	}

	if !named.FetchPositionLength().Verify() {
		panic("not allowed")
	}
	ref := &PrimitiveTypeReference{named: named, primitiveType: primitiveType}

	primitiveType.AddReferee(ref)

	return ref
}

func (g *PrimitiveTypeReference) FetchPositionLength() token.SourceFileReference {
	return g.named.FetchPositionLength()
}

func (g *PrimitiveTypeReference) WasReferenced() bool {
	return false
}
