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
