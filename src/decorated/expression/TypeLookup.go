package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

type TypeLookup struct {
	moduleImports *ModuleImports
	localTypes    *ModuleTypes
	importedTypes *ExposedTypes
}

func NewTypeLookup(moduleImports *ModuleImports, localTypes *ModuleTypes, importedTypes *ExposedTypes) *TypeLookup {
	return &TypeLookup{
		moduleImports: moduleImports,
		localTypes:    localTypes,
		importedTypes: importedTypes,
	}
}

func (l *TypeLookup) FindType(typeIdentifier *ast.TypeIdentifier) (dtype.Type, decshared.DecoratedError) {
	foundLocalType := l.localTypes.FindType(typeIdentifier)
	if foundLocalType != nil {
		return foundLocalType, nil
	}

	foundLocalType = l.importedTypes.FindType(typeIdentifier)
	if foundLocalType == nil {
		return nil, NewInternalError(fmt.Errorf("could not find imported type %v", typeIdentifier))
	}

	return foundLocalType, nil
}

func (l *TypeLookup) FindTypeScoped(typeIdentifier *ast.TypeIdentifierScoped) (dtype.Type, decshared.DecoratedError) {
	moduleFound := l.moduleImports.FindModule(typeIdentifier.ModuleReference())
	if moduleFound == nil {
		return nil, NewInternalError(fmt.Errorf("could not find module %v", typeIdentifier.ModuleReference()))
	}

	foundExposedType := moduleFound.ExposedTypes().FindType(typeIdentifier.Symbol())
	if foundExposedType == nil {
		return nil, NewInternalError(fmt.Errorf("could not find exposed type %v in module %v", typeIdentifier, moduleFound))
	}

	return foundExposedType, nil
}

func (l *TypeLookup) CreateTypeReference(typeIdentifier *ast.TypeIdentifier) (*dectype.TypeReference, decshared.DecoratedError) {
	foundType, err := l.FindType(typeIdentifier)
	if err != nil {
		return nil, err
	}

	return dectype.NewTypeReference(typeIdentifier, foundType), nil
}

func (l *TypeLookup) CreateTypeScopedReference(typeIdentifier *ast.TypeIdentifierScoped) (*dectype.TypeReferenceScoped, decshared.DecoratedError) {
	foundType, err := l.FindTypeScoped(typeIdentifier)
	if err != nil {
		return nil, err
	}

	return dectype.NewTypeReferenceScoped(typeIdentifier, foundType), nil
}

func (l *TypeLookup) CreateSomeTypeReference(someTypeIdentifier ast.TypeIdentifierNormalOrScoped) (dectype.TypeReferenceScopedOrNormal, decshared.DecoratedError) {
	scoped, wasScope := someTypeIdentifier.(*ast.TypeIdentifierScoped)
	if wasScope {
		return l.CreateTypeScopedReference(scoped)
	}

	normal, wasNormal := someTypeIdentifier.(*ast.TypeIdentifier)
	if !wasNormal {
		return nil, NewInternalError(fmt.Errorf("not sure of this type identifier %T", someTypeIdentifier))
	}

	return l.CreateTypeReference(normal)
}
