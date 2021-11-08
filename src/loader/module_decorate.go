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
	rootModule *decorated.Module
	verbose    verbosity.Verbosity
	forceStyle bool
}

func NewWorldDecorator(forceStyle bool, verbose verbosity.Verbosity) (*WorldDecorator, decshared.DecoratedError) {
	rootModule, rootModuleErr := deccy.CreateDefaultRootModule(true)
	if rootModuleErr != nil {
		return nil, rootModuleErr
	}
	return &WorldDecorator{verbose: verbose, forceStyle: forceStyle, rootModule: rootModule}, nil
}

func (w *WorldDecorator) RootModules() *decorated.Module {
	return w.rootModule
}

func (w *WorldDecorator) RunesToModule(moduleType decorated.ModuleType, moduleRepository deccy.ModuleRepository, moduleName dectype.ArtifactFullyQualifiedModuleName, absoluteFilename string, str string) (*decorated.Module, decshared.DecoratedError) {
	const errorsAsWarnings = false
	return deccy.InternalCompileToModule(moduleType, moduleRepository, w.rootModule, moduleName, absoluteFilename, str,
		w.forceStyle, w.verbose, errorsAsWarnings)
}
