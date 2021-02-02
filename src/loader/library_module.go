/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/file"
	"github.com/swamp/compiler/src/settings"
)

type LibraryReaderAndDecorator struct {
}

func NewLibraryReaderAndDecorator() *LibraryReaderAndDecorator {
	return &LibraryReaderAndDecorator{}
}

func (r *LibraryReaderAndDecorator) checkSettings(world *World, repository deccy.ModuleRepository, swampDirectory string, verboseFlag bool) decshared.DecoratedError {
	settingsFilename := filepath.Join(swampDirectory, ".swamp.toml")

	mapping := make(map[string]string)

	settingsFile, settingsFileErr := os.Open(settingsFilename)
	if settingsFileErr != nil {
		return nil
	}

	settingsReader := bufio.NewReader(settingsFile)
	foundSettings, loadErr := settings.Load(settingsReader, swampDirectory)
	if loadErr != nil {
		return decorated.NewInternalError(loadErr)
	}
	for _, m := range foundSettings.Module {
		if verboseFlag {
			fmt.Printf("  * found mapping %s => %s\n", m.Name, m.Path)
		}
		mapping[m.Name] = m.Path
	}

	for packageRootModuleNameString, packagePath := range mapping {
		dependencyFilePrefix := packagePath
		if !filepath.IsAbs(packagePath) {
			dependencyFilePrefix = filepath.Join(swampDirectory, packagePath)
		}

		rootNamespace := dectype.MakePackageRootModuleNameFromString(packageRootModuleNameString)
		if !file.IsDir(dependencyFilePrefix) {
			panic(fmt.Sprintf("could not find directory '%v' '%v'", swampDirectory, packagePath))
			// moduleNameToFetch = dectype.MakePackageRelativeModuleName(nil)
		}
		_, moduleErr := r.ReadLibraryModule(world, repository, dependencyFilePrefix, rootNamespace)
		if moduleErr != nil {
			return moduleErr
		}
	}

	return nil
}

func (r *LibraryReaderAndDecorator) ReadLibraryModule(world *World, repository deccy.ModuleRepository, absoluteDirectory string, namespacePrefix dectype.PackageRootModuleName) (*decorated.Module, decshared.DecoratedError) {
	const verboseFlag = false
	if strings.HasSuffix(absoluteDirectory, ".swamp") {
		panic("problem")
	}
	if verboseFlag {
		fmt.Printf("* read library %v -> %v  \n", namespacePrefix, absoluteDirectory)
	}
	const enforceStyle = true

	if err := r.checkSettings(world, repository, absoluteDirectory, verboseFlag); err != nil {
		return nil, err
	}

	fileLoader := NewLoader(absoluteDirectory)

	worldDecorator, worldDecoratorErr := NewWorldDecorator(enforceStyle, verboseFlag)
	if worldDecoratorErr != nil {
		return nil, nil
	}

	moduleReader := NewModuleReaderAndDecorator(fileLoader, worldDecorator)
	newRepository := NewModuleRepository(world, namespacePrefix, moduleReader)

	//green := color.New(color.FgHiGreen)
	//green.Fprintf(os.Stderr, "* compiling library '%v' {%v}\n", absoluteDirectory, namespacePrefix)

	return newRepository.FetchMainModuleInPackage(verboseFlag)
}
