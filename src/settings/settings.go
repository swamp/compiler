/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package settings

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Module struct {
	Name string
	Path string
}

type Settings struct {
	Name   string
	Module []Module
}

func Load(reader io.Reader, rootDirectory string) (Settings, error) {
	data, dataErr := ioutil.ReadAll(reader)
	if dataErr != nil {
		return Settings{}, dataErr
	}

	settings := Settings{}

	unmarshalErr := toml.Unmarshal(data, &settings)
	if unmarshalErr != nil {
		return Settings{}, unmarshalErr
	}

	for index, mod := range settings.Module {
		convertedPath := os.ExpandEnv(mod.Path)
		if !filepath.IsAbs(convertedPath) {
			convertedPath = filepath.Join(rootDirectory, convertedPath)
		}
		settings.Module[index].Path = convertedPath
	}

	return settings, nil
}
