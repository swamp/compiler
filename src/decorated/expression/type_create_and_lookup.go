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

func (l *TypeCreateAndLookup) AddTypeAlias(alias *dectype.Alias) TypeError {
	return l.localTypes.AddTypeAlias(alias)
}

func (l *TypeCreateAndLookup) AddCustomType(customType *dectype.CustomTypeAtom) TypeError {
	return l.localTypes.AddCustomType(customType)
}

func (l *TypeCreateAndLookup) CreateSomeTypeReference(someTypeIdentifier ast.TypeIdentifierNormalOrScoped) (dectype.TypeReferenceScopedOrNormal, decshared.DecoratedError) {
	return l.lookup.CreateSomeTypeReference(someTypeIdentifier)
}

func (l *TypeCreateAndLookup) FindBuiltInType(s string) dtype.Type {
	return l.localTypes.FindBuiltInType(s)
}

func (l *TypeCreateAndLookup) SourceModule() *Module {
	return l.localTypes.sourceModule
}
