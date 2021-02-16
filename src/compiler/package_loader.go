/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package swampcompiler

import (
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/loader"
)

type PackageLoader struct {
	repository *loader.ModuleRepository
}

func NewPackageLoader(filePathPrefix string, moduleNamespace dectype.PackageRootModuleName, world *loader.World, worldDecorator *loader.WorldDecorator) *PackageLoader {
	fileLoader := loader.NewLoader(filePathPrefix)
	loaderAndDecorator := loader.NewModuleReaderAndDecorator(fileLoader, worldDecorator)
	moduleRepo := loader.NewModuleRepository(world, moduleNamespace, loaderAndDecorator)

	return &PackageLoader{repository: moduleRepo}
}
