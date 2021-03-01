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

type World struct {
	moduleLookup       map[string]*decorated.Module
	absolutePathLookup map[LocalFileSystemPath]*decorated.Module
	modules            []*decorated.Module
}

func NewWorld() *World {
	return &World{moduleLookup: make(map[string]*decorated.Module), absolutePathLookup: make(map[LocalFileSystemPath]*decorated.Module)}
}

func (w *World) AllModules() []*decorated.Module {
	return w.modules
}

func (w *World) FindModule(moduleName dectype.ArtifactFullyQualifiedModuleName) *decorated.Module {
	return w.moduleLookup[moduleName.String()]
}

func (w *World) FindModuleFromAbsoluteFilePath(absolutePath LocalFileSystemPath) *decorated.Module {
	return w.absolutePathLookup[absolutePath]
}

func (w *World) AddModule(moduleName dectype.ArtifactFullyQualifiedModuleName, module *decorated.Module) {
	if module == nil {
		panic("not a good module")
	}
	w.moduleLookup[moduleName.String()] = module

	localFilePath, convertErr := module.DocumentURI().ToLocalFilePath()
	if convertErr != nil {
		panic(convertErr)
	}
	w.absolutePathLookup[LocalFileSystemPath(localFilePath)] = module
	w.modules = append(w.modules, module)
}

func (w *World) String() string {
	s := ""
	for key, module := range w.modules {
		s += fmt.Sprintf("\n\n===== MODULE %v =======\n", key)
		s += module.String()
	}
	return s
}
