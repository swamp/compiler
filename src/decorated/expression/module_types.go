/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"log"
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type ModuleTypes struct {
	identifierToType map[string]dtype.Type
	sourceModule     *Module
}

func NewModuleTypes(sourceModule *Module) *ModuleTypes {
	t := &ModuleTypes{sourceModule: sourceModule, identifierToType: make(map[string]dtype.Type)}
	return t
}

func (t *ModuleTypes) AllTypes() map[string]dtype.Type {
	return t.identifierToType
}

func (t *ModuleTypes) SourceModule() *Module {
	return t.sourceModule
}

// -----------------------------------------------------
//                    Find
// -----------------------------------------------------
//
func (t *ModuleTypes) FindType(identifier *ast.TypeIdentifier) dtype.Type {
	found := t.identifierToType[identifier.Name()]
	return found
}

func (t *ModuleTypes) FindBuiltInType(name string) dtype.Type {
	found := t.identifierToType[name]
	return found
}

func (t *ModuleTypes) CreateTypeReference(typeIdentifier *ast.TypeIdentifier) dtype.Type {
	foundType := t.FindType(typeIdentifier)
	if foundType == nil {
		return nil
	}
	ref := dectype.NewTypeReference(typeIdentifier, foundType)

	return ref
}

// -----------------------------------------------------
//                    Declare
// -----------------------------------------------------

func (t *ModuleTypes) internalAddType(typeIdentifier *ast.TypeIdentifier, realType dtype.Type) error {
	existingType := t.FindType(typeIdentifier)
	if existingType != nil {
		return fmt.Errorf("sorry, '%v' already declared", existingType)
	}
	t.internalAdd(typeIdentifier, realType)
	return nil
}

func (t *ModuleTypes) InternalAddPrimitive(typeIdentifier *ast.TypeIdentifier, atom *dectype.PrimitiveAtom) error {
	return t.internalAddType(typeIdentifier, atom)
}

type TypeError interface {
	Error() string
}

func (t *ModuleTypes) AddTypeAlias(alias *ast.Alias, concreteType dtype.Type, localComments []ast.LocalComment) (*dectype.Alias, TypeError) {
	artifactTypeName := t.sourceModule.fullyQualifiedModuleName.JoinTypeIdentifier(alias.Identifier())
	newType := dectype.NewAliasType(alias, artifactTypeName, concreteType)
	t.internalAddType(alias.Identifier(), newType)
	return newType, nil
}

func (t *ModuleTypes) internalAddVariantConstructorType(constructor *dectype.CustomTypeVariantConstructorType) {
	t.internalAddType(constructor.Variant().Name(), constructor)
}

func (t *ModuleTypes) addCustomTypeVariantConstructors(customType *dectype.CustomTypeAtom) {
	for _, variant := range customType.Variants() {
		constructorType := dectype.NewCustomTypeVariantConstructorType(variant)
		t.internalAddVariantConstructorType(constructorType)
	}
}

func (t *ModuleTypes) AddCustomType(customType *dectype.CustomTypeAtom) TypeError {
	t.internalAddType(customType.TypeIdentifier(), customType)
	t.addCustomTypeVariantConstructors(customType)
	return nil
}

// -----------------------------------------------------
//                    Other
// -----------------------------------------------------

func (t *ModuleTypes) DebugOutput() {
	fmt.Println(t.DebugString())
}

func TraverseToString(t dtype.Type) string {
	s := t.DecoratedName()
	concreteType, isConcreteType := t.(dtype.Type)
	if isConcreteType {
		next := concreteType.Next()
		if next != nil {
			s += " => " + TraverseToString(next)
		}
	}
	return s
}

func (t *ModuleTypes) DebugString() string {
	s := "Type Repo:\n"
	for k, v := range t.identifierToType {
		s += fmt.Sprintf(".. %p %v : %v\n", v, k, TraverseToString(v))
	}
	return s
}

func (t *ModuleTypes) String() string {
	return t.DebugString()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.HasPrefix(a, b) {
			return true
		}
	}
	return false
}

func isTypeToIgnoreForDebugOutput(repoType dtype.Type) bool {
	_, isPrimitive := repoType.(*dectype.PrimitiveAtom)
	if isPrimitive {
		return true
	}

	if stringInSlice(repoType.DecoratedName(), []string{"Nothing", "Just", "Maybe(a)", "Maybe", "List", "Array", "TypeRef", "{test:", "{atest:"}) {
		return true
	}
	return false
}

func (t *ModuleTypes) ShortString() string {
	s := ""
	for _, repoType := range t.identifierToType {
		if isTypeToIgnoreForDebugOutput(repoType) {
			continue
		}

		s += fmt.Sprintf("%v : %v\n", repoType.DecoratedName(), repoType.ShortString())
	}

	return s
}

func (t *ModuleTypes) internalAdd(identifier *ast.TypeIdentifier, realType dtype.Type) {
	t.internalAddWithString(identifier.Name(), realType)
}

func (t *ModuleTypes) internalAddWithString(name string, realType dtype.Type) {
	log.Printf("Adding type %v\n", name)
	hasType := t.identifierToType[name]
	if hasType != nil {
		panic("already have name " + name)
	}
	t.identifierToType[name] = realType
}

func (t *ModuleTypes) CopyTypes(realTypes map[string]dtype.Type) decshared.DecoratedError {
	for nameOfType, copyType := range realTypes {
		symbol := token.NewTypeSymbolToken(nameOfType, t.sourceModule.FetchPositionLength(), 0)
		fakeIdentifier := ast.NewTypeIdentifier(symbol)
		copyErr := t.CopyType(fakeIdentifier, copyType)
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

func (t *ModuleTypes) CopyType(nameOfType *ast.TypeIdentifier, realType dtype.Type) decshared.DecoratedError {
	existingType := t.FindType(nameOfType)
	if existingType != nil {
		return NewInternalError(fmt.Errorf("copy: sorry, '%v' already declared", existingType))
	}

	log.Printf("copying %v, %T", nameOfType.Name(), realType)

	t.internalAddWithString(nameOfType.Name(), realType)
	return nil
}
