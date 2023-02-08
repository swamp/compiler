/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type TypeCreateAndLookup struct {
	lookup           *TypeLookup
	localTypeLookup  *dectype.TypeParameterContext
	localNameContext *dectype.LocalTypeNameContext
	localTypes       *ModuleTypes
}

func NewTypeCreateAndLookup(lookup *TypeLookup, localTypes *ModuleTypes, localTypeContext *dectype.TypeParameterContext) *TypeCreateAndLookup {
	return &TypeCreateAndLookup{localTypes: localTypes, lookup: lookup, localTypeLookup: localTypeContext}
}

func (l *TypeCreateAndLookup) MakeLocalNameContext(localNameContext *dectype.LocalTypeNameContext) TypeAddAndReferenceMaker {
	if localNameContext == nil {
		panic(fmt.Errorf("can not set to nil local name context"))
	}

	return &TypeCreateAndLookup{localTypes: l.localTypes, lookup: l.lookup, localTypeLookup: l.localTypeLookup, localNameContext: localNameContext}
}

func (l *TypeCreateAndLookup) AddTypeAlias(alias *dectype.Alias) TypeError {
	return l.localTypes.AddTypeAlias(alias)
}

func (l *TypeCreateAndLookup) AddCustomType(customType *dectype.CustomTypeAtom) TypeError {
	return l.localTypes.AddCustomType(customType)
}

func (l *TypeCreateAndLookup) AddCustomTypeWrappedInNameOnlyContext(customTypeWrappedInContext *dectype.LocalTypeNameContext) TypeError {
	return l.localTypes.AddCustomTypeWrappedInNameOnlyContext(customTypeWrappedInContext)
}

func (l *TypeCreateAndLookup) CreateSomeTypeReference(someTypeIdentifier ast.TypeIdentifierNormalOrScoped) (dectype.TypeReferenceScopedOrNormal, decshared.DecoratedError) {
	return l.lookup.CreateSomeTypeReference(someTypeIdentifier)
}

func (l *TypeCreateAndLookup) FindBuiltInType(s string) dtype.Type {
	identifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken(s, token.NewInternalSourceFileReference(), 0))
	foundType, _, err := l.lookup.FindType(identifier)
	if err != nil {
		panic(fmt.Errorf("could not find %v", identifier))
	}

	return foundType
}

func (l *TypeCreateAndLookup) CreateLocalTypeNameOnlyReference(some *ast.LocalTypeNameReference) (*dectype.LocalTypeNameReference, decshared.DecoratedError) {
	if l.localNameContext == nil {
		panic(fmt.Errorf("no localNameContext set"))
	}
	found, err := l.localNameContext.ReferenceNameOnly(some)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return found, nil
}

func (l *TypeCreateAndLookup) LookupLocalTypeName(some *ast.LocalTypeName) (dtype.Type, decshared.DecoratedError) {
	found, err := l.localTypeLookup.LookupTypeName(some)
	if err != nil {
		return nil, NewInternalError(err)
	}
	return found, nil
}

func (l *TypeCreateAndLookup) SourceModule() *Module {
	return l.localTypes.sourceModule
}
