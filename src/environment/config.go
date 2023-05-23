/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package environment

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	"github.com/swamp/compiler/src/file"
)

type Package struct {
	Name string
	Path string
}

type Environment struct {
	Version string
	Package []*Package
}

func (c Environment) Lookup(name string) string {
	for _, x := range c.Package {
		if x.Name == name {
			return x.Path
		}
	}

	return ""
}

func (c *Environment) AddOrSet(name string, path string) {
	for _, x := range c.Package {
		if x.Name == name {
			x.Path = path
			return
		}
	}

	c.Package = append(c.Package, &Package{Name: name, Path: path})
}

func Load(reader io.Reader) (Environment, error) {
	data, dataErr := io.ReadAll(reader)
	if dataErr != nil {
		return Environment{}, dataErr
	}

	config := Environment{}

	unmarshalErr := toml.Unmarshal(data, &config)
	if unmarshalErr != nil {
		return Environment{}, unmarshalErr
	}

	for _, entry := range config.Package {
		entry.Path = filepath.ToSlash(entry.Path)
	}

	return config, nil
}

func getConfigFilename() (string, error) {
	configTomlFilename, err := EnvironmentTomlFilename()
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

func EnvironmentTomlFilename() (string, error) {
	configDirectoryName, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	completePath := path.Join(configDirectoryName, "swamp/")
	configTomlFile := path.Join(completePath, "env.toml")
	return configTomlFile, nil
}

func LoadFromConfig() (Environment, bool, error) {
	configTomlFile, err := getConfigFilename()
	if err != nil {
		return Environment{}, false, err
	}

	if !file.HasFile(configTomlFile) {
		return Environment{}, false, nil
	}

	tomlFile, settingsFileErr := os.Open(configTomlFile)
	if settingsFileErr != nil {
		return Environment{}, false, fmt.Errorf("couldn't load config file %s %w", configTomlFile, settingsFileErr)
	}

	configReader := bufio.NewReader(tomlFile)

	config, err := Load(configReader)
	if err != nil {
		return Environment{}, false, err
	}

	return config, true, err
}

func (c Environment) SaveToConfig() error {
	configTomlFile, err := getConfigFilename()
	if err != nil {
		return err
	}

	data, marshalErr := toml.Marshal(&c)
	if marshalErr != nil {
		return marshalErr
	}

	return os.WriteFile(configTomlFile, data, 0o755)
}
