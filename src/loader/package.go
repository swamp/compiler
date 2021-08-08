/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"fmt"

	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

type Package struct {
	moduleLookup       map[string]*decorated.Module
	absolutePathLookup map[LocalFileSystemPath]*decorated.Module
	modules            []*decorated.Module
	root               LocalFileSystemRoot
	name               string
}

func NewPackage(root LocalFileSystemRoot, name string) *Package {
	return &Package{root: root, name: name, moduleLookup: make(map[string]*decorated.Module), absolutePathLookup: make(map[LocalFileSystemPath]*decorated.Module)}
}

func (w *Package) Root() LocalFileSystemRoot {
	return w.root
}

func (w *Package) Name() string {
	return w.name
}

func (w *Package) AllModules() []*decorated.Module {
	return w.modules
}

func (w *Package) FindModule(moduleName dectype.ArtifactFullyQualifiedModuleName) *decorated.Module {
	return w.moduleLookup[moduleName.String()]
}

func (w *Package) FindModuleFromAbsoluteFilePath(absolutePath LocalFileSystemPath) *decorated.Module {
	return w.absolutePathLookup[absolutePath]
}

func (w *Package) AddModule(moduleName dectype.ArtifactFullyQualifiedModuleName, module *decorated.Module) {
	if module == nil {
		panic("not a good module")
	}
	if _, hasExisting := w.moduleLookup[moduleName.String()]; hasExisting {
		panic(fmt.Errorf("tried to add an already existing module '%v'", moduleName))
	}

	w.moduleLookup[moduleName.String()] = module

	localFilePath, convertErr := module.Document().Uri.ToLocalFilePath()
	if convertErr != nil {
		panic(convertErr)
	}
	localFilePathForThisModule := LocalFileSystemPath(localFilePath)
	w.absolutePathLookup[localFilePathForThisModule] = module
	w.modules = append(w.modules, module)
}

func (w *Package) String() string {
	s := ""
	for key, module := range w.modules {
		s += fmt.Sprintf("\n\n===== MODULE %v =======\n", key)
		s += module.String()
	}
	return s
}
