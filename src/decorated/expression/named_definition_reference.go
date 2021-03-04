package decorated

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

type NamedDefinitionReference struct {
	optionalModuleReference *ModuleReference
	ident                   *ast.VariableIdentifier
}

func NewNamedDefinitionReference(optionalModuleReference *ModuleReference, ident *ast.VariableIdentifier) *NamedDefinitionReference {
	return &NamedDefinitionReference{
		optionalModuleReference: optionalModuleReference,
		ident:                   ident,
	}
}

func (r *NamedDefinitionReference) ModuleReference() *ModuleReference {
	return r.optionalModuleReference
}

func (r *NamedDefinitionReference) AstIdentifier() *ast.VariableIdentifier {
	return r.ident
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
