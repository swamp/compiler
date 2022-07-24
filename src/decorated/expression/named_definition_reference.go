package decorated

import (
	"fmt"
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

func (r *NamedDefinitionReference) FullyQualified() *FullyQualifiedPackageVariableName {
	return NewFullyQualifiedVariableName(r.optionalModuleReference.module, ast.NewVariableIdentifier(r.ident.Symbol()))
}

func (r *NamedDefinitionReference) FullyQualifiedName() string {
	if r.optionalModuleReference != nil {
		return r.optionalModuleReference.module.fullyQualifiedModuleName.String() + "." + r.ident.Symbol().Name()
	}

	return r.ident.Name()
}

func moduleRefToString(reference *ModuleReference) string {
	if reference == nil {
		return ""
	}
	return reference.module.fullyQualifiedModuleName.String()
}

func (r *NamedDefinitionReference) String() string {
	return fmt.Sprintf("[NamedDefinitionReference %v/%v]", moduleRefToString(r.optionalModuleReference), r.ident.Symbol().Name())
}

func (r *NamedDefinitionReference) DebugString() string {
	return fmt.Sprintf("[NamedDefinitionReference %v/%v]", moduleRefToString(r.optionalModuleReference), r.ident.Symbol().Name())
}

func (r *NamedDefinitionReference) FetchPositionLength() token.SourceFileReference {
	return r.ident.FetchPositionLength()
}
