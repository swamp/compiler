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

type LocalTypeNameOnlyContextReference struct {
	named       *NamedDefinitionTypeReference `debug:"true"`
	nameContext *LocalTypeNameOnlyContext     `debug:"true"`
}

func (g *LocalTypeNameOnlyContextReference) Type() dtype.Type {
	return g.nameContext
}

func (g *LocalTypeNameOnlyContextReference) String() string {
	return fmt.Sprintf("[LocalTypeNameOnlyContextReference %v %v]", g.named, g.nameContext)
}

func (g *LocalTypeNameOnlyContextReference) Next() dtype.Type {
	return g.nameContext
}

func (g *LocalTypeNameOnlyContextReference) HumanReadable() string {
	return g.nameContext.HumanReadable()
}

func (g *LocalTypeNameOnlyContextReference) LocalTypeNameContext() *LocalTypeNameOnlyContext {
	return g.nameContext
}

func (g *LocalTypeNameOnlyContextReference) AstIdentifier() ast.TypeReferenceScopedOrNormal {
	return g.named.ident
}

func (g *LocalTypeNameOnlyContextReference) NameReference() *NamedDefinitionTypeReference {
	return g.named
}

func (g *LocalTypeNameOnlyContextReference) ParameterCount() int {
	return g.nameContext.ParameterCount()
}

func (g *LocalTypeNameOnlyContextReference) Resolve() (dtype.Atom, error) {
	return g.nameContext.Resolve()
}

func NewLocalTypeNameContextReference(named *NamedDefinitionTypeReference, context *LocalTypeNameOnlyContext) *LocalTypeNameOnlyContextReference {
	ref := &LocalTypeNameOnlyContextReference{named: named, nameContext: context}

	// TODO: context.AddReferee(ref)

	return ref
}

func (g *LocalTypeNameOnlyContextReference) FetchPositionLength() token.SourceFileReference {
	return g.named.FetchPositionLength()
}

func (g *LocalTypeNameOnlyContextReference) WasReferenced() bool {
	return false // You can not reference references
}
