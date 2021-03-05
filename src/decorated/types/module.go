/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
)

// ModuleName is similar to a ModuleReference, however path can be nil
type ModuleName struct {
	path          *ast.ModuleReference
	precalculated string
}

func MakeModuleName(path *ast.ModuleReference) ModuleName {
	if path != nil {
		for _, p := range path.Parts() {
			if p == nil {
				panic("nil")
			}
		}
	}
	return ModuleName{path: path, precalculated: CalculateString(path)}
}

func MakeModuleNameFromString(name string) ModuleName {
	parts := strings.Split(name, ".")
	if len(parts) == 0 {
		return MakeModuleName(nil)
	}
	var nameParts []*ast.ModuleNamePart
	for _, part := range parts {
		typeIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken(part, token.SourceFileReference{}, 0))
		part := ast.NewModuleNamePart(typeIdentifier)
		nameParts = append(nameParts, part)
	}

	return MakeModuleName(ast.NewModuleReference(nameParts))
}

func CalculateString(ref *ast.ModuleReference) string {
	s := ""

	if ref != nil {
		for index, p := range ref.Parts() {
			if index > 0 {
				s += "."
			}

			s += p.TypeIdentifier().Name()
		}
	}

	return s
}

func (m ModuleName) String() string {
	return m.precalculated
}

func (m ModuleName) Path() *ast.ModuleReference {
	return m.path
}

type PackageRelativeModuleName struct {
	ModuleName
}

func NewPackageRelativeModuleName(path *ast.ModuleReference) PackageRelativeModuleName {
	return PackageRelativeModuleName{ModuleName: MakeModuleName(path)}
}

type ArtifactFullyQualifiedModuleName struct {
	ModuleName
}

type ArtifactFullyQualifiedTypeName struct {
	ModuleName
}

func MakeArtifactFullyQualifiedModuleName(path *ast.ModuleReference) ArtifactFullyQualifiedModuleName {
	return ArtifactFullyQualifiedModuleName{MakeModuleName(path)}
}

type PackageRootModuleName struct {
	ModuleName
}

type SingleModuleName struct {
	ModuleName
}

func MakeSingleModuleName(path *ast.TypeIdentifier) SingleModuleName {
	if path == nil {
		return SingleModuleName{MakeModuleName(nil)}
	}
	return SingleModuleName{MakeModuleName(ast.NewModuleReference([]*ast.ModuleNamePart{ast.NewModuleNamePart(path)}))}
}

func (m SingleModuleName) IsEmpty() bool {
	return m.path == nil
}

func (n *ArtifactFullyQualifiedModuleName) String() string {
	return n.ModuleName.String()
}

func MakePackageRootModuleName(path *ast.ModuleReference) PackageRootModuleName {
	return PackageRootModuleName{MakeModuleName(path)}
}

func MakePackageRootModuleNameFromString(name string) PackageRootModuleName {
	return PackageRootModuleName{MakeModuleNameFromString(name)}
}

func (n PackageRootModuleName) Join(relative PackageRelativeModuleName) ArtifactFullyQualifiedModuleName {
	var existingParts []*ast.ModuleNamePart

	var partsToAdd []*ast.ModuleNamePart

	if n.path != nil {
		existingParts = n.path.Parts()
	}

	if relative.Path() != nil {
		partsToAdd = relative.Path().Parts()
	}

	paths := append(existingParts, partsToAdd...)

	var modRef *ast.ModuleReference

	if len(paths) > 0 {
		modRef = ast.NewModuleReference(paths)
	}

	return MakeArtifactFullyQualifiedModuleName(modRef)
}

func (n ArtifactFullyQualifiedModuleName) JoinLocalName(relative *ast.VariableIdentifier) string {
	if n.String() == "" {
		return relative.Name()
	}

	str := n.String() + "." + relative.Name()

	return str
}

func (n ArtifactFullyQualifiedModuleName) JoinTypeIdentifier(relative *ast.TypeIdentifier) ArtifactFullyQualifiedTypeName {
	var newPaths []*ast.ModuleNamePart

	if n.path != nil {
		newPaths = append(newPaths, n.path.Parts()...)
	}

	newPaths = append(newPaths, ast.NewModuleNamePart(relative))

	return ArtifactFullyQualifiedTypeName{ModuleName{path: ast.NewModuleReference(newPaths)}}
}

func (n ArtifactFullyQualifiedModuleName) JoinTypeIdentifierScoped(relative *ast.TypeIdentifierScoped) ArtifactFullyQualifiedTypeName {
	var newPaths []*ast.ModuleNamePart

	if n.path != nil {
		newPaths = append(newPaths, n.path.Parts()...)
	}

	for _, mRef := range relative.ModuleReference().Parts() {
		newPaths = append(newPaths, mRef)
	}
	newPaths = append(newPaths, ast.NewModuleNamePart(relative.Symbol()))

	return ArtifactFullyQualifiedTypeName{ModuleName{path: ast.NewModuleReference(newPaths)}}
}

func MakePackageRelativeModuleName(path *ast.ModuleReference) PackageRelativeModuleName {
	return PackageRelativeModuleName{MakeModuleName(path)}
}

func MakePackageRelativeModuleNameFromString(name string) PackageRelativeModuleName {
	return PackageRelativeModuleName{MakeModuleNameFromString(name)}
}

func (n PackageRelativeModuleName) JoinLocalName(relative *ast.VariableIdentifier) string {
	if n.String() == "" {
		return relative.Name()
	}

	str := n.String() + "." + relative.Name()

	// fmt.Printf("packageRelative and add local name %v + %v = %v\n", n, relative, str)

	return str
}
