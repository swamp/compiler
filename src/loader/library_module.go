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

type LibraryReaderAndDecorator struct{}

func NewLibraryReaderAndDecorator() *LibraryReaderAndDecorator {
	return &LibraryReaderAndDecorator{}
}

func FindSettingsDirectory(swampDirectory string) (string, error) {
	tryCount := 0
	testDirectory := swampDirectory
	if !file.IsDir(testDirectory) {
		return "", fmt.Errorf("wasn't a directory %q", swampDirectory)
	}
	for file.IsDir(testDirectory) && tryCount < 3 {
		if file.HasFile(filepath.Join(testDirectory, ".swamp.toml")) {
			return testDirectory, nil
		}
		testDirectory = filepath.Dir(testDirectory)
		tryCount++
	}

	return "", fmt.Errorf("sorry, couldn't find settings file from %v and up. Not a library", swampDirectory)
}

func (r *LibraryReaderAndDecorator) loadAndApplySettings(world *Package, repository deccy.ModuleRepository, swampDirectory string, documentProvider DocumentProvider, verboseFlag bool) decshared.DecoratedError {
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
		_, moduleErr := r.ReadLibraryModule(world, repository, dependencyFilePrefix, rootNamespace, documentProvider)
		if moduleErr != nil {
			return moduleErr
		}
	}

	return nil
}

func (r *LibraryReaderAndDecorator) ReadLibraryModule(world *Package, repository deccy.ModuleRepository, absoluteDirectory string, namespacePrefix dectype.PackageRootModuleName, documentProvider DocumentProvider) (*decorated.Module, decshared.DecoratedError) {
	const verboseFlag = false
	if strings.HasSuffix(absoluteDirectory, ".swamp") {
		panic("problem")
	}
	if verboseFlag {
		fmt.Printf("* read library %v -> %v  \n", namespacePrefix, absoluteDirectory)
	}
	const enforceStyle = true

	if err := r.loadAndApplySettings(world, repository, absoluteDirectory, documentProvider, verboseFlag); err != nil {
		return nil, err
	}

	fileLoader := NewLoader(absoluteDirectory, documentProvider)

	worldDecorator, worldDecoratorErr := NewWorldDecorator(enforceStyle, verboseFlag)
	if worldDecoratorErr != nil {
		return nil, nil
	}

	moduleReader := NewModuleReaderAndDecorator(fileLoader, worldDecorator)
	newRepository := NewModuleRepository(world, namespacePrefix, moduleReader)

	// green := color.New(color.FgHiGreen)
	// green.Fprintf(os.Stderr, "* compiling library '%v' {%v}\n", absoluteDirectory, namespacePrefix)

	return newRepository.FetchMainModuleInPackage(verboseFlag)
}

func (r *LibraryReaderAndDecorator) CompileAllInLibrary(world *Package, repository deccy.ModuleRepository, absoluteDirectory string, documentProvider DocumentProvider, namespacePrefix dectype.PackageRootModuleName) (*Package, decshared.DecoratedError) {
	const verboseFlag = false
	if strings.HasSuffix(absoluteDirectory, ".swamp") {
		panic("problem")
	}
	if verboseFlag {
		fmt.Printf("* read library %v -> %v  \n", namespacePrefix, absoluteDirectory)
	}
	const enforceStyle = true

	if err := r.loadAndApplySettings(world, repository, absoluteDirectory, documentProvider, verboseFlag); err != nil {
		return nil, err
	}

	fileLoader := NewLoader(absoluteDirectory, documentProvider)

	worldDecorator, worldDecoratorErr := NewWorldDecorator(enforceStyle, verboseFlag)
	if worldDecoratorErr != nil {
		return nil, nil
	}

	moduleReader := NewModuleReaderAndDecorator(fileLoader, worldDecorator)
	newRepository := NewModuleRepository(world, namespacePrefix, moduleReader)

	// green := color.New(color.FgHiGreen)
	// green.Fprintf(os.Stderr, "* compiling library '%v' {%v}\n", absoluteDirectory, namespacePrefix)

	_, err := newRepository.FetchMainModuleInPackage(verboseFlag)
	if err != nil {
		return nil, err
	}

	return world, nil
}

func (r *LibraryReaderAndDecorator) CompileAllInLibraryFindSettings(world *Package, repository deccy.ModuleRepository, absoluteDirectory string, documentProvider DocumentProvider, namespacePrefix dectype.PackageRootModuleName) (*Package, decshared.DecoratedError) {
	foundSettingsDirectory, err := FindSettingsDirectory(absoluteDirectory)
	if err != nil {
		return nil, decorated.NewInternalError(err)
	}

	return r.CompileAllInLibrary(world, repository, foundSettingsDirectory, documentProvider, namespacePrefix)
}
