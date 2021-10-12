package decorated

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type NamedDefinitionReference struct {
	optionalModuleReference *ModuleReference
	ident                   ast.ScopedOrNormalVariableIdentifier
}

func NewNamedDefinitionReference(optionalModuleReference *ModuleReference, ident ast.ScopedOrNormalVariableIdentifier) *NamedDefinitionReference {
	return &NamedDefinitionReference{
		optionalModuleReference: optionalModuleReference,
		ident:                   ident,
	}
}

func (r *NamedDefinitionReference) ModuleReference() *ModuleReference {
	return r.optionalModuleReference
}

func (r *NamedDefinitionReference) AstIdentifier() ast.ScopedOrNormalVariableIdentifier {
	return r.ident
}

func (r *NamedDefinitionReference) FullyQualifiedName() string {
	if r.optionalModuleReference != nil {
		return r.optionalModuleReference.module.fullyQualifiedModuleName.String() + "." + r.ident.Symbol().Name()
	}

	return r.ident.Name()
}

func (r *NamedDefinitionReference) String() string {
	return "named definition reference"
}

func (r *NamedDefinitionReference) DebugString() string {
	return "named definition reference"
}

func (r *NamedDefinitionReference) FetchPositionLength() token.SourceFileReference {
	return r.ident.FetchPositionLength()
}
