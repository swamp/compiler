/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package loader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

type Loader struct {
	rootPath string
}

func moduleNameToRelativeFilePath(moduleName dectype.PackageRelativeModuleName) string {
	fixed := strings.Replace(moduleName.String(), ".", "/", -1)
	if fixed == "" {
		fixed += "Main"
	}
	if fixed != "" {
		fixed = fixed + ".swamp"
	}

	return fixed
}

func NewLoader(rootPath string) *Loader {
	return &Loader{rootPath: rootPath}
}

func (l *Loader) Load(relativeModuleName dectype.PackageRelativeModuleName, verboseFlag bool) (string, string, decshared.DecoratedError) {
	relativePath := moduleNameToRelativeFilePath(relativeModuleName)
	if filepath.IsAbs(relativePath) {
		return "", "", decorated.NewInternalError(fmt.Errorf("loader wants relative paths, can not use absolute ones '%s'", relativeModuleName))
	}
	fullPath := filepath.Join(l.rootPath, relativePath)
	completeFilename, completeFilenameErr := filepath.Abs(fullPath)
	if completeFilenameErr != nil {
		return "", "", decorated.NewInternalError(completeFilenameErr)
	}
	if verboseFlag {
		fmt.Printf("* loading file %v\n", completeFilename)
	}

	octets, readFileErr := ioutil.ReadFile(completeFilename)
	if readFileErr != nil {
		switch readFileErr {
		case os.ErrInvalid:
		case os.ErrPermission:
		case os.ErrNotExist:
			return "", "", decorated.NewInternalError(fmt.Errorf("file '%s' didn't exist (%v)", completeFilename, readFileErr))
		default:
			switch v := readFileErr.(type) {
			case *os.PathError:
				return "", "", decorated.NewInternalError(fmt.Errorf("couldn't find relative module '%v', root '%v', file '%s' (%v)", relativeModuleName, l.rootPath, completeFilename, v))
			default:
				return "", "", decorated.NewInternalError(fmt.Errorf("couldn't open file '%s'", completeFilename))
			}
		}
	}
	str := string(octets)

	return fullPath, str, nil
}
