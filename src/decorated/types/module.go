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

type ModuleName struct {
	path []*ast.TypeIdentifier
}

func MakeModuleName(path []*ast.TypeIdentifier) ModuleName {
	for _, p := range path {
		if p == nil {
			panic("nil")
		}
	}
	return ModuleName{path: path}
}

func MakeModuleNameFromString(name string) ModuleName {
	parts := strings.Split(name, ".")
	var typeIdents []*ast.TypeIdentifier
	for _, part := range parts {
		typeIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken(part, token.PositionLength{}, 0))
		typeIdents = append(typeIdents, typeIdentifier)
	}

	return ModuleName{typeIdents}
}

func (m ModuleName) String() string {
	s := ""

	for index, p := range m.path {
		if index > 0 {
			s += "."
		}

		s += p.Name()
	}

	return s
}

func (m ModuleName) Path() []*ast.TypeIdentifier {
	return m.path
}

type PackageRelativeModuleName struct {
	ModuleName
}

func NewPackageRelativeModuleName(path []*ast.TypeIdentifier) PackageRelativeModuleName {
	return PackageRelativeModuleName{ModuleName: MakeModuleName(path)}
}

type ArtifactFullyQualifiedModuleName struct {
	ModuleName
}

type ArtifactFullyQualifiedTypeName struct {
	ModuleName
}

func MakeArtifactFullyQualifiedModuleName(path []*ast.TypeIdentifier) ArtifactFullyQualifiedModuleName {
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
	return SingleModuleName{MakeModuleName([]*ast.TypeIdentifier{path})}
}

func (m SingleModuleName) IsEmpty() bool {
	return len(m.ModuleName.path) == 0
}

func (n *ArtifactFullyQualifiedModuleName) String() string {
	return n.ModuleName.String()
}

func MakePackageRootModuleName(path []*ast.TypeIdentifier) PackageRootModuleName {
	return PackageRootModuleName{MakeModuleName(path)}
}

func MakePackageRootModuleNameFromString(name string) PackageRootModuleName {
	return PackageRootModuleName{MakeModuleNameFromString(name)}
}

func (n PackageRootModuleName) Join(relative PackageRelativeModuleName) ArtifactFullyQualifiedModuleName {
	paths := append(n.ModuleName.path, relative.path...)
	return MakeArtifactFullyQualifiedModuleName(paths)
}

func (n ArtifactFullyQualifiedModuleName) JoinLocalName(relative *ast.VariableIdentifier) string {
	if n.String() == "" {
		return relative.Name()
	}

	str := n.String() + "." + relative.Name()

	return str
}

func (n ArtifactFullyQualifiedModuleName) JoinTypeIdentifier(relative *ast.TypeIdentifier) ArtifactFullyQualifiedTypeName {
	var newPaths []*ast.TypeIdentifier

	newPaths = append(newPaths, n.path...)

	if relative.ModuleReference() != nil {
		for _, mRef := range relative.ModuleReference().Parts() {
			newPaths = append(newPaths, mRef.TypeIdentifier())
		}
	}
	newPaths = append(newPaths, relative)

	return ArtifactFullyQualifiedTypeName{ModuleName{path: newPaths}}
}

func MakePackageRelativeModuleName(path []*ast.TypeIdentifier) PackageRelativeModuleName {
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
