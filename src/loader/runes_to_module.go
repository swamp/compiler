/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"fmt"

	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/verbosity"
)

type RunesToModuleConverter interface {
	RunesToModule(moduleRepository deccy.ModuleRepository, moduleName dectype.ArtifactFullyQualifiedModuleName, relativeFilename string, str string) (*decorated.Module, decshared.DecoratedError)
}

type ModuleRunes interface {
	Load(moduleName dectype.PackageRelativeModuleName, verboseFlag verbosity.Verbosity) (string, string, decshared.DecoratedError)
}

type ModuleReaderAndDecorator struct {
	runesLoader   ModuleRunes
	runesToModule RunesToModuleConverter
}

func NewModuleReaderAndDecorator(runesLoader ModuleRunes, runesToModule RunesToModuleConverter) *ModuleReaderAndDecorator {
	return &ModuleReaderAndDecorator{runesLoader: runesLoader, runesToModule: runesToModule}
}

func (r *ModuleReaderAndDecorator) ReadModule(repository deccy.ModuleRepository, moduleName dectype.PackageRelativeModuleName, namespacePrefix dectype.PackageRootModuleName) (*decorated.Module, decshared.DecoratedError) {
	const verboseFlag = verbosity.None
	if verboseFlag > verbosity.None {
		fmt.Printf("* read module %v\n", moduleName)
	}

	fullyQualifiedName := namespacePrefix.Join(moduleName)

	absoluteFilename, runes, loadErr := r.runesLoader.Load(moduleName, verboseFlag)
	if loadErr != nil {
		return nil, loadErr
	}

	// green := color.New(color.FgHiGreen)
	// filepathToShow := pathutil.TryToMakeRelativePath(absoluteFilename)
	// green.Fprintf(os.Stderr, "* compiling module '%v' %v\n", filepathToShow, fullyQualifiedName)

	return r.runesToModule.RunesToModule(repository, fullyQualifiedName, absoluteFilename, runes)
}
