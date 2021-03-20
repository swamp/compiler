package decorated

import (
	"fmt"
	"log"

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

func (l *TypeLookup) FindType(typeIdentifier *ast.TypeIdentifier) (dtype.Type, *dectype.NamedDefinitionTypeReference, decshared.DecoratedError) {
	ref := ast.NewTypeReference(typeIdentifier, nil)
	namedTypeRef := dectype.NewNamedDefinitionTypeReference(nil, ref)

	foundLocalType := l.localTypes.FindType(typeIdentifier)
	if foundLocalType != nil {
		return foundLocalType, namedTypeRef, nil
	}

	foundLocalType = l.importedTypes.FindType(typeIdentifier)
	if foundLocalType == nil {
		return nil, nil, NewInternalError(fmt.Errorf("could not find imported type %v", typeIdentifier))
	}

	return foundLocalType, namedTypeRef, nil
}

func (l *TypeLookup) FindTypeScoped(typeIdentifier *ast.TypeIdentifierScoped) (dtype.Type, *dectype.NamedDefinitionTypeReference, decshared.DecoratedError) {
	moduleFound := l.moduleImports.FindModule(typeIdentifier.ModuleReference())
	if moduleFound == nil {
		return nil, nil, NewInternalError(fmt.Errorf("could not find module %v", typeIdentifier.ModuleReference()))
	}

	moduleReference := NewModuleReference(typeIdentifier.ModuleReference(), moduleFound)
	typeReferenceScoped := ast.NewScopedTypeReference(typeIdentifier, nil)
	namedTypeRef := dectype.NewNamedDefinitionTypeReference(moduleReference, typeReferenceScoped)
	foundExposedType := moduleFound.ExposedTypes().FindType(typeIdentifier.Symbol())
	if foundExposedType == nil {
		return nil, nil, NewInternalError(fmt.Errorf("could not find exposed type %v in module %v", typeIdentifier, moduleFound))
	}

	return foundExposedType, namedTypeRef, nil
}

/*
func (l *TypeLookup) CreateTypeReference(typeIdentifier *ast.TypeIdentifier) (*dectype.TypeReference, *dectype.NamedDefinitionTypeReference, decshared.DecoratedError) {
	foundType, namedRef, err := l.FindType(typeIdentifier)
	if err != nil {
		return nil, namedRef, err
	}

	switch t := foundType.(type) {
	case *dectype.Alias:
		dectype.NewAliasReference(typeIdentifier, t)
	}
	typeRef := dectype.NewTypeReference(typeIdentifier, foundType)

	return typeRef, namedRef, nil
}

func (l *TypeLookup) CreateTypeScopedReference(typeIdentifier *ast.TypeIdentifierScoped) (*dectype.TypeReferenceScoped, *dectype.NamedDefinitionTypeReference, decshared.DecoratedError) {
	foundType, module, err := l.FindTypeScoped(typeIdentifier)
	if err != nil {
		return nil, nil, err
	}

	return dectype.NewTypeReferenceScoped(typeIdentifier, foundType), module, nil
}

*/

func (l *TypeLookup) CreateSomeTypeReference(someTypeIdentifier ast.TypeIdentifierNormalOrScoped) (dectype.TypeReferenceScopedOrNormal, decshared.DecoratedError) {
	var lookedUpType dtype.Type
	var named *dectype.NamedDefinitionTypeReference

	scoped, wasScope := someTypeIdentifier.(*ast.TypeIdentifierScoped)
	if wasScope {
		foundType, foundNamed, err := l.FindTypeScoped(scoped)
		if err != nil {
			return nil, err
		}
		lookedUpType = foundType
		named = foundNamed
	} else {
		normal, wasNormal := someTypeIdentifier.(*ast.TypeIdentifier)
		if !wasNormal {
			return nil, NewInternalError(fmt.Errorf("not sure of this type identifier %T", someTypeIdentifier))
		}

		foundType, foundNamed, typeErr := l.FindType(normal)
		if typeErr != nil {
			return nil, typeErr
		}
		lookedUpType = foundType
		named = foundNamed
	}

	var reference dectype.TypeReferenceScopedOrNormal
	switch t := lookedUpType.(type) {
	case *dectype.Alias:
		reference = dectype.NewAliasReference(named, t)
	case *dectype.CustomTypeVariant:
		reference = dectype.NewCustomTypeVariantReference(named, t)
	case *dectype.PrimitiveAtom:
		reference = dectype.NewPrimitiveTypeReference(named, t)
	case *dectype.CustomTypeAtom:
		reference = dectype.NewCustomTypeReference(named, t)
	case *dectype.CustomTypeVariantConstructorType:
		reference = dectype.NewCustomTypeVariantReference(named, t.Variant())
	default:
		log.Printf("typelookup: what is this type: %T", t)
	}

	return reference, nil
}
