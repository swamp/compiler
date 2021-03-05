package decorated

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

type TypeCreateAndLookup struct {
	lookup     *TypeLookup
	localTypes *ModuleTypes
}

func NewTypeCreateAndLookup(lookup *TypeLookup, localTypes *ModuleTypes) *TypeCreateAndLookup {
	return &TypeCreateAndLookup{localTypes: localTypes, lookup: lookup}
}

func (l *TypeCreateAndLookup) AddTypeAlias(alias *ast.Alias, concreteType dtype.Type, localComments []ast.LocalComment) (*dectype.Alias, TypeError) {
	return l.localTypes.AddTypeAlias(alias, concreteType, localComments)
}

func (l *TypeCreateAndLookup) AddCustomType(customType *dectype.CustomTypeAtom) TypeError {
	return l.localTypes.AddCustomType(customType)
}

func (l *TypeCreateAndLookup) CreateTypeReference(typeIdentifier *ast.TypeIdentifier) (*dectype.TypeReference, decshared.DecoratedError) {
	return l.lookup.CreateTypeReference(typeIdentifier)
}

func (l *TypeCreateAndLookup) CreateTypeScopedReference(typeIdentifier *ast.TypeIdentifierScoped) (*dectype.TypeReferenceScoped, decshared.DecoratedError) {
	return l.lookup.CreateTypeScopedReference(typeIdentifier)
}

func (l *TypeCreateAndLookup) FindBuiltInType(s string) dtype.Type {
	return l.localTypes.FindBuiltInType(s)
}

func (l *TypeCreateAndLookup) SourceModule() *Module {
	return l.localTypes.sourceModule
}
