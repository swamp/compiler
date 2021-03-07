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

func NewPackageLoader(filePathPrefix string, documentProvider loader.DocumentProvider, moduleNamespace dectype.PackageRootModuleName, world *loader.Package, worldDecorator *loader.WorldDecorator) *PackageLoader {
	fileLoader := loader.NewLoader(filePathPrefix, documentProvider)
	loaderAndDecorator := loader.NewModuleReaderAndDecorator(fileLoader, worldDecorator)
	moduleRepo := loader.NewModuleRepository(world, moduleNamespace, loaderAndDecorator)

	return &PackageLoader{repository: moduleRepo}
}
