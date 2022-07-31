/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"log"

	"github.com/swamp/compiler/src/parser"

	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/verbosity"
)

type RunesToModuleConverter interface {
	RunesToModule(moduleType decorated.ModuleType, moduleRepository deccy.ModuleRepository, moduleName dectype.ArtifactFullyQualifiedModuleName, relativeFilename string, str string) (*decorated.Module, decshared.DecoratedError)
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

func (r *ModuleReaderAndDecorator) ReadModule(moduleType decorated.ModuleType, repository deccy.ModuleRepository, moduleName dectype.PackageRelativeModuleName, namespacePrefix dectype.PackageRootModuleName) (*decorated.Module, decshared.DecoratedError) {
	const verboseFlag = verbosity.None
	if verboseFlag > verbosity.None {
		log.Printf("* read module %v\n", moduleName)
	}

	fullyQualifiedName := namespacePrefix.Join(moduleName)

	absoluteFilename, runes, loadErr := r.runesLoader.Load(moduleName, verboseFlag)
	var errors []decshared.DecoratedError
	if loadErr != nil {
		if parser.IsCompileErr(loadErr) {
			return nil, loadErr
		}
		errors = append(errors, loadErr)
	}

	// green := color.New(color.FgHiGreen)
	// filepathToShow := pathutil.TryToMakeRelativePath(absoluteFilename)
	// green.Fprintf(os.Stderr, "* compiling module '%v' %v\n", filepathToShow, fullyQualifiedName)

	loadedModule, runesErr := r.runesToModule.RunesToModule(moduleType, repository, fullyQualifiedName, absoluteFilename, runes)
	if runesErr != nil {
		if parser.IsCompileErr(runesErr) {
			return nil, runesErr
		}
		errors = append(errors, runesErr)
	}

	var returnErr decshared.DecoratedError

	if len(errors) > 0 {
		returnErr = decorated.NewMultiErrors(errors)
	}

	return loadedModule, returnErr
}
