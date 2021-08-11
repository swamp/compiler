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
	"github.com/swamp/compiler/src/environment"
	"github.com/swamp/compiler/src/file"
	"github.com/swamp/compiler/src/settings"
	"github.com/swamp/compiler/src/verbosity"
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

func ModuleTypeFromMapped(moduleMap settings.ModuleMap) decorated.ModuleType {
	switch moduleMap {
	case settings.ModuleFromEnvironment:
		return decorated.ModuleTypeFromEnvironment
	case settings.ModuleFromPath:
		return decorated.ModuleTypeFromPath
	}
	panic("unknown mapped")
}

func (r *LibraryReaderAndDecorator) loadAndApplySettings(world *Package, repository deccy.ModuleRepository, swampDirectory string, documentProvider DocumentProvider, configuration environment.Environment, verboseFlag verbosity.Verbosity) decshared.DecoratedError {
	settingsFilename := filepath.Join(swampDirectory, ".swamp.toml")

	mapping := make(map[string]settings.Module)

	settingsFile, settingsFileErr := os.Open(settingsFilename)
	if settingsFileErr != nil {
		return nil
	}

	settingsReader := bufio.NewReader(settingsFile)

	foundSettings, loadErr := settings.Load(settingsReader, swampDirectory, configuration)
	if loadErr != nil {
		return decorated.NewInternalError(loadErr)
	}

	for _, m := range foundSettings.Module {
		if verboseFlag >= verbosity.Mid {
			fmt.Printf("  * found mapping %s => %s\n", m.Name, m.Path)
		}

		mapping[m.Name] = m
	}

	for packageRootModuleNameString, packagePath := range mapping {
		dependencyFilePrefix := packagePath.Path
		if !filepath.IsAbs(packagePath.Path) {
			dependencyFilePrefix = filepath.Join(swampDirectory, packagePath.Path)
		}

		rootNamespace := dectype.MakePackageRootModuleNameFromString(packageRootModuleNameString)
		if !file.IsDir(dependencyFilePrefix) {
			full, _ := filepath.Abs(dependencyFilePrefix)
			return decorated.NewInternalError(fmt.Errorf("could not find directory '%v' '%v' ('%v' '%v')", full, dependencyFilePrefix, swampDirectory, packagePath))
		}
		_, moduleErr := r.ReadLibraryModule(ModuleTypeFromMapped(packagePath.Mapped), world, repository, dependencyFilePrefix, rootNamespace, documentProvider, configuration)
		if moduleErr != nil {
			return moduleErr
		}
	}

	return nil
}

func (r *LibraryReaderAndDecorator) ReadLibraryModule(moduleType decorated.ModuleType, world *Package, repository deccy.ModuleRepository, absoluteDirectory string, namespacePrefix dectype.PackageRootModuleName, documentProvider DocumentProvider, configuration environment.Environment) (*decorated.Module, decshared.DecoratedError) {
	const verboseFlag = verbosity.Low
	if strings.HasSuffix(absoluteDirectory, ".swamp") {
		panic("problem")
	}
	if verboseFlag >= verbosity.Mid {
		fmt.Printf("* read library %v -> %v  \n", namespacePrefix, absoluteDirectory)
	}
	const enforceStyle = true

	if err := r.loadAndApplySettings(world, repository, absoluteDirectory, documentProvider, configuration, verboseFlag); err != nil {
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

	return newRepository.FetchMainModuleInPackage(moduleType, verboseFlag)
}

func (r *LibraryReaderAndDecorator) CompileAllInLibrary(moduleType decorated.ModuleType, world *Package, repository deccy.ModuleRepository, absoluteDirectory string, documentProvider DocumentProvider, namespacePrefix dectype.PackageRootModuleName, configuration environment.Environment) (*Package, decshared.DecoratedError) {
	const verboseFlag = verbosity.None
	if strings.HasSuffix(absoluteDirectory, ".swamp") {
		panic("problem")
	}
	if verboseFlag >= verbosity.Mid {
		fmt.Printf("* read library %v -> %v  \n", namespacePrefix, absoluteDirectory)
	}
	const enforceStyle = true

	if err := r.loadAndApplySettings(world, repository, absoluteDirectory, documentProvider, configuration, verboseFlag); err != nil {
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

	_, err := newRepository.FetchMainModuleInPackage(moduleType, verboseFlag)
	if err != nil {
		return nil, err
	}

	return world, nil
}

func (r *LibraryReaderAndDecorator) CompileAllInLibraryFindSettings(moduleType decorated.ModuleType, world *Package, repository deccy.ModuleRepository, absoluteDirectory string, documentProvider DocumentProvider, namespacePrefix dectype.PackageRootModuleName, configuration environment.Environment) (*Package, decshared.DecoratedError) {
	foundSettingsDirectory, err := FindSettingsDirectory(absoluteDirectory)
	if err != nil {
		return nil, decorated.NewInternalError(err)
	}

	return r.CompileAllInLibrary(moduleType, world, repository, foundSettingsDirectory, documentProvider, namespacePrefix, configuration)
}
