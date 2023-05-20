/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/modref"
	"github.com/swamp/compiler/src/token"
)

type NamedDefinitionTypeReference struct {
	optionalModuleReference modref.ModuleReferencer
	ident                   ast.TypeReferenceScopedOrNormal
}

func NewNamedDefinitionTypeReference(optionalModuleReference modref.ModuleReferencer, ident ast.TypeReferenceScopedOrNormal) *NamedDefinitionTypeReference {
	if !ident.FetchPositionLength().Verify() {
		//panic(fmt.Errorf("stop, wrong type %T %v", ident, ident))
	}
	return &NamedDefinitionTypeReference{
		optionalModuleReference: optionalModuleReference,
		ident:                   ident,
	}
}

func (r *NamedDefinitionTypeReference) ModuleReference() modref.ModuleReferencer {
	return r.optionalModuleReference
}

func (r *NamedDefinitionTypeReference) AstIdentifier() ast.TypeReferenceScopedOrNormal {
	return r.ident
}

func moduleReferenceString(referencer modref.ModuleReferencer) string {
	if referencer == nil {
		return ""
	}
	return referencer.String()
}

func (r *NamedDefinitionTypeReference) String() string {
	return fmt.Sprintf("[NamedDefTypeRef %v:%v]", moduleReferenceString(r.optionalModuleReference), r.ident)
}

func (r *NamedDefinitionTypeReference) DebugString() string {
	return fmt.Sprintf("[NamedDefTypeRef %v:%v]", moduleReferenceString(r.optionalModuleReference), r.ident)
}

func (r *NamedDefinitionTypeReference) FetchPositionLength() token.SourceFileReference {
	return r.ident.FetchPositionLength()
}

func MakeFakeNamedDefinitionTypeReference(sourceFileRef token.SourceFileReference, primitiveName string) *NamedDefinitionTypeReference {
	typeIdent := ast.NewTypeIdentifier(token.NewTypeSymbolToken(primitiveName, sourceFileRef, 0))
	typeRef := ast.NewTypeReference(typeIdent, nil)
	return NewNamedDefinitionTypeReference(nil, typeRef)
}
