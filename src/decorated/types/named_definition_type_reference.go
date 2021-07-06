package dectype

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/modref"
	"github.com/swamp/compiler/src/token"
)

type NamedDefinitionTypeReference struct {
	optionalModuleReference modref.ModuleReferencer
	ident                   ast.TypeReferenceScopedOrNormal
}

func NewNamedDefinitionTypeReference(optionalModuleReference modref.ModuleReferencer, ident ast.TypeReferenceScopedOrNormal) *NamedDefinitionTypeReference {
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

func (r *NamedDefinitionTypeReference) String() string {
	return "named definition type reference"
}

func (r *NamedDefinitionTypeReference) DebugString() string {
	return "named definition type reference"
}

func (r *NamedDefinitionTypeReference) FetchPositionLength() token.SourceFileReference {
	return r.ident.FetchPositionLength()
}
