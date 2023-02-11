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

type LocalTypeNameContextReference struct {
	named       *NamedDefinitionTypeReference
	nameContext *LocalTypeNameContext
}

func (g *LocalTypeNameContextReference) Type() dtype.Type {
	return g.nameContext
}

func (g *LocalTypeNameContextReference) String() string {
	return fmt.Sprintf("[LocalTypeNameContextReference %v]", g.named)
}

func (g *LocalTypeNameContextReference) Next() dtype.Type {
	return g.nameContext
}

func (g *LocalTypeNameContextReference) HumanReadable() string {
	return g.nameContext.HumanReadable()
}

func (g *LocalTypeNameContextReference) LocalTypeNameContext() *LocalTypeNameContext {
	return g.nameContext
}

func (g *LocalTypeNameContextReference) AstIdentifier() ast.TypeReferenceScopedOrNormal {
	return g.named.ident
}

func (g *LocalTypeNameContextReference) NameReference() *NamedDefinitionTypeReference {
	return g.named
}

func (g *LocalTypeNameContextReference) ParameterCount() int {
	return g.nameContext.ParameterCount()
}

func (g *LocalTypeNameContextReference) Resolve() (dtype.Atom, error) {
	return g.nameContext.Resolve()
}

func NewLocalTypeNameContextReference(named *NamedDefinitionTypeReference, context *LocalTypeNameContext) *LocalTypeNameContextReference {
	ref := &LocalTypeNameContextReference{named: named, nameContext: context}

	// TODO: context.AddReferee(ref)

	return ref
}

func (g *LocalTypeNameContextReference) FetchPositionLength() token.SourceFileReference {
	return g.named.FetchPositionLength()
}

func (g *LocalTypeNameContextReference) WasReferenced() bool {
	return false // You can not reference references
}
