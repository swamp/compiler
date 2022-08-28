/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package solution

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pelletier/go-toml"
	"github.com/swamp/compiler/src/file"
)

type Package struct {
	Name string
	Path string
}

type Settings struct {
	Name     string
	Packages []string
}

func Load(reader io.Reader, solutionFileDirectory string) (Settings, error) {
	data, dataErr := io.ReadAll(reader)
	if dataErr != nil {
		return Settings{}, dataErr
	}

	settings := Settings{}

	unmarshalErr := toml.Unmarshal(data, &settings)
	if unmarshalErr != nil {
		return Settings{}, unmarshalErr
	}

	return settings, nil
}

func fileFromDirectory(rootDirectory string) string {
	return path.Join(rootDirectory, "swamp.solution.toml")
}

func LoadIfExists(rootDirectory string) (Settings, error) {
	tomlFilename := fileFromDirectory(rootDirectory)
	if !file.HasFile(tomlFilename) {
		return Settings{}, fmt.Errorf("didn't find solution file in %s", rootDirectory)
	}

	tomlFile, settingsFileErr := os.Open(tomlFilename)
	if settingsFileErr != nil {
		return Settings{}, fmt.Errorf("couldn't load solution file %s %w", tomlFilename, settingsFileErr)
	}

	settingsReader := bufio.NewReader(tomlFile)

	return Load(settingsReader, rootDirectory)
}
