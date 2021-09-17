/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/verbosity"
)

type WorldDecorator struct {
	rootModules   []*decorated.Module
	importModules []*decorated.Module
	verbose       verbosity.Verbosity
	forceStyle    bool
}

func NewWorldDecorator(forceStyle bool, verbose verbosity.Verbosity) (*WorldDecorator, decshared.DecoratedError) {
	rootModules, importModules, rootModuleErr := deccy.CreateDefaultRootModule(false)
	if rootModuleErr != nil {
		return nil, rootModuleErr
	}
	return &WorldDecorator{verbose: verbose, forceStyle: forceStyle, rootModules: rootModules, importModules: importModules}, nil
}

func (w *WorldDecorator) RootModules() []*decorated.Module {
	return w.rootModules
}

func (w *WorldDecorator) ImportModules() []*decorated.Module {
	return w.importModules
}

func (w *WorldDecorator) RunesToModule(moduleType decorated.ModuleType, moduleRepository deccy.ModuleRepository, moduleName dectype.ArtifactFullyQualifiedModuleName, absoluteFilename string, str string) (*decorated.Module, decshared.DecoratedError) {
	const errorsAsWarnings = false
	return deccy.InternalCompileToModule(moduleType, moduleRepository, w.rootModules, w.importModules, moduleName, absoluteFilename, str,
		w.forceStyle, w.verbose, errorsAsWarnings)
}
