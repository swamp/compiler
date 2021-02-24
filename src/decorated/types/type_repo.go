/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
)

type TypeRepo struct {
	identifierToType    map[string]dtype.Type
	types               []dtype.Type
	parentImportedTypes *ExposedTypes
	moduleName          ArtifactFullyQualifiedModuleName
}

func NewTypeRepo(moduleName ArtifactFullyQualifiedModuleName, imported *ExposedTypes) *TypeRepo {
	t := &TypeRepo{parentImportedTypes: imported, moduleName: moduleName, identifierToType: make(map[string]dtype.Type)}
	return t
}

func (t *TypeRepo) AllTypes() map[string]dtype.Type {
	return t.identifierToType
}

func (t *TypeRepo) AllLocalTypes() []dtype.Type {
	return t.types
}

func (t *TypeRepo) ImportedTypes() *ExposedTypes {
	return t.parentImportedTypes
}

func (t *TypeRepo) ModuleName() ArtifactFullyQualifiedModuleName {
	return t.moduleName
}

// -----------------------------------------------------
//                    Find
// -----------------------------------------------------
func (t *TypeRepo) FindTypeFromSignature(complete string) dtype.Type {
	found := t.identifierToType[complete]
	if found == nil {
		if t.parentImportedTypes != nil {
			return t.parentImportedTypes.FindTypeFromSignature(complete)
		}
	}
	return found
}

func (t *TypeRepo) FindType(typeToSearchFor dtype.Type) dtype.Type {
	return t.FindTypeFromSignature(typeToSearchFor.DecoratedName())
}

func (t *TypeRepo) FindTypeFromAlias(alias string) dtype.Type {
	return t.FindTypeFromSignature(alias)
}

func (t *TypeRepo) FindTypeFromName(alias string) dtype.Type {
	foundType := t.FindTypeFromSignature(alias)
	return foundType
}

func (t *TypeRepo) CreateTypeReference(typeIdentifier *ast.TypeIdentifier) dtype.Type {
	foundType := t.FindTypeFromSignature(typeIdentifier.Name())
	if foundType == nil {
		return nil
	}
	return NewTypeReference(typeIdentifier, foundType)
}

// -----------------------------------------------------
//                    Declare
// -----------------------------------------------------

func (t *TypeRepo) DeclareType(realType dtype.Type) error {
	existingType := t.FindType(realType)
	if existingType != nil {
		return fmt.Errorf("sorry, '%v' already declared", existingType)
	}
	t.internalAdd(realType)
	return nil
}

func (t *TypeRepo) AddFunctionAtom(astFunctionType *ast.FunctionType, parameterTypes []dtype.Type) *FunctionAtom {
	newType := NewFunctionAtom(astFunctionType, parameterTypes)
	existing := t.FindTypeFromSignature(newType.DecoratedName())
	if existing != nil {
		return existing.(*FunctionAtom)
	}
	t.internalAdd(newType)
	return newType
}

func (t *TypeRepo) DeclareRecordType(r *RecordAtom) *RecordAtom {
	combinedName := r.DecoratedName()
	existing := t.FindTypeFromSignature(combinedName)
	if existing != nil {
		return existing.(*RecordAtom)
	}

	t.internalAdd(r)
	return r
}

type DecoratedTypeError interface {
	Error() string
}

func (t *TypeRepo) DeclareTypeAlias(alias *ast.TypeIdentifier, concreteType dtype.Type) (*Alias, DecoratedTypeError) {
	artifactTypeName := t.moduleName.JoinTypeIdentifier(alias)
	newType := NewAliasType(alias, artifactTypeName, concreteType)
	t.internalAdd(newType)
	return newType, nil
}

func (t *TypeRepo) DeclareAlias(alias *ast.TypeIdentifier, referencedType dtype.Type, localComments []ast.LocalComment) (dtype.Type, DecoratedTypeError) {
	if referencedType == nil {
		panic("alias nil")
	}
	foundType := t.FindTypeFromAlias(alias.Name())
	if foundType != nil {
		//if foundType.AliasReferencedType() != referencedType {
		//return nil, NewDifferentAliasTypes(foundType, referencedType)
		//}
		return foundType, nil
	}

	return t.DeclareTypeAlias(alias, referencedType)
}

// -----------------------------------------------------
//                    Other
// -----------------------------------------------------

func (t *TypeRepo) DebugOutput() {
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

func (t *TypeRepo) DebugString() string {
	s := "Type Repo:\n"
	for k, v := range t.types {
		s += fmt.Sprintf(".. %p %v : %v\n", v, k, TraverseToString(v))
	}

	s += "Imported types:\n"

	s += t.parentImportedTypes.DebugString()
	return s
}

func (t *TypeRepo) String() string {
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
	_, isPrimitive := repoType.(*PrimitiveAtom)
	if isPrimitive {
		return true
	}

	if stringInSlice(repoType.DecoratedName(), []string{"Nothing", "Just", "Maybe(a)", "Maybe", "List", "Array", "TypeRef", "{test:", "{atest:"}) {
		return true
	}
	return false
}

func (t *TypeRepo) ShortString() string {
	s := ""
	for _, repoType := range t.types {
		if isTypeToIgnoreForDebugOutput(repoType) {
			continue
		}

		s += fmt.Sprintf("%v : %v\n", repoType.DecoratedName(), repoType.ShortString())
	}

	return s
}

func (t *TypeRepo) internalAdd(realType dtype.Type) {
	t.internalAddWithString(realType.DecoratedName(), realType)
}

func (t *TypeRepo) internalAddWithString(name string, realType dtype.Type) {
	// fmt.Printf("Adding type %v\n", realType.Name())
	hasType := t.identifierToType[name]
	if hasType != nil {
		panic("already have name " + name)
	}
	t.identifierToType[name] = realType
	t.types = append(t.types, realType)
}

func (t *TypeRepo) CopyTypes(realTypes []dtype.Type) error {
	for _, copyType := range realTypes {
		copyErr := t.CopyType(copyType)
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

func (t *TypeRepo) CopyType(realType dtype.Type) error {
	existingType := t.FindTypeFromName(realType.DecoratedName())
	if existingType != nil {
		return fmt.Errorf("copy: sorry, '%v' already declared", existingType)
	}
	t.internalAddWithString(realType.DecoratedName(), realType)
	return nil
}
