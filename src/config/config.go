/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package settings

import (
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

func getHomeConfigDirectory() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	completePath := path.Join(dirname, ".config/swamp/")
	log.Printf("found path %v", completePath)

	if !file.IsDir(completePath) {
		os.MkdirAll(completePath, 0o755)
	}

	return "", fmt.Errorf("not a ")
}

func LoadFromHome() (Config, error) {
	configDirectory, er := getHomeConfigDirectory()
	configFile := path.Join(completePath, "env.toml")

	if !file.HasFile(configFile) {
		return ""
	}

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
