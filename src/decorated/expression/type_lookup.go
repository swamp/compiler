/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

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

	importedType := l.importedTypes.FindType(typeIdentifier)
	if importedType == nil {
		return nil, nil, NewUnknownImportedType(typeIdentifier)
	}

	importedType.MarkAsReferenced()

	foundLocalType = importedType.referencedType

	return foundLocalType, namedTypeRef, nil
}

func (l *TypeLookup) FindTypeScoped(typeIdentifier *ast.TypeIdentifierScoped) (dtype.Type, *dectype.NamedDefinitionTypeReference, decshared.DecoratedError) {
	moduleFound := l.moduleImports.FindModule(typeIdentifier.ModuleReference())
	if moduleFound == nil {
		return nil, nil, NewUnknownModule(typeIdentifier.ModuleReference())
	}

	NewImportStatementReference(moduleFound.ImportStatementInModule())

	moduleFound.MarkAsReferenced()
	moduleReference := NewModuleReference(typeIdentifier.ModuleReference(), moduleFound.referencedModule)
	typeReferenceScoped := ast.NewScopedTypeReference(typeIdentifier, nil)
	namedTypeRef := dectype.NewNamedDefinitionTypeReference(moduleReference, typeReferenceScoped)
	foundExposedType := moduleFound.referencedModule.ExposedTypes().FindType(typeIdentifier.Symbol())
	if foundExposedType == nil {
		return nil, nil, NewUnknownExposedType(typeIdentifier.Symbol())
	}
	foundExposedType.MarkAsReferenced()

	return foundExposedType.referencedType, namedTypeRef, nil
}

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
			return nil, NewInternalError(fmt.Errorf("not sure of this type functionParameter %T", someTypeIdentifier))
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
	case *dectype.CustomTypeVariantAtom:
		reference = dectype.NewCustomTypeVariantReference(named, t)
	case *dectype.PrimitiveAtom:
		reference = dectype.NewPrimitiveTypeReference(named, t)
	case *dectype.CustomTypeAtom:
		reference = dectype.NewCustomTypeReference(named, t)
	case *dectype.LocalTypeNameOnlyContext:
		reference = dectype.NewLocalTypeNameContextReference(named, t)
		// TODO: FIX THIS
		//panic(fmt.Errorf("what is this %v", named))
	default:
		log.Printf("typelookup: what is this type: %T", t)
	}

	//	log.Printf("reference is %T %v", reference, someTypeIdentifier)

	return reference, nil
}
