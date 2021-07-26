/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package config

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/pelletier/go-toml"
	"github.com/swamp/compiler/src/file"
)

type PackagePath struct {
	Name string
	Path string
}

type Config struct {
	Version     string
	PackagePath []PackagePath
}

func (c Config) Lookup(name string) string {
	for _, x := range c.PackagePath {
		if x.Name == name {
			return x.Path
		}
	}

	return ""
}

func Load(reader io.Reader) (Config, error) {
	data, dataErr := ioutil.ReadAll(reader)
	if dataErr != nil {
		return Config{}, dataErr
	}

	config := Config{}

	unmarshalErr := toml.Unmarshal(data, &config)
	if unmarshalErr != nil {
		return Config{}, unmarshalErr
	}

	return config, nil
}

func getConfigFilename() (string, error) {
	configTomlFilename, err := ConfigTomlFilename()
	if err != nil {
		return "", err
	}

	parentDirectory := path.Dir(configTomlFilename)
	if !file.IsDir(parentDirectory) {
		log.Printf("creating config directory '%v'", parentDirectory)
		os.MkdirAll(parentDirectory, 0o755)
	}

	return configTomlFilename, nil
}

func ConfigTomlFilename() (string, error) {
	configDirectoryName, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	completePath := path.Join(configDirectoryName, "swamp/")
	configTomlFile := path.Join(completePath, "env.toml")
	return configTomlFile, nil
}

func LoadFromConfig() (Config, bool, error) {
	configTomlFile, err := getConfigFilename()
	if err != nil {
		return Config{}, false, err
	}

	if !file.HasFile(configTomlFile) {
		return Config{}, false, nil
	}

	tomlFile, settingsFileErr := os.Open(configTomlFile)
	if settingsFileErr != nil {
		return Config{}, false, fmt.Errorf("couldn't load config file %s %w", configTomlFile, settingsFileErr)
	}

	configReader := bufio.NewReader(tomlFile)

	config, err := Load(configReader)
	if err != nil {
		return Config{}, false, err
	}

	return config, true, err
}
