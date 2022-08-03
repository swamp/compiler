/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/swamp/compiler/src/parser"

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
		if file.HasFile(path.Join(testDirectory, ".swamp.toml")) {
			return testDirectory, nil
		}
		testDirectory = path.Dir(testDirectory)
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
	settingsFilename := path.Join(swampDirectory, ".swamp.toml")

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
			log.Printf("  * found mapping %s => %s\n", m.Name, m.Path)
		}
	}

	var errors decshared.DecoratedError

	for _, packagePath := range foundSettings.Module {
		dependencyFilePrefix := packagePath.Path
		if !filepath.IsAbs(packagePath.Path) {
			dependencyFilePrefix = path.Join(swampDirectory, packagePath.Path)
		}

		rootNamespace := dectype.MakePackageRootModuleNameFromString(packagePath.Name)
		dependencyFilePrefix = filepath.ToSlash(dependencyFilePrefix)
		if !file.IsDir(dependencyFilePrefix) {
			full, err := filepath.Abs(dependencyFilePrefix)
			if err != nil {
				return decorated.NewInternalError(fmt.Errorf("could not do abs of '%w'", err))
			}
			full = filepath.ToSlash(full)
			return decorated.NewInternalError(fmt.Errorf("could not find slashed directory '%v' '%v' ('%v' '%v')", full, dependencyFilePrefix, swampDirectory, packagePath))
		}
		_, moduleErr := r.ReadLibraryModule(ModuleTypeFromMapped(packagePath.Mapped), world, repository, dependencyFilePrefix, rootNamespace, documentProvider, configuration)
		if moduleErr != nil {
			if parser.IsCompileError(moduleErr) {
				return moduleErr
			}
			errors = decorated.AppendError(errors, moduleErr)
		}
	}

	return errors
}

func (r *LibraryReaderAndDecorator) ReadLibraryModule(moduleType decorated.ModuleType, world *Package, repository deccy.ModuleRepository, absoluteDirectory string, namespacePrefix dectype.PackageRootModuleName, documentProvider DocumentProvider, configuration environment.Environment) (*decorated.Module, decshared.DecoratedError) {
	const verboseFlag = verbosity.Low
	if strings.HasSuffix(absoluteDirectory, ".swamp") {
		panic(fmt.Sprintf("the directory should not end with .swamp '%v'", absoluteDirectory))
	}
	if verboseFlag >= verbosity.Mid {
		log.Printf("* read library %v -> %v  \n", namespacePrefix, absoluteDirectory)
	}
	const enforceStyle = true

	var errors decshared.DecoratedError

	if err := r.loadAndApplySettings(world, repository, absoluteDirectory, documentProvider, configuration, verboseFlag); err != nil {
		if parser.IsCompileErr(err) {
			return nil, err
		}
		errors = decorated.AppendError(errors, err)
	}

	fileLoader := NewLoader(absoluteDirectory, documentProvider)

	worldDecorator, worldDecoratorErr := NewWorldDecorator(enforceStyle, verboseFlag)
	errors = decorated.AppendError(errors, worldDecoratorErr)
	if parser.IsCompileErr(worldDecoratorErr) {
		return nil, worldDecoratorErr
	}

	moduleReader := NewModuleReaderAndDecorator(fileLoader, worldDecorator)
	newRepository := NewModuleRepository(world, namespacePrefix, moduleReader)

	// green := color.New(color.FgHiGreen)
	// green.Fprintf(os.Stderr, "* compiling library '%v' {%v}\n", absoluteDirectory, namespacePrefix)

	fetchedModule, fetchErr := newRepository.FetchMainModuleInPackage(moduleType, verboseFlag)
	if parser.IsCompileErr(fetchErr) {
		return nil, fetchErr
	}
	errors = decorated.AppendError(errors, fetchErr)

	return fetchedModule, errors
}

func (r *LibraryReaderAndDecorator) CompileAllInLibrary(moduleType decorated.ModuleType, world *Package, repository deccy.ModuleRepository, absoluteDirectory string, documentProvider DocumentProvider, namespacePrefix dectype.PackageRootModuleName, configuration environment.Environment) (*Package, decshared.DecoratedError) {
	const verboseFlag = verbosity.None
	if strings.HasSuffix(absoluteDirectory, ".swamp") {
		panic("problem")
	}
	if verboseFlag >= verbosity.Mid {
		log.Printf("* read library %v -> %v  \n", namespacePrefix, absoluteDirectory)
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
		if parser.IsCompileErr(err) {
			return nil, err
		}
	}

	return world, err
}

func (r *LibraryReaderAndDecorator) CompileAllInLibraryFindSettings(moduleType decorated.ModuleType, world *Package, repository deccy.ModuleRepository, absoluteDirectory string, documentProvider DocumentProvider, namespacePrefix dectype.PackageRootModuleName, configuration environment.Environment) (*Package, decshared.DecoratedError) {
	foundSettingsDirectory, err := FindSettingsDirectory(absoluteDirectory)
	if err != nil {
		return nil, decorated.NewInternalError(err)
	}

	return r.CompileAllInLibrary(moduleType, world, repository, foundSettingsDirectory, documentProvider, namespacePrefix, configuration)
}
