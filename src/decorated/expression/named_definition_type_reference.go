package decorated

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type NamedDefinitionTypeReference struct {
	optionalModuleReference *ModuleReference
	ident                   ast.TypeIdentifierNormalOrScoped
}

func NewNamedDefinitionTypeReference(optionalModuleReference *ModuleReference, ident ast.TypeIdentifierNormalOrScoped) *NamedDefinitionTypeReference {
	return &NamedDefinitionTypeReference{
		optionalModuleReference: optionalModuleReference,
		ident:                   ident,
	}
}

func (r *NamedDefinitionTypeReference) ModuleReference() *ModuleReference {
	return r.optionalModuleReference
}

func (r *NamedDefinitionTypeReference) AstIdentifier() ast.TypeIdentifierNormalOrScoped {
	return r.ident
}

func (r *NamedDefinitionTypeReference) String() string {
	return "named definition reference"
}

func (r *NamedDefinitionTypeReference) DebugString() string {
	return "named definition reference"
}

func (r *NamedDefinitionTypeReference) FetchPositionLength() token.SourceFileReference {
	return r.ident.FetchPositionLength()
}
