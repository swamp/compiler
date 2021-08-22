/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package settings

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/swamp/compiler/src/environment"
)

type ModuleMap = int

const (
	ModuleFromEnvironment ModuleMap = iota
	ModuleFromPath
)

type Module struct {
	Name   string
	Path   string
	Mapped ModuleMap
}

type Settings struct {
	Name   string
	Module []Module
}

func Load(reader io.Reader, rootDirectory string, configuration environment.Environment) (Settings, error) {
	data, dataErr := ioutil.ReadAll(reader)
	if dataErr != nil {
		return Settings{}, dataErr
	}

	settings := Settings{}
	mapped := ModuleFromPath

	unmarshalErr := toml.Unmarshal(data, &settings)
	if unmarshalErr != nil {
		return Settings{}, unmarshalErr
	}

	for index, mod := range settings.Module {
		cleanedUpPath := strings.TrimSpace(mod.Path)
		convertedPath := cleanedUpPath
		if strings.HasPrefix(cleanedUpPath, "${") {
			endIndex := strings.Index(cleanedUpPath, "}")
			if endIndex == -1 {
				return Settings{}, fmt.Errorf("bad format '%v'", cleanedUpPath)
			}
			packageName := cleanedUpPath[2:endIndex]
			suffix := cleanedUpPath[endIndex+1:]
			convertedPath = configuration.Lookup(packageName)
			if convertedPath == "" {
				fileName, _ := environment.EnvironmentTomlFilename()
				return Settings{}, fmt.Errorf("could not resolve external package name '%v', please add it to '%v' file", packageName, fileName)
			}
			convertedPath = convertedPath + suffix
			mapped = ModuleFromEnvironment
		} else {
			mapped = ModuleFromPath
		}
		if !filepath.IsAbs(convertedPath) {
			convertedPath = path.Join(rootDirectory, convertedPath)
		}

		settings.Module[index].Path = convertedPath
		settings.Module[index].Mapped = mapped
	}

	return settings, nil
}
